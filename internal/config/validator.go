package config

// Validate performs non-fatal sanity checks on a Config and returns a list of
// informational/warning messages.
func Validate(c *Config) []Message {
	var msgs []Message
	if c.Provider == "" {
		msgs = append(msgs, Message{Level: "warning", Message: "provider is empty; defaulting to mock at call sites"})
	}
	if c.Timeout <= 0 {
		msgs = append(msgs, Message{Level: "warning", Message: "timeout must be positive; using call-site default"})
	}
	switch c.REPL.SecurityLevel {
	case "", "standard", "high", "low":
	default:
		msgs = append(msgs, Message{Level: "warning", Message: "unknown security_level; expected low|standard|high"})
	}
	if c.Budget.Daily < 0 || c.Budget.Monthly < 0 {
		msgs = append(msgs, Message{Level: "warning", Message: "budget values must be non-negative"})
	}
	return msgs
}
