#!/bin/sh
# entrypoint.sh starts the selected Kitty command.
# Usage: entrypoint.sh [repl|cli|api|gui] [args...]
set -eu

mode="${1:-api}"
shift 2>/dev/null || true

case "$mode" in
  repl)   exec /kitty-repl "$@" ;;
  cli)    exec /kitty-cli "$@" ;;
  api)    exec /kitty-api "$@" ;;
  gui)    exec /kitty-gui "$@" ;;
  *)      echo "unknown mode: $mode" >&2; exit 2 ;;
esac