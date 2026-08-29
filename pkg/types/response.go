package types

import "strings"

// Response represents an assistant's reply. Immutable after creation.
type Response struct {
	content string
}

// NewResponse creates a Response, validating that content is non-blank.
func NewResponse(content string) (Response, error) {
	if strings.TrimSpace(content) == "" {
		return Response{}, errEmptyContent("response")
	}
	return Response{content: content}, nil
}

// Content returns the response content.
func (r Response) Content() string { return r.content }

// String implements fmt.Stringer.
func (r Response) String() string { return r.content }
