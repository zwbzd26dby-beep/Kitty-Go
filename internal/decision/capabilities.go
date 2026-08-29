package decision

// requirementsFor maps a task kind to model-capability requirements.
func requirementsFor(kind string) []Requirement {
	switch kind {
	case "code":
		return []Requirement{
			{Name: "code-gen", Capability: "CodeGeneration", Description: "generating or editing code"},
			{Name: "code-reasoning", Capability: "CodeReasoning", Description: "reasoning about code"},
			{Name: "tools", Capability: "Tools", Description: "driving tool calls"},
		}
	case "math":
		return []Requirement{
			{Name: "reasoning", Capability: "CodeReasoning", Description: "multi-step derivation"},
		}
	case "vision":
		return []Requirement{
			{Name: "vision", Capability: "Vision", Description: "understanding images"},
		}
	case "chat":
		return []Requirement{
			{Name: "tools", Capability: "Tools", Description: "optional tool use"},
		}
	default:
		return []Requirement{{Name: "chat", Description: "general dialogue"}}
	}
}
