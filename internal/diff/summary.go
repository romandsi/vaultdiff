package diff

import (
	"fmt"
	"io"
)

// Summary holds aggregated counts of diff results.
type Summary struct {
	Added   int
	Removed int
	Changed int
	Total   int
}

// Summarize computes a Summary from a slice of Diffs.
func Summarize(diffs []Diff) Summary {
	s := Summary{Total: len(diffs)}
	for _, d := range diffs {
		switch d.Type {
		case DiffAdded:
			s.Added++
		case DiffRemoved:
			s.Removed++
		case DiffChanged:
			s.Changed++
		}
	}
	return s
}

// WriteSummary writes a human-readable summary line to w.
func WriteSummary(w io.Writer, s Summary) {
	if s.Total == 0 {
		fmt.Fprintln(w, "Summary: no differences found.")
		return
	}
	fmt.Fprintf(w, "Summary: %d difference(s) — ", s.Total)
	parts := []string{}
	if s.Added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", s.Added))
	}
	if s.Removed > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", s.Removed))
	}
	if s.Changed > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", s.Changed))
	}
	for i, p := range parts {
		if i > 0 {
			fmt.Fprint(w, ", ")
		}
		fmt.Fprint(w, p)
	}
	fmt.Fprintln(w)
}
