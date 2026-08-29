# Architektura Kitty-Go

## Warstwy

```
┌──────────────────────────────────────────────┐
│ interface: cli / repl / api / gui            │
├──────────────────────────────────────────────┤
│ agent → orchestrator                         │
│   decision → router/model → router/compute   │
│   → execution (local / remote / distributed)  │
├──────────────────────────────────────────────┤
│ provider / model / cost / limit / budget     │
│ compute / memory / tools / security / config │
│ observability                                 │
├──────────────────────────────────────────────┤
│ pkg/types, pkg/llm, pkg/utils                │
└──────────────────────────────────────────────┘
```

## Przepływ żądania

```
Agent.Process(userInput)
 → Orchestrator.Orchestrate(Task)
 → DecisionEngine.Decide(Task)
 → ModelRouter.Select(Decision)
 → ComputeRouter.Select(Model, Task)
 → Executor.Execute(Task, Model, Device)
 → Provider.Complete(Request) / DistributedExecutor
 → Response
```

## Rozszerzanie

Każda faza rozwija swój pakiet w `internal/`, eksportuje interfejs i dostarcza
testy w `tests/<obszar>/`. Komunikacja między modułami wyłącznie przez interfejsy
(sekcja 31 Master Architecture). Suite zawsze green przed przejściem dalej.

## Faza 0 (obecna)

Realizuje fundament analogiczny do Python Phase 1:

- `pkg/types` — `Message`, `Response`, `Conversation`, `Task`
- `pkg/llm` — `Client` (interfejs), `MockClient`, `HistoryEchoClient`
- `cmd/kitty-repl` — REPL z historią
