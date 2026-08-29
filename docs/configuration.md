# Konfiguracja Kitty-Go

## Pliki konfiguracyjne

| Plik                      | Rola                                   |
|---------------------------|----------------------------------------|
| `configs/default.yaml`    | domyślne wartości startowe (bez sekretów) |
| `configs/user.yaml`       | nadpisania użytkownika                 |
| `configs/secrets.yaml`    | klucze API (NIGDY nie commitować)      |

Przykłady: `configs/user.yaml.example`, `configs/secrets.yaml.example`.

## Ładowanie (internal/config)

1. `Load(path)` czyta YAML (gopkg.in/yaml.v3).
2. `ParseUser` wyłuskuje profil użytkownika.
3. `Validate` sprawdza spójność (np. niepuste podstawowe pola).
4. Zmienne środowiskowe mają pierwszeństwo nad plikami.

## Klucze API (internal/security)

Klucze są trzymane w pamięci managera (`security.Manager`) i **nigdy nie są
logowane ani zapisywane przez pakiet**. Rozwiązywanie klucza:

1. jawnie ustawiony klucz (`SetAPIKey`),
2. zmienna środowiskowa z mapy fallbacków,
3. błąd braku klucza.

Domyślna mapa provider → env:

| Provider    | Zmienna           |
|-------------|-------------------|
| openai      | `OPENAI_API_KEY`  |
| kimi        | `KIMI_API_KEY`    |
| deepseek    | `DEEPSEEK_API_KEY`|
| openrouter  | `OPENROUTER_API_KEY` |
| ollama      | `OLLAMA_API_KEY`  |
| opencode    | `OPENCODE_API_KEY`|

## Zmienne środowiskowe aplikacji

| Zmienna           | Domyślne     | Opis                        |
|-------------------|--------------|-----------------------------|
| `OPENCODE_API_KEY`| (brak)       | klucz dostępu OpenCode Zen  |
| `KITTY_HOST`      | `127.0.0.1`  | host API                 |
| `KITTY_PORT`      | `8080`       | port API                    |
| `KITTY_API_TOKEN` | (brak)       | token Bearer API            |
| `KITTY_GUI_ADDR`  | `127.0.0.1:8081` | adres GUI              |

## Limity / Budżet

- `internal/limit`: rate limiter (na minutę) + quota (dzienny, miesięczny).
- `internal/budget`: dzienny i miesięczny limit kosztów; przekroczenie blokuje.
- `internal/cost`: kalkulacja i śledzenie kosztów według modelu.