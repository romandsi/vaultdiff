package diff

import "sort"

// SecretMap holds key-value pairs for a single secret path.
type SecretMap map[string]string

// Result represents the diff outcome between two environments for a single path.
type Result struct {
	Path    string
	Added   map[string]string // keys present in right but not left
	Removed map[string]string // keys present in left but not right
	Changed map[string]Change // keys present in both but with different values
}

// Change holds the before/after values for a modified key.
type Change struct {
	Left  string
	Right string
}

// HasDiff returns true if there are any differences.
func (r Result) HasDiff() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Changed) > 0
}

// Compare computes the diff between two SecretMaps for a given path.
func Compare(path string, left, right SecretMap) Result {
	result := Result{
		Path:    path,
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string]Change),
	}

	for k, lv := range left {
		if rv, ok := right[k]; !ok {
			result.Removed[k] = lv
		} else if lv != rv {
			result.Changed[k] = Change{Left: lv, Right: rv}
		}
	}

	for k, rv := range right {
		if _, ok := left[k]; !ok {
			result.Added[k] = rv
		}
	}

	return result
}

// SortedKeys returns the keys of a map in sorted order.
func SortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
