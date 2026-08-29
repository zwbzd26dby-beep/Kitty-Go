# Kitty (Go) — Specyfikacja Faz 0–20

Uzupełnienie do `HANDOFF_GO_TRANSITION.md`. Interfejsy Go referencjonowane
poniżej (np. `Provider`, `ModelRouter`) są zdefiniowane w pełni w dokumencie
**MASTER ARCHITECTURE**, sekcje 4–24. Kolumna "Python ref." wskazuje
odpowiadający moduł/testy w repo Python, które służą jako specyfikacja
behawioralna (przypadki brzegowe, nie kod do kopiowania).

Zasada przejścia: faza ukończona gdy odpowiadający zestaw testów w `tests/<obszar>/`
jest zielony i moduł wyeksportowany zgodnie z `internal/`/`pkg/` z §3 Master
Architecture.

---

## Faza 0 — Fundament

**Cel:** minimalny szkielet end-to-end: REPL + mock LLM + typy bazowe.

**Pakiety:** `cmd/kitty-repl/`, `pkg/types/{message,response,conversation,task}.go`, `pkg/llm/{client,mock}.go`

**Zadania:**
- `pkg/types` — struktury `Message`, `Response`, `Conversation`, `Task`.
- `pkg/llm.Client` interfejs + `MockClient` (`MockClient` deterministyczny +
  `HistoryEchoClient` — **decyzja: dwa mocki**, 1:1 z Pythonem).
- `cmd/kitty-repl` — pętla read-eval-print wołająca `MockClient`.
- `go.mod` init, `Makefile` (build/test/run), `README.md`.

**Testy:** `tests/unit/types_test.go`, `tests/unit/mock_client_test.go`, smoke REPL.

**Python ref.:** Ph1 (`test_phase1.py`).

---

## Faza 1 — Core Architecture (warstwy)

**Cel:** wydzielenie Agent Core / Orchestrator / Execution.

**Pakiety:** `internal/agent/{agent,intent,context,history}.go`, `internal/orchestrator/{orchestrator,decision,coordinator}.go`, `internal/execution/{executor,local,job,status}.go`

**Zadania:**
- `agent.Agent` interfejs (sekcja 5): `Process`, `GetConversation`, `ClearHistory`.
- `orchestrator.Orchestrator` (sekcja 6) ze szkieletem `Orchestrate` (pass-through do `MockClient` przez `Executor`).
- `execution.Executor` — na razie tylko `ExecuteLocal`.
- Przepięcie `cmd/kitty-repl` na nowy stos: `REPL → Agent → Orchestrator → Executor → mock`.

**Testy:** `tests/unit/agent_test.go`, `tests/integration/agent_orchestrator_test.go`.

**Python ref.:** Ph3 (`test_phase3.py`).

---

## Faza 2 — Provider System

**Cel:** provider-agnostic warstwa komunikacji z LLM, w tym **OpenCode Zen** (wymaganie spoza dokumentu 1).

**Pakiety:** `internal/provider/{provider,openai,kimi,deepseek,openrouter,ollama,opencode}/`

**Zadania:**
- `provider.Provider` (sekcja 9): `Complete`, `Stream`, `GetPricing`, `GetLimits`, `HealthCheck`.
- `OpenAIProvider`, `KimiProvider`, `DeepSeekProvider`, `OpenRouterProvider`, `OllamaProvider`.
- **`OpenCodeProvider`** — OpenAI-kompatybilny, base URL `https://opencode.ai/zen/v1`, modele `big-pickle` (kod/reasoning), `mimo-v2-pro-free`, `minimax-m2.5-free`; klucz `OPENCODE_API_KEY`; zadania kodowe → `big-pickle`.
- Migracja `MockClient` do `provider.Provider` jako `MockProvider`.

**Testy:** `tests/provider/*_test.go` per provider (mockowany transport HTTP), `tests/provider/opencode_test.go`.

**Python ref.:** `providers/`, `test_providers_opencode.py`.

---

## Faza 3 — Real LLM Integration

**Cel:** realne API z obsługą błędów/retry/timeout.

**Zadania:**
- Współdzielony HTTP klient (`internal/provider/http.go`) z retry + exponential backoff.
- Testy integracyjne z realnymi kluczami (pomijane w CI bez sekretów).
- `context.Context` z timeoutem na każdym `Provider.Complete`.

**Testy:** `tests/integration/live_provider_test.go` (skip bez env), `tests/unit/retry_test.go`.

**Python ref.:** brak dedykowanej fazy; retry w `providers/http.py`.

---

## Faza 4 — Configuration & Secrets

**Cel:** trójpoziomowa konfiguracja (default/user/secret) + Security Manager.

**Pakiety:** `internal/config/{default,user,secret,loader,validator}.go`, `internal/security/{manager,apikey,secrets}.go`

**Zadania:**
- `config.Config` (sekcja 22).
- Ładowanie z `~/.kitty/config.yaml` (decyzja: dopuszczamy jedną zależność YAML, np. `gopkg.in/yaml.v3`).
- `SecurityManager.GetAPIKey/SetAPIKey` z fallbackiem do zmiennych środowiskowych.
- Walidacja configu przy starcie.

**Testy:** `tests/unit/config_test.go`, `tests/security/apikey_test.go`.

**Python ref.:** Ph4; uwaga o braku PyYAML w Pythonie — w Go decyzja otwarta (dopuszczamy zależność).

---

## Faza 5 — Model Registry

**Cel:** centralny rejestr modeli z metadanymi.

**Pakiety:** `internal/model/{registry,model,capabilities,pricing}.go`

**Zadania:**
- `model.ModelRegistry` (sekcja 10): `Register`, `Get`, `List`, `UpdatePricing`, `UpdateAvailability`.
- Statyczna rejestracja modeli (w tym OpenCode Zen).
- Hook dynamicznego pobierania metadanych z providera.

**Testy:** `tests/router/registry_test.go`.

**Python ref.:** Ph5.

---

## Faza 6 — Model Router

**Cel:** wybór modelu na podstawie `Decision` (koszt/jakość/szybkość/dostępność, scoring, fallback).

**Pakiety:** `internal/router/model/{router,selector,cost_aware,fallback}.go`

**Zadania:**
- `router.ModelRouter` (sekcja 8), algorytm filtrowania + scoringu (sekcja 33).
- `GetFallback` — łańcuch alternatyw.
- Zadania kodowe preferują `big-pickle` (przeniesione z Fazy 2 — właściwe miejsce logiki routingu).

**Testy:** `tests/router/model_router_test.go` — w tym przypadki: brak dostępnych modeli, przekroczony budżet, przekroczony limit.

**Python ref.:** Ph5/Ph6 (`routing/`).

---

## Faza 7 — Cost & Limits & Budget

**Cel:** kalkulacja/śledzenie kosztów, rate limiting, budżety z blokowaniem.

**Pakiety:** `internal/cost/{manager,calculator,tracker,alert}.go`, `internal/limit/{manager,rate,quota,retry}.go`, `internal/budget/{manager,daily,monthly,block}.go`

**Zadania:**
- `CostManager`, `LimitManager`, `BudgetManager` (sekcje 11–13).
- Przepływ (sekcja 35): `Complete` → koszt → `TrackCost` → budżet → alert/block.
- Komendy REPL `/budget`, `/cost`.

**Testy:** `tests/cost/*_test.go`, `tests/limit/*_test.go` (w tym blokada po przekroczeniu budżetu).

**Python ref.:** Ph6.

---

## Faza 8 — Compute Fabric & Device Registry

**Cel:** pula urządzeń — rejestracja, discovery w LAN, health monitoring.

**Pakiety:** `internal/compute/{fabric,device_registry,discovery,health,resource,worker}.go`

**Zadania:**
- `ComputeFabric`, `DeviceRegistry` (sekcje 15–16).
- Discovery: statycznie (config), potem LAN (mDNS / broadcast UDP).
- Health check (goroutine + ticker).

**Testy:** `tests/compute/device_registry_test.go`, `tests/compute/health_test.go`.

**Python ref.:** Ph9.

---

## Faza 9 — Distributed Execution

**Cel:** wykonywanie zadań na zdalnych urządzeniach.

**Pakiety:** `internal/execution/remote.go` + transport w `internal/compute/`.

**Zadania:**
- `DistributedExecutor` (sekcja 17): `Execute`, `Stream`, `Ping`, `Authenticate`.
- Serializacja `Job`/`JobResult` (gob/JSON przez HTTP/gRPC — do ustalenia).
- Autentykacja urządzeń przez `SecurityManager`.

**Testy:** `tests/distributed/remote_exec_test.go` (symulowane zdalne urządzenie).

**Python ref.:** Ph10.

---

## Faza 10 — Compute Router

**Cel:** wybór urządzenia (zasoby, obciążenie, latency, fallback).

**Pakiety:** `internal/router/compute/{router,selector,resource_aware,scheduler}.go`

**Zadania:**
- `ComputeRouter` (sekcja 14), algorytm (sekcja 34): filtrowanie VRAM/RAM, scoring po obciążeniu/latency, fallback.

**Testy:** `tests/router/compute_router_test.go`.

**Python ref.:** Ph7 (część compute).

---

## Faza 11 — Execution Layer (pełny)

**Cel:** jednolita warstwa wykonawcza (lokalne + zdalne, retry/fallback/timeout).

**Pakiety:** `internal/execution/` (rozbudowa Fazy 1 i 9)

**Zadania:**
- Pełna `Executor.Execute` z fallbackami (sekcja 36).
- Ujednolicenie ścieżki lokalnej i zdalnej.

**Testy:** `tests/e2e/execution_full_test.go` — sukces, retry+sukces, fallback+sukces, porażka.

**Python ref.:** Ph7.

---

## Faza 12 — Memory

**Cel:** pamięć krótkoterminowa, długoterminowa, semantyczna + RAG.

**Pakiety:** `internal/memory/{short_term,long_term,semantic,embeddings,rag}.go`

**Zadania:**
- `Memory`, `ShortTermMemory`, `LongTermMemory`, `SemanticMemory` (sekcja 19).
- Embeddingi: lokalny model (Ollama) lub TF-IDF jako fallback.
- `SimilaritySearch` — brute-force cosine similarity na start.

**Testy:** `tests/memory/*_test.go` (w tym RAG end-to-end).

**Python ref.:** Ph8.

---

## Faza 13 — Tools

**Cel:** system narzędzi: terminal, filesystem, web, Git, GitHub, Docker, SSH, sysinfo.

**Pakiety:** `internal/tools/{registry,permissions,sandbox,terminal,filesystem,web,git,github,docker,ssh,sysinfo}.go`

**Zadania:**
- `Tool`, `ToolRegistry` (sekcja 20).
- Sandboxing dla terminal/filesystem.
- Uprawnienia per narzędzie przez `SecurityManager`.
- Komendy REPL `/tools`, `/tool <nazwa>`.

**Testy:** `tests/tools/*_test.go` (w tym próba wyjścia poza sandbox musi się nie udać).

**Python ref.:** Ph11.

---

## Faza 14 — Security

**Cel:** sekrety, uprawnienia, autentykacja urządzeń, audyt.

**Pakiety:** `internal/security/{manager,apikey,secrets,auth,permissions,audit}.go`

**Zadania:**
- Rozbudowa `SecurityManager` (z Fazy 4): `ValidatePermissions`, `AuthenticateDevice`, `CreateAuditLog`.
- Audit log append-only (JSON lines).
- Żaden klucz API nigdy nie trafia do repo/logów w czystej postaci.

**Testy:** `tests/security/*_test.go` (audit rejestruje każde wywołanie narzędzia).

**Python ref.:** Ph12.

---

## Faza 15 — Observability

**Cel:** logging, metryki, tracing, usage/cost, audyt — cross-cutting.

**Pakiety:** `internal/observability/{logger,metrics,tracer,usage,cost,audit}.go`

**Zadania:**
- `Observability` (sekcja 23).
- Integracja ze wszystkimi modułami (hooki/middleware).
- Komenda REPL `/metrics`.

**Testy:** `tests/observability/*_test.go`.

**Python ref.:** Ph13.

---

## Faza 16 — Decision Engine

**Cel:** analiza zadania → `Decision` (wymagania, priorytet, budżet, capabilities).

**Pakiety:** `internal/decision/{engine,rules,priority,capabilities}.go`

**Zadania:**
- `DecisionEngine` (sekcja 7) — reguły + wagi + scoring (nie jeden if/else).
- Podłączenie do `Orchestrator` w miejsce pass-through z Fazy 1.

**Testy:** `tests/unit/decision_engine_test.go`.

**Python ref.:** Ph3/Ph6 (`core/policy`+`core/decision`).

---

## Faza 17 — Orchestrator (pełny)

**Cel:** pełna koordynacja Decision → Model Router → Compute Router → Execution.

**Zadania:**
- Zamiana uproszczonego `Orchestrate` z Fazy 1 na pełny przepływ (sekcja 32).
- Cost/Limit/Budget checks przed i po wykonaniu.
- Obsługa błędów (sekcja 37).

**Testy:** `tests/e2e/orchestrator_full_test.go`.

**Python ref.:** spięcie wszystkich faz Python.

---

## Faza 18 — CLI / REPL / API

**Cel:** wiele interfejsów na wspólnym Core.

**Pakiety:** `cmd/kitty-{cli,repl,api}/`, `internal/interface/{cli,repl,api}/`

**Zadania:**
- REPL: `/help /clear /model /models /budget /cost /tools /tool /memory /devices /metrics /security /exit /quit` — 1:1 z Python UX.
- API: `GET /health`, `GET /metrics`, `POST /agent/respond`, Bearer (`KITTY_API_TOKEN`), `KITTY_HOST`/`KITTY_PORT`.
- CLI one-shot.

**Testy:** `tests/cli/*_test.go`, `tests/integration/api_test.go`.

**Python ref.:** Ph14.

---

## Faza 19 — GUI

**Cel:** interfejs graficzny (TUI lub prosty web) na wspólnym Core.

**Pakiety:** `cmd/kitty-gui/`, `internal/interface/gui/`

**Zadania:**
- Odtworzenie menu/czat z Python `gui/` (zakres minimalny, bez rozbudowy).

**Testy:** `tests/cli/gui_smoke_test.go`.

**Python ref.:** Ph14 (`gui/`).

---

## Faza 20 — Production Hardening

**Cel:** gotowość produkcyjna.

**Zadania:**
- `Dockerfile`, `docker-compose.yml`, `entrypoint.sh` — wzór z Python.
- CI: GitHub Actions — build + `go test ./...` (matrix wersji Go).
- `docs/architecture.md`, `docs/api.md`, `docs/configuration.md`, `docs/deployment.md`, `docs/development.md`.
- Profiling (`pprof`), optymalizacja gorących ścieżek.

**Testy:** pełny `tests/e2e/` zielony w CI na dwóch wersjach Go.

**Python ref.:** Ph15.

---

## Dependencje między fazami

- **Faza 6** potrzebuje Fazy 5 i częściowo Fazy 7 (stuby Limit/Budget).
- **Fazy 10** i **16** zasilają Fazę 17 — mogą być równoległe.
- **Faza 15** wpinana przyrostowo od Fazy 7 (hooki), dopięcie w Fazie 15.
