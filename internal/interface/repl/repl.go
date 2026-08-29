// Package repl implements the interactive REPL engine. It is decoupled from
// the LLM client so different providers can be plugged in (Phase 2+).
package repl

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/llm"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"
)

// Engine drives an interactive session over the provided reader/writer.
type Engine struct {
	client   llm.Client
	in       *bufio.Reader
	out      io.Writer
	conv     *types.Conversation
	ctx      context.Context
	commands map[string]CommandHandler
}

// CommandHandler processes a REPL command (arguments passed without the slash).
type CommandHandler func(args []string) error

// NewEngine creates a REPL engine with the given client and I/O streams.
func NewEngine(client llm.Client, r io.Reader, w io.Writer) *Engine {
	e := &Engine{
		client:   client,
		in:       bufio.NewReader(r),
		out:      w,
		conv:     types.NewConversation(),
		ctx:      context.Background(),
		commands: map[string]CommandHandler{},
	}
	e.registerBuiltins()
	return e
}

func (e *Engine) registerBuiltins() {
	e.commands["help"] = e.cmdHelp
	e.commands["clear"] = e.cmdClear
	e.commands["exit"] = e.cmdExit
	e.commands["quit"] = e.cmdExit
}

// RegisterCommand adds or overrides a slash command handler.
func (e *Engine) RegisterCommand(name string, fn CommandHandler) {
	e.commands[strings.ToLower(name)] = fn
}

func (e *Engine) cmdHelp([]string) error {
	fmt.Fprintln(e.out, "Komendy: /help /clear /exit /quit")
	return nil
}

func (e *Engine) cmdClear([]string) error {
	e.conv.Clear()
	fmt.Fprintln(e.out, "Historia wyczyszczona.")
	return nil
}

func (e *Engine) cmdExit([]string) error {
	return errExit{}
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
	msg, err := types.NewMessage(userInput)
	if err != nil {
		return err
	}
	e.conv.AddMessage(msg)
	history := e.conv.GetHistory()
	resp, err := e.client.Generate(e.ctx, msg, history)
	if err != nil {
		return err
	}
	e.conv.AddResponse(resp)
	fmt.Fprintf(e.out, "Kitty: %s\n", resp.Content())
	return nil
}
