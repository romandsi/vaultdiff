package output_test

import (
	"testing"

	"github.com/yourorg/vaultdiff/internal/diff"
	"github.com/yourorg/vaultdiff/internal/output"
)

func TestResolveExitCode_NoDiff(t *testing.T) {
	code := output.ResolveExitCode(nil)
	if code != output.ExitOK {
		t.Errorf("expected ExitOK (%d), got %d", output.ExitOK, code)
	}
}

func TestResolveExitCode_EmptySlice(t *testing.T) {
	code := output.ResolveExitCode([]diff.Result{})
	if code != output.ExitOK {
		t.Errorf("expected ExitOK (%d), got %d", output.ExitOK, code)
	}
}

func TestResolveExitCode_WithDiffs(t *testing.T) {
	results := []diff.Result{
		{Path: "secret/app", Key: "TOKEN", Kind: diff.Changed},
	}
	code := output.ResolveExitCode(results)
	if code != output.ExitDiff {
		t.Errorf("expected ExitDiff (%d), got %d", output.ExitDiff, code)
	}
}

func TestResolveExitCode_MultipleResults(t *testing.T) {
	results := []diff.Result{
		{Path: "secret/db", Key: "PASSWORD", Kind: diff.Removed},
		{Path: "secret/db", Key: "HOST", Kind: diff.Added},
	}
	code := output.ResolveExitCode(results)
	if code != output.ExitDiff {
		t.Errorf("expected ExitDiff (%d), got %d", output.ExitDiff, code)
	}
}
