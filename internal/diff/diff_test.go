package diff_test

import (
	"testing"

	"github.com/your-org/vaultdiff/internal/diff"
)

func TestCompare_NoDiff(t *testing.T) {
	left := diff.SecretMap{"key1": "value1", "key2": "value2"}
	right := diff.SecretMap{"key1": "value1", "key2": "value2"}

	result := diff.Compare("secret/myapp", left, right)

	if result.HasDiff() {
		t.Errorf("expected no diff, got: added=%v removed=%v changed=%v",
			result.Added, result.Removed, result.Changed)
	}
}

func TestCompare_Added(t *testing.T) {
	left := diff.SecretMap{"key1": "value1"}
	right := diff.SecretMap{"key1": "value1", "key2": "new"}

	result := diff.Compare("secret/myapp", left, right)

	if v, ok := result.Added["key2"]; !ok || v != "new" {
		t.Errorf("expected key2 to be added with value 'new', got: %v", result.Added)
	}
	if len(result.Removed) != 0 || len(result.Changed) != 0 {
		t.Errorf("unexpected removed or changed entries")
	}
}

func TestCompare_Removed(t *testing.T) {
	left := diff.SecretMap{"key1": "value1", "key2": "old"}
	right := diff.SecretMap{"key1": "value1"}

	result := diff.Compare("secret/myapp", left, right)

	if v, ok := result.Removed["key2"]; !ok || v != "old" {
		t.Errorf("expected key2 to be removed with value 'old', got: %v", result.Removed)
	}
}

func TestCompare_Changed(t *testing.T) {
	left := diff.SecretMap{"key1": "before"}
	right := diff.SecretMap{"key1": "after"}

	result := diff.Compare("secret/myapp", left, right)

	change, ok := result.Changed["key1"]
	if !ok {
		t.Fatal("expected key1 to be in Changed")
	}
	if change.Left != "before" || change.Right != "after" {
		t.Errorf("unexpected change values: %+v", change)
	}
}

func TestCompare_Path(t *testing.T) {
	result := diff.Compare("secret/test", diff.SecretMap{}, diff.SecretMap{})
	if result.Path != "secret/test" {
		t.Errorf("expected path 'secret/test', got '%s'", result.Path)
	}
}

func TestCompare_MultipleChanges(t *testing.T) {
	left := diff.SecretMap{"key1": "before", "key2": "old", "key3": "keep"}
	right := diff.SecretMap{"key1": "after", "key3": "keep", "key4": "new"}

	result := diff.Compare("secret/myapp", left, right)

	if len(result.Changed) != 1 {
		t.Errorf("expected 1 changed key, got %d", len(result.Changed))
	}
	if len(result.Removed) != 1 {
		t.Errorf("expected 1 removed key, got %d", len(result.Removed))
	}
	if len(result.Added) != 1 {
		t.Errorf("expected 1 added key, got %d", len(result.Added))
	}
	if !result.HasDiff() {
		t.Error("expected HasDiff to return true")
	}
}

func TestSortedKeys(t *testing.T) {
	m := map[string]string{"z": "1", "a": "2", "m": "3"}
	keys := diff.SortedKeys(m)
	expected := []string{"a", "m", "z"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("expected %s at index %d, got %s", expected[i], i, k)
		}
	}
}
