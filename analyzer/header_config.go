package analyzer

import (
	"fmt"
	"slices"
)

// StyleConfig configures optional style rules without enabling them.
type StyleConfig struct {
	Header HeaderConfig `json:"header"`
}

// HeaderConfig is shared by header formatting and required metadata checks.
// Nil Methods preserves the built-in set; an empty list disables metadata
// recognition. Parallel remains a formatting header element.
type HeaderConfig struct {
	Methods           []string   `json:"methods"`
	Required          [][]string `json:"required"`
	IncludeBeforeEach *bool      `json:"include-before-each"`
}

var defaultHeaderMethods = []string{"Title", "Titlef", "Tags", "Description", "Descriptionf", "ID", "Labels", "Links", "Owner", "Severity", "Epic", "Feature", "Story", "Parameters", "Flaky", "Muted", "Known"}

func (h HeaderConfig) methods() []string {
	if h.Methods == nil {
		return defaultHeaderMethods
	}
	return h.Methods
}

func (h HeaderConfig) validate() error {
	seen := map[string]bool{}
	for _, name := range h.methods() {
		if !slices.Contains(defaultHeaderMethods, name) {
			return fmt.Errorf("testolint: unknown header method %q", name)
		}
		if seen[name] {
			return fmt.Errorf("testolint: duplicate header method %q", name)
		}
		seen[name] = true
	}
	for _, group := range h.Required {
		if len(group) == 0 {
			return fmt.Errorf("testolint: required header group must not be empty")
		}
		members := map[string]bool{}
		for _, name := range group {
			if !seen[name] {
				return fmt.Errorf("testolint: required header method %q is not recognized", name)
			}
			if members[name] {
				return fmt.Errorf("testolint: duplicate required header method %q", name)
			}
			members[name] = true
		}
	}
	return nil
}

func (h HeaderConfig) clone() HeaderConfig {
	h.Methods = slices.Clone(h.Methods)
	h.Required = slices.Clone(h.Required)
	for i := range h.Required {
		h.Required[i] = slices.Clone(h.Required[i])
	}
	if h.IncludeBeforeEach != nil {
		value := *h.IncludeBeforeEach
		h.IncludeBeforeEach = &value
	}
	return h
}
