# HANDOFF — Przejście Kitty: Python → Go

**Status:** decyzja podjęta — Go jest docelowym językiem/architekturą Kitty
**Data:** 2026-08-29
**Źródło architektury:** dokument "KITTY — MASTER ARCHITECTURE" (pełna specyfikacja, 20 faz)
**Źródło stanu obecnego:** HANDOFF.md (Python, stdlib-only, 386/386 testów, 20 commitów)
**Adresat:** OpenCode (implementacja na Fedorze)

Dokument opisuje decyzję o przejściu oraz zasady, według których rozwijana jest
implementacja Go w tym repozytorium.

> **AKTUALIZACJA (2026-08-29):** wszystkie fazy 0–20 zostały zrealizowane.
> Pełny stan implementacji Go — patrz **§8. Stan realizacji faz 0–20** poniżej.

---

## 1. Decyzja

Dotychczasowa implementacja w Pythonie (`src/kitty/...`, wyłącznie stdlib)
**nie jest** dalej rozwijana jako architektura docelowa. Obowiązuje struktura,
nazewnictwo i kontrakty z dokumentu **MASTER ARCHITECTURE** (Go), wraz z pełną
roadmapą 20 faz opisaną w tamtym dokumencie.

Python **nie jest kasowany** — pełni rolę referencyjną (patrz §3).

## 2. Dlaczego

- Master Architecture jest kompletną specyfikacją: interfejsy, warstwy,
  przepływy (request, routing, koszty, fallbacki, błędy), pełna struktura
  katalogów i 20-fazowa roadmapa gotowa do wdrożenia bez dalszego projektowania.
- Realna implementacja Python rozjechała się ze specyfikacją na poziomie
  fundamentalnym: inny język, inny podział modułów (w Go Decision Engine +
  Model Router to osobne moduły, w Pythonie są scalone), inne ograniczenia
  (Python wymuszał "stdlib only", co powodowało konflikty, np. nieśledzony
  `config/user.py` importował PyYAML).
- Równoległy rozwój obu ścieżek pogłębiałby rozjazd.

## 3. Co się dzieje z istniejącym kodem Python

Traktujemy Python jako **specyfikację behawioralną**, nie bazę do portowania 1:1:

1. **Zamrożenie gałęzi.** Repozytorium Python (`Kitty`, `main`) zostaje otagowane
   jako punkt odniesienia, np. `archive/python-v1`.
2. **386 testów jako źródło wymagań.** Testy `tests/test_phase1..15.py` +
   `test_repl_integration.py` + `test_providers_opencode.py` opisują zachowania,
   które system Go musi odtworzyć (logika wyboru modelu, fallbacki, budżety/
   limity, RAG).
3. **Provider OpenCode Zen (`big-pickle`) jest wymaganiem funkcjonalnym** mimo
   braku w Master Architecture — musi zostać dodany jako implementacja
   `Provider` w Fazie 2.
4. **3 niewypchnięte commity Python** (`c2e0d0a`, `d26ed81`, `8cba66a`) —
   wypchnąć na `origin/main` przed zamrożeniem gałęzi.
5. **Nieśledzone pliki Python** (`src/kitty/config/` z PyYAML, puste
   `interfaces/`, `utils/`, `KITTY_PROJECT_DOCUMENTATION.md`) — nie przenosimy
   ich do Go.

## 4. Struktura nowego repozytorium (Go)

Zgodnie z Master Architecture §3 (pełny szablon w `docs/architecture.md`):
`cmd/{kitty-cli,kitty-repl,kitty-api,kitty-gui}/`, `internal/{agent,orchestrator,
decision,router,execution,provider,model,cost,limit,budget,compute,memory,tools,
security,config,observability,interface}/`, `pkg/{types,llm,utils}/`,
`scentralizowane tests/{unit,integration,provider,router,cost,limit,memory,tools,
security,compute,distributed,cli,e2e}/`, `scripts/`, `docs/`, `configs/`.

### Lokalizacja testów — decyzja

Scentralizowane `tests/` zgodnie z MASTER ARCHITECTURE §3. Nie trzymamy testów
w pakietach obok kodu (mimo że to idiomatyczny Go) — priorytetem jest zgodność
ze specyfikacją faz 0–20, która odwołuje się do ścieżek `tests/<obszar>/*_test.go`.

## 5. Kolejność faz

Roadmapa 1:1 z Master Architecture §28 oraz specyfikacją `docs/KITTY_GO_PHASES.md`.
Kolejność wiążąca: Faza 0 → 1 → 2 → … → 20.

Zależności do pilnowania:
- Faza 6 (Model Router) potrzebuje Fazy 5 (Registry) i częściowo Fazy 7
  (Limit/Budget — na start stuby zwracające "brak limitu").
- Fazy 10 (Compute Router) i 16 (Decision Engine) mogą być rozwijane równolegle,
  obie zasilają Fazę 17 (Orchestrator pełny).
- Fazę 15 (Observability) wpinamy przyrostowo od Fazy 7 (hooki), dopięcie w Fazie 15.

## 6. Decyzje otwarte do rozstrzygnięcia

1. Nowe repo (`Kitty-Go`) — **rozstrzygnięte: tak, osobne repo.**
2. Reguła "stdlib only" — **rozstrzygnięte: dopuszczamy wybrane zależności**
   (np. `gopkg.in/yaml.v3` do configu, do testów).
3. Numeracja commitów — wzorzec `feat(phaseN): implement <nazwa>`.
4. `docs/AGENTS.md` — przenosimy z Python z aktualizacją zasad pod Go.

## 7. Workflow

Planowanie/architektura → Kimi/DeepSeek/ChatGPT.
Implementacja/testy/commity → OpenCode, lokalnie na Fedorze.
GitHub — źródło prawdy.

Ten dokument jest punktem startowym dla OpenCode do realizacji faz w Go.

---

## 8. Stan realizacji faz 0–20 (AKTUALIZACJA 2026-08-29)

**Status:** wszystkie fazy zrealizowane. Repo `Kitty-Go`, `main` `8a90b39`,
18 commitów, **179 testów** (w tym `-race`), pełne `go build`/`go vet` czyste.

### 8.1 Pakiety

```
cmd/kitty-{cli,repl,api,gui}/   internal/{agent,budget,compute,config,cost,
decision,execution,interface/{api,cli,gui,repl},limit,memory,model,
observability,orchestrator,provider/{openai,kimi,deepseek,openrouter,ollama,
opencode,mock},router/{model,compute},security,tools}/   pkg/{types,llm,utils}/
tests/{unit,integration,provider,router,cost,limit,memory,tools,security,
compute,distributed,cli,e2e,observability}/
```

### 8.2 Fazy — implementacja i commity

| Faza | Kommit | Co powstało |
|------|--------|-------------|
| 0 | `69b67f0` | fundament: `pkg/types`, `pkg/llm` mocks, CLI, REPL, docs |
| 1 | `a653c0f` | `internal/{agent,orchestrator,execution}`: pętla Agent→Orch→Exec |
| 2 | `26f7c6b` | system providerów (OpenAI/kimi/DeepSeek/OpenRouter/Ollama/OpenCode) |
| 3 | `846f960` | realne LLM (retry + backoff, `HTTPError`), config/secrets |
| 4 | `846f960` | `internal/config` (YAML), `internal/security` (klucze API) |
| 5 | `d0d6344` | `internal/model`: registry, capabilities, pricing, knowns |
| 6 | `d0d6344` | `internal/router/model`: `ModelRouter` + preferencja big-pickle |
| 7 | `30088f0` | `internal/{cost,limit,budget}`: koszty, limity, budżet |
| 8 | `038a513` | `internal/compute`: fabric (devices, health, discovery, worker) |
| 9 | `038a513` | `internal/execution/remote.go`: `DistributedExecutor` + transport |
| 10 | `038a513` | `internal/router/compute`: `ComputeRouter` (scoring, scheduler) |
| 11 | `19d1805` | pełny `Execution Layer`: `Runner`, retry, fallback (§36) |
| 12 | `19d1805` | `internal/memory`: STM/LTM/semantic/RAG |
| 13 | `c3d30be` | `internal/tools`: registry, allowlist, safe builtiny |
| 14 | `021b449` | sandbox, redactor, prompt guard |
| 15 | `22dc8a7` | `internal/observability`: logger, metrics, trace, usage, cost, audit |
| 16 | `b6f16a8` | `internal/decision`: `DecisionEngine` (reguły + wagi + scoring) |
| 17 | `9fa86f8` | pełny orchestrator: Decision→Model→Compute→Exec + budget/limit/cost |
| 18 | `8d0659f` | CLI/REPL (13 komend)/API (Bearer) na wspólnym Core |
| 19 | `8d0659f` | `internal/interface/gui`: web chat (`cmd/kitty-gui`) |
| 20 | `007125f` | Dockerfile, docker-compose, CI (2 wersje Go), docs, pprof |

### 8.3 Rozstrzygnięcia techniczne

- **Zależności:** dozwolone wybrane zewnętrzne — w użyciu tylko
  `gopkg.in/yaml.v3` (config). Go `1.26.7`.
- **Testy:** scentralizowane w `tests/<obszar>/` (decyzja z §4).
- **Routing big-pickle:** reguła przeniesiona z Fazy 2 do `router/model`
  (Faza 6); `opencode.Route` usunięty.
- **Retry:** `utils.RetryWith` (backoff śledzony); użyty przez providery
  (`HTTPError` staje się przejściowy dla 429/5xx) i `limit.RetryWithBackoff`.
- **Determinizm scheduler:** tie-break po `Device.ID` w `sortByScore`
  (flaky round-robin naprawiony).
- **Observability:** `observability.New(nil)` → `io.Discard` (brak paniki).
- **REPL `/model`:** dopasowanie po samym ID modelu (przeszukanie `List()`).

### 8.4 API / GUI

- `POST /agent/respond` (Bearer `KITTY_API_TOKEN`), `GET /health`,
  `GET /metrics`, pprof under `/debug/pprof/`.
- Zmienne: `KITTY_HOST`, `KITTY_PORT`, `KITTY_GUI_ADDR`.
- GUI: web chat w `cmd/kitty-gui`, `internal/interface/gui`.
- Docs: `docs/{api,configuration,deployment,architecture,development}.md`.

### 8.5 Metryka testów (per obszar)

`unit:34, provider:15, router:21, cli:21, integration:12, e2e:13,
security:13, cost:10, tools:9, compute:8, observability:7, limit:7,
distributed:5, memory:4` — razem **179**.

### 8.6 Co dalej (przyrostowe)

1. Realny provider OpenAI/OpenCode zamiast mocków (smoke e2e live).
2. Keys do configu: sekcje `configs/secrets.yaml` + `security.Manager`.
3. ComputeRegistry z faktyczną telemetrią urządzeń (`/metrics` per device).
4. RAG zasilany prawdziwym embeddingiem zamiast `HashEmbedder`.
5. Testy live `tests/integration/live_provider_test.go` (włączane przez
   `OPENCODE_API_KEY`).

---
