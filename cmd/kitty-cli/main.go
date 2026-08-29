// Command kitty-cli runs a one-shot Kitty request using the deterministic
// mock provider.
package main

import (
	"context"
	"os"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/interface/cli"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/llm"
)

func main() {
	client := llm.NewMockClient()
	os.Exit(cli.Run(context.Background(), client, os.Args[1:], os.Stdout))
}
