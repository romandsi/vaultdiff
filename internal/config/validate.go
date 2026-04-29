package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// ValidationError holds a list of validation issues found in a Config.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("config validation failed:\n  - %s", strings.Join(e.Issues, "\n  - "))
}

// Validate checks that a loaded Config is semantically valid.
// It returns a *ValidationError listing every problem found, or nil.
func Validate(cfg *Config) error {
	var issues []string

	if len(cfg.Environments) < 2 {
		issues = append(issues, "at least two environments are required")
	}

	seenNames := make(map[string]bool)
	for i, env := range cfg.Environments {
		prefix := fmt.Sprintf("environment[%d]", i)

		if strings.TrimSpace(env.Name) == "" {
			issues = append(issues, prefix+": name must not be empty")
		} else if seenNames[env.Name] {
			issues = append(issues, fmt.Sprintf("%s: duplicate environment name %q", prefix, env.Name))
		} else {
			seenNames[env.Name] = true
		}

		if strings.TrimSpace(env.Address) == "" {
			issues = append(issues, prefix+": address must not be empty")
		} else if _, err := url.ParseRequestURI(env.Address); err != nil {
			issues = append(issues, fmt.Sprintf("%s: address %q is not a valid URL", prefix, env.Address))
		}

		if strings.TrimSpace(env.Token) == "" {
			issues = append(issues, prefix+": token must not be empty")
		}
	}

	if len(cfg.Paths) == 0 {
		issues = append(issues, "at least one secret path is required")
	}

	for i, p := range cfg.Paths {
		if strings.TrimSpace(p) == "" {
			issues = append(issues, fmt.Sprintf("paths[%d]: path must not be empty", i))
		}
	}

	if len(issues) > 0 {
		return &ValidationError{Issues: issues}
	}
	return nil
}

// ErrInvalidConfig is a sentinel for callers that only need to detect
// validation failures without inspecting individual issues.
var ErrInvalidConfig = errors.New("invalid config")
