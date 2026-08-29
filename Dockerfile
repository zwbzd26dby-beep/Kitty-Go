# Build stage
FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/kitty-repl ./cmd/kitty-repl \
 && CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/kitty-cli ./cmd/kitty-cli \
 && CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/kitty-api ./cmd/kitty-api \
 && CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/kitty-gui ./cmd/kitty-gui

# Runtime stage (alpine provides sh for entrypoint.sh)
FROM alpine:3.20
RUN apk add --no-cache ca-certificates && adduser -D -u 10001 kitty
COPY --from=build /out/kitty-* /usr/local/bin/
COPY --from=build /src/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
EXPOSE 8080
ENV KITTY_HOST=0.0.0.0
ENV KITTY_PORT=8080
USER kitty
ENTRYPOINT ["/entrypoint.sh"]
CMD ["api"]