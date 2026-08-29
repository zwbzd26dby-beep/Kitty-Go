# Kitty-Go

Agentowa warstwa orkiestracji AI w Go. System abstrahuje modele, providerów,
koszty, limity, pamięć, narzędzia, urządzenia i zasoby obliczeniowe.

Implementacja zgodna z dokumentem **MASTER ARCHITECTURE** (20 faz). Kod jest
rozwijany fazami; obecny stan: **Faza 0** (fundament: `pkg/types`, `pkg/llm`,
`cmd/kitty-repl`, pełny scaffold katalogów).

## Struktura

- `cmd/` — entry pointy: `kitty-cli`, `kitty-repl`, `kitty-api`, `kitty-gui`
- `internal/` — moduły wewnętrzne (agent, orchestrator, decision, router,
  execution, provider, model, cost, limit, budget, compute, memory, tools,
  security, config, observability, interface)
- `pkg/` — pakiety publiczne: `types`, `llm`, `utils`
- `tests/` — testy: unit, integration, provider, router, cost, limit, memory,
  tools, security, compute, distributed, cli, e2e
- `configs/` — pliki konfiguracyjne (default, przykład user, secrets)
- `scripts/` — skrypty build/dev/deply

## Uruchomienie (Faza 0)

```bash
# REPL z mock-providerem (history echo)
go run ./cmd/kitty-repl

# CLI one-shot
go run ./cmd/kitty-cli "Hej Kitty"
```

## Testy

```bash
make test
# albo
go test ./...
```

## Fazy

| Faza | Zakres | Status |
|------|--------|--------|
| 0 | Fundament: types, llm (mock), kitty-repl | ✅ w trakcie |
| 1 | Warstwy: agent / orchestrator / execution | ⏳ |
| 2 | Provider System (+ OpenCode Zen `big-pickle`) | ⏳ |
| 3+ | pełna roadmapa 20 faz (patrz docs) | ⏳ |

Szczegóły: `docs/architecture.md`, `docs/development.md`.
