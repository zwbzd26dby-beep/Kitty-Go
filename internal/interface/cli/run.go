// Package cli implements the one-shot command-line interaction, decoupled from
// the process entry point so it can be tested from tests/cli.
package cli

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/llm"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"
)

// Run executes a one-shot request and writes the reply to w.
// It returns a process-style exit code: 0 on success, 2 on usage/validation
// errors, 1 on provider errors.
func Run(ctx context.Context, client llm.Client, args []string, w io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(w, "Użycie: kitty-cli <wiadomość>")
		return 2
	}
	input := strings.Join(args, " ")
	msg, err := types.NewMessage(input)
	if err != nil {
		fmt.Fprintln(w, err)
		return 2
	}
	resp, err := client.Generate(ctx, msg, nil)
	if err != nil {
		fmt.Fprintln(w, err)
		return 1
	}
	fmt.Fprintf(w, "Kitty: %s\n", resp.Content())
	return 0
}
