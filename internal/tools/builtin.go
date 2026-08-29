package tools

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Calculator computes basic arithmetic safely (no eval).
type Calculator struct{}

// NewCalculator instantiates the calculator tool.
func NewCalculator() *Calculator { return &Calculator{} }

// Name reports "calculator".
func (*Calculator) Name() string { return "calculator" }

// Description explains the tool.
func (*Calculator) Description() string { return "evaluate a basic arithmetic expression" }

// Params declares the expression argument.
func (*Calculator) Params() []Param {
	return []Param{{Name: "expr", Type: "string", Description: "expression", Required: true}}
}

// Execute parses and evaluates + - * / ( ) safely.
func (*Calculator) Execute(_ context.Context, args map[string]string) (string, error) {
	val, err := evalExpr(strings.TrimSpace(args["expr"]))
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(val, 'g', 10, 64), nil
}

// HashTool computes a SHA-256 digest.
type HashTool struct{}

// NewHashTool instantiates the hash tool.
func NewHashTool() *HashTool { return &HashTool{} }

// Name reports "sha256".
func (*HashTool) Name() string { return "sha256" }

// Description explains the tool.
func (*HashTool) Description() string { return "compute a SHA-256 hex digest" }

// Params declares the input argument.
func (*HashTool) Params() []Param {
	return []Param{{Name: "input", Type: "string", Description: "text to hash", Required: true}}
}

// Execute hashes the input.
func (*HashTool) Execute(_ context.Context, args map[string]string) (string, error) {
	sum := sha256.Sum256([]byte(args["input"]))
	return fmt.Sprintf("%x", sum), nil
}

// ClockTool returns the current UTC time.
type ClockTool struct{}

// NewClockTool instantiates the clock tool.
func NewClockTool() *ClockTool { return &ClockTool{} }

// Name reports "clock".
func (*ClockTool) Name() string { return "clock" }

// Description explains the tool.
func (*ClockTool) Description() string { return "return the current UTC timestamp" }

// Params declares no arguments.
func (*ClockTool) Params() []Param { return nil }

// Execute returns the current time.
func (*ClockTool) Execute(_ context.Context, _ map[string]string) (string, error) {
	return time.Now().UTC().Format(time.RFC3339), nil
}

// RegisterBuiltins registers the safe built-in tools.
func RegisterBuiltins(r *Registry) {
	r.Register(NewCalculator())
	r.Register(NewHashTool())
	r.Register(NewClockTool())
}
