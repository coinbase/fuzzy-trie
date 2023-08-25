package trie

// Fuzzable describes an object that can participate in a fuzzy search
type Fuzzable interface {
	// GetPrimaryDistanceFactor gets the primary distance factor, if any.
	GetPrimaryDistanceFactor() *float64

	// GetSecondaryDistances gets any and all secondary instances in the event of otherwise-equal-weights to be used to elevate
	// the relevance of a term
	GetSecondaryDistances() []*int

	// SortingGroup gets the zero-based index of the sorting group for purposes of determining rank of an asset relative to others.
	// This allows members within a particular sorting group to be ranked relative to each other.
	SortingGroup() int
}
