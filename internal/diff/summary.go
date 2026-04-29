package diff

import (
	"fmt"
	"io"
	"sort"
)

// Summary holds aggregated statistics about a diff result.
type Summary struct {
	Total   int
	Added   int
	Removed int
	Changed int
	Unchanged int
	ByPath  map[string]*PathSummary
}

// PathSummary holds per-path diff statistics.
type PathSummary struct {
	Added   int
	Removed int
	Changed int
}

// Summarize computes a Summary from a slice of DiffEntry values.
func Summarize(entries []DiffEntry) Summary {
	s := Summary{
		ByPath: make(map[string]*PathSummary),
	}

	for _, e := range entries {
		s.Total++

		if _, ok := s.ByPath[e.Path]; !ok {
			s.ByPath[e.Path] = &PathSummary{}
		}
		ps := s.ByPath[e.Path]

		switch e.Status {
		case StatusAdded:
			s.Added++
			ps.Added++
		case StatusRemoved:
			s.Removed++
			ps.Removed++
		case StatusChanged:
			s.Changed++
			ps.Changed++
		case StatusUnchanged:
			s.Unchanged++
		}
	}

	return s
}

// WriteSummary writes a human-readable summary to w.
func WriteSummary(w io.Writer, s Summary) {
	fmt.Fprintf(w, "Summary: %d total keys ", s.Total)
	fmt.Fprintf(w, "(%d added, %d removed, %d changed, %d unchanged)\n",
		s.Added, s.Removed, s.Changed, s.Unchanged)

	if len(s.ByPath) == 0 {
		return
	}

	paths := make([]string, 0, len(s.ByPath))
	for p := range s.ByPath {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	fmt.Fprintln(w, "\nPer-path breakdown:")
	for _, p := range paths {
		ps := s.ByPath[p]
		fmt.Fprintf(w, "  %s: +%d -%d ~%d\n", p, ps.Added, ps.Removed, ps.Changed)
	}
}
