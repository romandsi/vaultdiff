package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/vaultdiff/internal/diff"
)

// Format represents the output format for diff results.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatMarkdown Format = "markdown"
)

// ParseFormat parses a string into a Format, returning an error if unknown.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "text", "":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	case "markdown", "md":
		return FormatMarkdown, nil
	default:
		return "", fmt.Errorf("unknown output format %q: must be one of text, json, markdown", s)
	}
}

// Write writes the diff results to w in the specified format.
func Write(w io.Writer, results []diff.Result, envA, envB string, format Format) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, results, envA, envB)
	case FormatMarkdown:
		return writeMarkdown(w, results, envA, envB)
	default:
		return writeText(w, results, envA, envB)
	}
}

func writeText(w io.Writer, results []diff.Result, envA, envB string) error {
	if len(results) == 0 {
		_, err := fmt.Fprintf(w, "No differences found between %s and %s.\n", envA, envB)
		return err
	}
	for _, r := range results {
		_, err := fmt.Fprintf(w, "[%s] %s: %s\n", r.Path, r.Key, r.Kind)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, results []diff.Result, envA, envB string) error {
	payload := map[string]interface{}{
		"env_a":   envA,
		"env_b":   envB,
		"results": results,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}

func writeMarkdown(w io.Writer, results []diff.Result, envA, envB string) error {
	_, err := fmt.Fprintf(w, "## Vault Diff: `%s` vs `%s`\n\n", envA, envB)
	if err != nil {
		return err
	}
	if len(results) == 0 {
		_, err = fmt.Fprintln(w, "_No differences found._")
		return err
	}
	_, err = fmt.Fprintln(w, "| Path | Key | Change |")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, "|------|-----|--------|")
	if err != nil {
		return err
	}
	for _, r := range results {
		_, err = fmt.Fprintf(w, "| `%s` | `%s` | %s |\n", r.Path, r.Key, r.Kind)
		if err != nil {
			return err
		}
	}
	return nil
}
