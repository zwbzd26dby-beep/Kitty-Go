# Rozwój Kitty-Go

## Wymagania

- Go 1.26+
- Bez obowiązkowej zasady "stdlib only" — wybrane zewnętrzne zależności są
  dozwolone (np. YAML config, test framework), jeśli śrwiadomie zaakceptowane.

## Testy

```bash
make test          # cały zestaw
make unit          # unit
make integration   # integracyjne / cli
```

## Zasady

- Każda faza kończy się zestawem testów w `tests/<obszar>/` (+ jednostkowe w
  pakietach) zanim przejdziemy dalej — "suite zawsze green".
- Pakiety w `internal/` komunikują się przez interfejsy, nie przez
  konkretne struktury.
- Sekrety nigdy nie trafiają do Gita (`configs/secrets.yaml` w `.gitignore`).
- Commit `feat(phaseN): implement <nazwa>`.

## Struktura katalogów

Pełny szablon: patrz `docs/architecture.md` oraz Master Architecture §3.
