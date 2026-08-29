// Package types defines the core immutable message types used across Kitty.
package types

import (
	"fmt"
	"strings"
)

// Message represents a single user message. It is immutable after creation.
type Message struct {
	content string
}

// NewMessage creates a Message, validating that content is a non-blank string.
func NewMessage(content string) (Message, error) {
	if strings.TrimSpace(content) == "" {
		return Message{}, fmt.Errorf("message content must not be empty")
	}
	return Message{content: content}, nil
}

// Content returns the message content.
func (m Message) Content() string { return m.content }

// String implements fmt.Stringer.
func (m Message) String() string { return m.content }
