# Deployment

## Obraz Docker

Multi-stage `Dockerfile` buduje cztery binarie (`kitty-repl`, `kitty-cli`,
`kitty-api`, `kitty-gui`) i pakuje je do lekkiego obrazu Alpine z `entrypoint.sh`.

```sh
docker build -t kitty-go:latest .
```

### Tryby entrypointa

```sh
docker run --rm kitty-go:latest api     # serwer HTTP (domyślny)
docker run --rm kitty-go:latest gui     # web chat
docker run --rm kitty-go:latest repl    # REPL (wymaga TTY)
docker run --rm kitty-go:latest cli "hello"
```

## docker-compose

```sh
export KITTY_API_TOKEN=$(openssl rand -hex 24)
export OPENCODE_API_KEY=sk-...
docker compose up -d --build
```

Usługi:
- `kitty-api` na `:8080` (z healthcheckiem),
- `kitty-gui` na `:8081`.

## Wdrożenie produkcyjne

- Klucze dostarczać przez sekrety CI / vault, nie przez zmienne w obrazie.
- `OPENCODE_API_KEY` tylko w środowisku runtime.
- Reverse proxy (np. Caddy/nginx) przed `kitty-api` dla TLS.
- Ustawić `KITTY_API_TOKEN` — inaczej `/agent/respond` jest otwarte.
- Monitorować `/health` i `/metrics`.

## Profiling

Podczas lokalnego testowania można wskazać pprof:

```sh
curl -o heap.out http://127.0.0.1:8080/debug/pprof/heap
go tool pprof heap.out
```