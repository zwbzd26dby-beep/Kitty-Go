# Kitty-Go HTTP API

Kitty-Go exposes a persistent HTTP API bound to the shared Core via
`internal/interface/api` and `cmd/kitty-api`.

## Uruchomienie

```sh
go run ./cmd/kitty-api
```

Konfiguracja przez zmienne środowiskowe:

| Zmienna          | Domyślne     | Opis                            |
|------------------|--------------|---------------------------------|
| `KITTY_HOST`     | `127.0.0.1`  | Adres nasłuchu                  |
| `KITTY_PORT`     | `8080`       | Port                            |
| `KITTY_API_TOKEN`| (brak)       | Token Bearer; puste = bez auth  |

## Endpointy

### `GET /health`

| Kod | Opis |
|-----|------|
| 200 | usługa żyje |

```json
{"status":"ok"}
```

### `GET /metrics`

Zwraca zrzut liczników i gauge'ów obserwowalności.

```json
{"counters":{"tasks_total":3},"gauges":{},"uptime_seconds":42.5}
```

### `POST /agent/respond`

Przetwarza wiadomość przez agenta.

Body:

```json
{"message":"Napisz test jednostkowy"}
```

Odpowiedź (200):

```json
{"content":"...","intent":"chat"}
```

| Kod | Opis |
|-----|------|
| 200 | sukces                        |
| 400 | błąd walidacji (pusta wiadomość, zły JSON) |
| 401 | brak / błędny token Bearer     |
| 405 | nieobsługiwana metoda          |
| 500 | błąd przetwarzania             |

## Autoryzacja

Jeśli ustawiono `KITTY_API_TOKEN`, żądania do `/agent/respond` wymagają:

```
Authorization: Bearer <token>
```

## Profiling (pprof)

Endpoints profilingowe dostępne pod `/debug/pprof/…`:

```sh
go tool pprof http://127.0.0.1:8080/debug/pprof/heap
```