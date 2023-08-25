package trie

import (
	"context"
	"time"
)

// Timer defines a means of recording time measurements of operations of a DistanceTrees.
type Timer interface {
	// RecordTreeSearch records the amount of time spent searching the nodes of a single Tree instance.
	RecordTreeSearch(ctx context.Context, searchDuration time.Duration) error
	// RecordNodeSearchIteration records the amount of time spent evaluating an individual Node's candidacy for matching a search phrase.
	RecordNodeSearchIteration(ctx context.Context, searchDuration time.Duration) error
	// RecordSortTime records the amount of time spent sorting the final results in a search operation.
	RecordSortTime(ctx context.Context, sortDuration time.Duration) error
}

// NoOpTimer is a Timer implementation that doesn't do anything
type NoOpTimer struct {
}

func (NoOpTimer) RecordNodeSearchIteration(_ context.Context, _ time.Duration) error {
	return nil
}

func (NoOpTimer) RecordSortTime(_ context.Context, _ time.Duration) error {
	return nil
}

func (NoOpTimer) RecordTreeSearch(_ context.Context, _ time.Duration) error {
	return nil
}
