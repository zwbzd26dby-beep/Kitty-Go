package agent

import "github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"

// historyAdapter provides history management for an agent backed by a
// types.Conversation. In later phases this is replaced by the full Memory
// system (Phase 12).
type historyAdapter struct {
	conversation *types.Conversation
}

// appendUserMessage records a user message in the history.
func (h *historyAdapter) appendUserMessage(msg types.Message) {
	if h.conversation != nil {
		h.conversation.AddMessage(msg)
	}
}

// appendAssistantResponse records an assistant response in the history.
func (h *historyAdapter) appendAssistantResponse(resp types.Response) {
	if h.conversation != nil {
		h.conversation.AddResponse(resp)
	}
}

// history returns the current conversation history (a defensive copy).
func (h *historyAdapter) history() []types.Turn {
	if h.conversation == nil {
		return nil
	}
	return h.conversation.GetHistory()
}

// clear wipes the conversation history.
func (h *historyAdapter) clear() {
	if h.conversation != nil {
		h.conversation.Clear()
	}
}
