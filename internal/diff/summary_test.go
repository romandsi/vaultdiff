package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestSummarize_Empty(t *testing.T) {
	s := Summarize(nil)
	if s.Total != 0 || s.Added != 0 || s.Removed != 0 || s.Changed != 0 || s.Unchanged != 0 {
		t.Errorf("expected all zeros for empty input, got %+v", s)
	}
}

func TestSummarize_Counts(t *testing.T) {
	entries := []DiffEntry{
		{Path: "secret/app", Key: "DB_PASS", Status: StatusChanged},
		{Path: "secret/app", Key: "API_KEY", Status: StatusAdded},
		{Path: "secret/app", Key: "OLD_KEY", Status: StatusRemoved},
		{Path: "secret/infra", Key: "TOKEN", Status: StatusUnchanged},
		{Path: "secret/infra", Key: "CERT", Status: StatusAdded},
	}

	s := Summarize(entries)

	if s.Total != 5 {
		t.Errorf("expected Total=5, got %d", s.Total)
	}
	if s.Added != 2 {
		t.Errorf("expected Added=2, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected Removed=1, got %d", s.Removed)
	}
	if s.Changed != 1 {
		t.Errorf("expected Changed=1, got %d", s.Changed)
	}
	if s.Unchanged != 1 {
		t.Errorf("expected Unchanged=1, got %d", s.Unchanged)
	}
}

func TestSummarize_ByPath(t *testing.T) {
	entries := []DiffEntry{
		{Path: "secret/app", Key: "A", Status: StatusAdded},
		{Path: "secret/app", Key: "B", Status: StatusChanged},
		{Path: "secret/infra", Key: "C", Status: StatusRemoved},
	}

	s := Summarize(entries)

	app, ok := s.ByPath["secret/app"]
	if !ok {
		t.Fatal("expected path secret/app in ByPath")
	}
	if app.Added != 1 || app.Changed != 1 || app.Removed != 0 {
		t.Errorf("unexpected secret/app stats: %+v", app)
	}

	infra, ok := s.ByPath["secret/infra"]
	if !ok {
		t.Fatal("expected path secret/infra in ByPath")
	}
	if infra.Removed != 1 || infra.Added != 0 {
		t.Errorf("unexpected secret/infra stats: %+v", infra)
	}
}

func TestWriteSummary_Output(t *testing.T) {
	entries := []DiffEntry{
		{Path: "secret/app", Key: "X", Status: StatusAdded},
		{Path: "secret/app", Key: "Y", Status: StatusRemoved},
	}
	s := Summarize(entries)

	var buf bytes.Buffer
	WriteSummary(&buf, s)
	out := buf.String()

	if !strings.Contains(out, "2 total keys") {
		t.Errorf("expected '2 total keys' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "secret/app") {
		t.Errorf("expected 'secret/app' in per-path breakdown, got:\n%s", out)
	}
}
