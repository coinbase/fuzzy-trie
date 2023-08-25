package trie

import "strings"

// normalizeTerm normalizes the given term to a format that ensures consistency of storage
func normalizeTerm(v string) string {
	return strings.ToUpper(v)
}
