package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteReport_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	diffs := []Diff{}
	WriteReport(&buf, diffs, false)
	out := buf.String()
	if !strings.Contains(out, "No differences") {
		t.Errorf("expected 'No differences' message, got: %s", out)
	}
}

func TestWriteReport_WithDiffs(t *testing.T) {
	var buf bytes.Buffer
	diffs := []Diff{
		{Path: "secret/app", Key: "DB_PASS", Type: DiffChanged, Left: "old", Right: "new"},
		{Path: "secret/app", Key: "API_KEY", Type: DiffAdded, Left: "", Right: "abc123"},
		{Path: "secret/app", Key: "OLD_KEY", Type: DiffRemoved, Left: "xyz", Right: ""},
	}
	WriteReport(&buf, diffs, false)
	out := buf.String()

	if !strings.Contains(out, "DB_PASS") {
		t.Errorf("expected DB_PASS in output")
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output")
	}
	if !strings.Contains(out, "OLD_KEY") {
		t.Errorf("expected OLD_KEY in output")
	}
}

func TestWriteReport_MaskValues(t *testing.T) {
	var buf bytes.Buffer
	diffs := []Diff{
		{Path: "secret/app", Key: "DB_PASS", Type: DiffChanged, Left: "supersecret", Right: "newpassword"},
	}
	WriteReport(&buf, diffs, true)
	out := buf.String()

	if strings.Contains(out, "supersecret") {
		t.Errorf("expected value to be masked, but found plain text")
	}
	if strings.Contains(out, "newpassword") {
		t.Errorf("expected value to be masked, but found plain text")
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected masked value '***' in output")
	}
}

func TestWriteReport_GroupsByPath(t *testing.T) {
	var buf bytes.Buffer
	diffs := []Diff{
		{Path: "secret/app", Key: "KEY1", Type: DiffAdded, Left: "", Right: "v1"},
		{Path: "secret/db", Key: "KEY2", Type: DiffRemoved, Left: "v2", Right: ""},
	}
	WriteReport(&buf, diffs, false)
	out := buf.String()

	if !strings.Contains(out, "secret/app") {
		t.Errorf("expected path 'secret/app' in output")
	}
	if !strings.Contains(out, "secret/db") {
		t.Errorf("expected path 'secret/db' in output")
	}
}
