package diff

import (
	"fmt"
	"io"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

// ReportOptions controls output formatting.
type ReportOptions struct {
	Color   bool
	MaskValues bool
}

// WriteReport writes a human-readable diff report for a slice of Results to w.
func WriteReport(w io.Writer, results []Result, opts ReportOptions) {
	for _, r := range results {
		if !r.HasDiff() {
			continue
		}
		fmt.Fprintf(w, "\n[%s]\n", r.Path)
		fmt.Fprintln(w, strings.Repeat("-", 40))

		for _, k := range SortedKeys(r.Added) {
			val := maskIf(r.Added[k], opts.MaskValues)
			line := fmt.Sprintf("  + %-20s = %s", k, val)
			fmt.Fprintln(w, colorize(line, colorGreen, opts.Color))
		}

		for _, k := range SortedKeys(r.Removed) {
			val := maskIf(r.Removed[k], opts.MaskValues)
			line := fmt.Sprintf("  - %-20s = %s", k, val)
			fmt.Fprintln(w, colorize(line, colorRed, opts.Color))
		}

		for k, ch := range r.Changed {
			lv := maskIf(ch.Left, opts.MaskValues)
			rv := maskIf(ch.Right, opts.MaskValues)
			line := fmt.Sprintf("  ~ %-20s : %s -> %s", k, lv, rv)
			fmt.Fprintln(w, colorize(line, colorYellow, opts.Color))
		}
	}
}

func colorize(s, color string, enabled bool) string {
	if !enabled {
		return s
	}
	return color + s + colorReset
}

func maskIf(val string, mask bool) string {
	if mask {
		return "***"
	}
	return val
}
