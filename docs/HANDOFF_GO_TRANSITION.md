# HANDOFF — Przejście Kitty: Python → Go

**Status:** decyzja podjęta — Go jest docelowym językiem/architekturą Kitty
**Data:** 2026-08-29
**Źródło architektury:** dokument "KITTY — MASTER ARCHITECTURE" (pełna specyfikacja, 20 faz)
**Źródło stanu obecnego:** HANDOFF.md (Python, stdlib-only, 386/386 testów, 20 commitów)
**Adresat:** OpenCode (implementacja na Fedorze)

Dokument opisuje decyzję o przejściu oraz zasady, według których rozwijana jest
implementacja Go w tym repozytorium.

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
