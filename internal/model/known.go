package model

// RegisterKnown seeds a Registry with the statically-known models across all
// providers, including the OpenCode Zen catalog.
func RegisterKnown(r ModelRegistry) error {
	known := []Model{
		// OpenCode Zen (functional requirement).
		{IDs: IDs{Provider: "opencode", Model: "big-pickle"}, Description: "OpenCode Zen coding/reasoning model", Capabilities: Capabilities{CodeGeneration: true, CodeReasoning: true, Tools: true, Streaming: true}, ContextWindow: 128000, MaxOutputTokens: 8192, Availability: Available},
		{IDs: IDs{Provider: "opencode", Model: "mimo-v2-pro-free"}, Description: "OpenCode Zen free general model", Capabilities: Capabilities{Streaming: true}, ContextWindow: 32000, MaxOutputTokens: 4096, Availability: Available},
		{IDs: IDs{Provider: "opencode", Model: "minimax-m2.5-free"}, Description: "OpenCode Zen free MiniMax model", Capabilities: Capabilities{Streaming: true}, ContextWindow: 32000, MaxOutputTokens: 4096, Availability: Available},

		// OpenAI.
		{IDs: IDs{Provider: "openai", Model: "gpt-4o"}, Description: "OpenAI GPT-4o", Capabilities: Capabilities{CodeGeneration: true, Vision: true, Tools: true, Streaming: true}, Pricing: Pricing{InputCostPerMillion: 2.5, OutputCostPerMillion: 10}, ContextWindow: 128000, MaxOutputTokens: 4096, Availability: Available},
		{IDs: IDs{Provider: "openai", Model: "gpt-4o-mini"}, Description: "OpenAI GPT-4o mini", Capabilities: Capabilities{CodeGeneration: true, Vision: true, Tools: true, Streaming: true}, Pricing: Pricing{InputCostPerMillion: 0.15, OutputCostPerMillion: 0.6}, ContextWindow: 128000, MaxOutputTokens: 4096, Availability: Available},

		// Kimi (Moonshot).
		{IDs: IDs{Provider: "kimi", Model: "kimi-k2"}, Description: "Kimi K2", Capabilities: Capabilities{CodeGeneration: true, Tools: true, Streaming: true}, ContextWindow: 128000, MaxOutputTokens: 8192, Availability: Available},

		// DeepSeek.
		{IDs: IDs{Provider: "deepseek", Model: "deepseek-chat"}, Description: "DeepSeek chat", Capabilities: Capabilities{CodeGeneration: true, Tools: true, Streaming: true}, ContextWindow: 64000, MaxOutputTokens: 8192, Availability: Available},

		// Ollama (local).
		{IDs: IDs{Provider: "ollama", Model: "llama3"}, Description: "Ollama llama3 local", Capabilities: Capabilities{CodeGeneration: true, Streaming: true}, ContextWindow: 8192, MaxOutputTokens: 2048, Availability: Available},
	}
	for _, m := range known {
		if err := r.Register(m); err != nil {
			return err
		}
	}
	return nil
}
