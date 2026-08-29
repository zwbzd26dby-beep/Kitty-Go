package clitest

import (
	"bytes"
	"context"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/interface/cli"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/llm"
)

func TestCLIOneShot(t *testing.T) {
	var out bytes.Buffer
	code := cli.Run(context.Background(), llm.NewMockClient(), []string{"Hej Kitty"}, &out)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if out.String() != "Kitty: Kitty LLM mock działa.\n" {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestCLINoArgs(t *testing.T) {
	var out bytes.Buffer
	if code := cli.Run(context.Background(), llm.NewMockClient(), nil, &out); code != 2 {
		t.Fatalf("expected exit 2 with no args, got %d", code)
	}
}

func TestCLIBlankMessage(t *testing.T) {
	var out bytes.Buffer
	if code := cli.Run(context.Background(), llm.NewMockClient(), []string{"   "}, &out); code != 2 {
		t.Fatalf("expected exit 2 for blank message, got %d", code)
	}
}
