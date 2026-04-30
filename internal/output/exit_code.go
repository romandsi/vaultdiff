package output

import "github.com/yourorg/vaultdiff/internal/diff"

// Exit codes used by the CLI.
const (
	// ExitOK indicates no differences were found.
	ExitOK = 0
	// ExitDiff indicates one or more differences were found.
	ExitDiff = 1
	// ExitError indicates a runtime error occurred.
	ExitError = 2
)

// ResolveExitCode returns the appropriate exit code based on diff results.
// It returns ExitDiff if any results are present, otherwise ExitOK.
func ResolveExitCode(results []diff.Result) int {
	if len(results) > 0 {
		return ExitDiff
	}
	return ExitOK
}
