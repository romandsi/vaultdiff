package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/vaultdiff/internal/diff"
	"github.com/yourorg/vaultdiff/internal/output"
)

var sampleResults = []diff.Result{
	{Path: "secret/app", Key: "DB_PASS", Kind: diff.Changed},
	{Path: "secret/app", Key: "API_KEY", Kind: diff.Added},
}

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		want     output.Format
	}{
		{"text", output.FormatText},
		{"", output.FormatText},
		{"json", output.FormatJSON},
		{"markdown", output.FormatMarkdown},
		{"md", output.FormatMarkdown},
	}
	for _, tc := range cases {
		got, err := output.ParseFormat(tc.input)
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseFormat(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := output.ParseFormat("xml")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestWrite_Text_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	if err := output.Write(&buf, nil, "dev", "prod", output.FormatText); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestWrite_Text_WithDiffs(t *testing.T) {
	var buf bytes.Buffer
	if err := output.Write(&buf, sampleResults, "dev", "prod", output.FormatText); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if !strings.Contains(got, "DB_PASS") || !strings.Contains(got, "API_KEY") {
		t.Errorf("expected keys in output, got: %s", got)
	}
}

func TestWrite_JSON(t *testing.T) {
	var buf bytes.Buffer
	if err := output.Write(&buf, sampleResults, "dev", "prod", output.FormatJSON); err != nil {
		t.Fatal(err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload["env_a"] != "dev" || payload["env_b"] != "prod" {
		t.Errorf("unexpected env fields: %v", payload)
	}
}

func TestWrite_Markdown(t *testing.T) {
	var buf bytes.Buffer
	if err := output.Write(&buf, sampleResults, "dev", "prod", output.FormatMarkdown); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if !strings.Contains(got, "## Vault Diff") {
		t.Errorf("expected markdown header, got: %s", got)
	}
	if !strings.Contains(got, "| Path |") {
		t.Errorf("expected markdown table header, got: %s", got)
	}
}
