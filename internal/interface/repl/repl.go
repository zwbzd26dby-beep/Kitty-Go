// Package repl implements the interactive REPL engine. From Phase 1 the
// engine drives an agent.Agent backed by the full Core stack
// (Agent → Orchestrator → Executor → provider), so providers can be swapped.
package repl

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/agent"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/execution"
	reg "github.com/zwbzd26dby-beep/Kitty-Go/internal/model"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/orchestrator"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/llm"
)

// Engine drives an interactive session over the provided reader/writer.
type Engine struct {
	agent    agent.Agent
	core     *Core
	in       *bufio.Reader
	out      io.Writer
	ctx      context.Context
	commands map[string]CommandHandler
	model    string
}

// CommandHandler processes a REPL command (arguments passed without the slash).
type CommandHandler func(args []string) error

// NewEngineWithCore creates a REPL engine driven by an agent.Agent and
// exposing the subsystem handles through slash commands.
func NewEngineWithCore(core *Core, r io.Reader, w io.Writer) *Engine {
	e := &Engine{
		agent:    core.Agent,
		core:     core,
		in:       bufio.NewReader(r),
		out:      w,
		ctx:      context.Background(),
		commands: map[string]CommandHandler{},
	}
	e.registerBuiltins()
	return e
}

// NewEngineWithAgent creates a REPL engine driven by an agent.Agent.
func NewEngineWithAgent(a agent.Agent, r io.Reader, w io.Writer) *Engine {
	e := &Engine{
		agent:    a,
		in:       bufio.NewReader(r),
		out:      w,
		ctx:      context.Background(),
		commands: map[string]CommandHandler{},
	}
	e.registerBuiltins()
	return e
}

// NewEngine creates a REPL engine on the full default stack built around the
// given LLM client: Agent → Orchestrator → Executor(Local) → client.
func NewEngine(client llm.Client, r io.Reader, w io.Writer) *Engine {
	exec := execution.NewLocalExecutor()
	orch := orchestrator.New(exec, client)
	a := agent.New(orch, agent.DefaultOptions{})
	return NewEngineWithAgent(a, r, w)
}

func (e *Engine) registerBuiltins() {
	e.commands["help"] = e.cmdHelp
	e.commands["clear"] = e.cmdClear
	e.commands["exit"] = e.cmdExit
	e.commands["quit"] = e.cmdExit
	e.commands["model"] = e.cmdModel
	e.commands["models"] = e.cmdModels
	e.commands["budget"] = e.cmdBudget
	e.commands["cost"] = e.cmdCost
	e.commands["tools"] = e.cmdTools
	e.commands["memory"] = e.cmdMemory
	e.commands["devices"] = e.cmdDevices
	e.commands["metrics"] = e.cmdMetrics
	e.commands["security"] = e.cmdSecurity
}

// RegisterCommand adds or overrides a slash command handler.
func (e *Engine) RegisterCommand(name string, fn CommandHandler) {
	e.commands[strings.ToLower(name)] = fn
}

func (e *Engine) cmdHelp([]string) error {
	fmt.Fprintln(e.out, "Komendy: /help /clear /model /models /budget /cost /tools /memory /devices /metrics /security /exit /quit")
	return nil
}

func (e *Engine) cmdClear([]string) error {
	e.agent.ClearHistory()
	if e.core != nil && e.core.Memory != nil {
		e.core.Memory.Clear()
	}
	fmt.Fprintln(e.out, "Historia wyczyszczona.")
	return nil
}

func (e *Engine) cmdExit([]string) error {
	return errExit{}
}

func (e *Engine) cmdModel(args []string) error {
	if e.core == nil || e.core.Models == nil {
		fmt.Fprintln(e.out, "Brak registry modeli.")
		return nil
	}
	if len(args) > 0 {
		if !modelExists(e.core.Models, args[0]) {
			fmt.Fprintf(e.out, "Nieznany model %q\n", args[0])
			return nil
		}
		e.model = args[0]
		fmt.Fprintf(e.out, "Model: %s\n", args[0])
		return nil
	}
	if e.model == "" {
		fmt.Fprintln(e.out, "Brak aktywnego modelu.")
		return nil
	}
	fmt.Fprintf(e.out, "Model: %s\n", e.model)
	return nil
}

func modelExists(r *reg.Registry, id string) bool {
	for _, m := range r.List() {
		if m.Model == id {
			return true
		}
	}
	return false
}

func (e *Engine) cmdModels([]string) error {
	if e.core == nil || e.core.Models == nil {
		fmt.Fprintln(e.out, "Brak registry modeli.")
		return nil
	}
	for _, m := range e.core.Models.List() {
		fmt.Fprintf(e.out, "%s/%s — %s\n", m.Provider, m.Model, m.Description)
	}
	return nil
}

func (e *Engine) cmdBudget([]string) error {
	if e.core == nil || e.core.Budget == nil {
		fmt.Fprintln(e.out, "Brak managera budżetu.")
		return nil
	}
	d, m := e.core.Budget.GetRemaining()
	fmt.Fprintf(e.out, "Pozostały budżet: dzienny %.4f, miesięczny %.4f\n", d, m)
	return nil
}

func (e *Engine) cmdCost([]string) error {
	if e.core == nil || e.core.Cost == nil {
		fmt.Fprintln(e.out, "Brak managera kosztów.")
		return nil
	}
	fmt.Fprintf(e.out, "Całkowity koszt: %.4f\n", e.core.Cost.GetTotalCost())
	return nil
}

func (e *Engine) cmdTools([]string) error {
	if e.core == nil || e.core.Tools == nil {
		fmt.Fprintln(e.out, "Brak registry narzędzi.")
		return nil
	}
	for _, t := range e.core.Tools.List() {
		status := "dost."
		if !e.core.Tools.IsAllowed(t.Name()) {
			status = "blok."
		}
		fmt.Fprintf(e.out, "%s [%s] — %s\n", t.Name(), status, t.Description())
	}
	return nil
}

func (e *Engine) cmdMemory([]string) error {
	if e.core == nil || e.core.Memory == nil {
		fmt.Fprintln(e.out, "Brak pamięci krótkoterminowej.")
		return nil
	}
	items := e.core.Memory.Recent(10)
	if len(items) == 0 {
		fmt.Fprintln(e.out, "Pamięć pusta.")
		return nil
	}
	for i, it := range items {
		fmt.Fprintf(e.out, "%02d: %s\n", i+1, it)
	}
	return nil
}

func (e *Engine) cmdDevices([]string) error {
	if e.core == nil || e.core.Devices == nil {
		fmt.Fprintln(e.out, "Brak registry urządzeń.")
		return nil
	}
	devs := e.core.Devices.List()
	if len(devs) == 0 {
		fmt.Fprintln(e.out, "Brak urządzeń.")
		return nil
	}
	for _, d := range devs {
		fmt.Fprintf(e.out, "%s (%s) [%s]\n", d.ID, d.Address, d.Health)
	}
	return nil
}

func (e *Engine) cmdMetrics([]string) error {
	if e.core == nil || e.core.Obs == nil {
		fmt.Fprintln(e.out, "Brak metryk.")
		return nil
	}
	cs, gs := e.core.Obs.Metrics.Snapshot()
	if len(cs) == 0 && len(gs) == 0 {
		fmt.Fprintln(e.out, "Brak zarejestrowanych metryk.")
		return nil
	}
	for k, v := range cs {
		fmt.Fprintf(e.out, "counter %s = %d\n", k, v)
	}
	for k, v := range gs {
		fmt.Fprintf(e.out, "gauge %s = %.2f\n", k, v)
	}
	return nil
}

func (e *Engine) cmdSecurity([]string) error {
	if e.core == nil || e.core.Security == nil {
		fmt.Fprintln(e.out, "Brak managera bezpieczeństwa.")
		return nil
	}
	fmt.Fprintln(e.out, "Sekrety: zamaskowane, nigdy nie logowane.")
	ev := e.core.Obs.Audit.Events()
	if len(ev) > 0 {
		fmt.Fprintf(e.out, "Zdarzenia audytowe: %d\n", len(ev))
	} else {
		fmt.Fprintln(e.out, "Zdarzenia audytowe: brak")
	}
	return nil
}

// errExit signals a clean REPL exit.
type errExit struct{}

func (errExit) Error() string { return "exit" }

var exitSentinel = errExit{}

// Run starts the interactive loop. It returns nil on exit, or an error.
func (e *Engine) Run() error {
	greet := "Hej Kitty (wpisz /help). Pusta linia kończy."
	fmt.Fprintln(e.out, greet)

	for {
		fmt.Fprint(e.out, "Ty: ")
		line, err := e.readLine()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if strings.TrimSpace(line) == "" {
			return nil
		}
		if strings.HasPrefix(line, "/") {
			if stop, cerr := e.handleCommand(line); stop {
				return cerr
			}
			continue
		}
		if err := e.respond(line); err != nil {
			fmt.Fprintf(e.out, "Błąd: %v\n", err)
		}
	}
}

func (e *Engine) readLine() (string, error) {
	return e.in.ReadString('\n')
}

// handleCommand returns (stop, err). stop indicates the loop should end.
func (e *Engine) handleCommand(line string) (bool, error) {
	parts := strings.Fields(strings.TrimPrefix(line, "/"))
	if len(parts) == 0 {
		return false, nil
	}
	name := strings.ToLower(parts[0])
	fn, ok := e.commands[name]
	if !ok {
		fmt.Fprintf(e.out, "Nieznana komenda /%s. Użyj /help.\n", name)
		return false, nil
	}
	err := fn(parts[1:])
	if _, isExit := err.(errExit); isExit {
		return true, nil
	}
	if err != nil {
		fmt.Fprintf(e.out, "Błąd: %v\n", err)
	}
	return false, nil
}

func (e *Engine) respond(userInput string) error {
	resp, err := e.agent.Process(e.ctx, userInput)
	if err != nil {
		return err
	}
	fmt.Fprintf(e.out, "Kitty: %s\n", resp.Content)
	return nil
}
