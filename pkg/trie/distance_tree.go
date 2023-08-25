package trie

import (
	"context"
	"fmt"
	"github.com/coinbase/fuzzy-trie/pkg/levenshtein"
	"math"
	"sort"
	"time"
)

var defaultTimer = &NoOpTimer{}

// ComparableFuzzable is a combinatory interface of Fuzzable and the comparable keyword
type ComparableFuzzable interface {
	comparable
	Fuzzable
}

// DistanceTrees is a search mechanism supporting fuzzy search across multiple trees, each
// one weighed lower in priority to the previous one.
// This allows, for example, the fuzzy searching of assets by their symbol and then by their name,
// weighing name as lower in priority of a match than the symbol.
type DistanceTrees[T ComparableFuzzable] struct {
	trees []*Tree[T]
	timer Timer
}

// DistanceResult is a result of a fuzzy search, containing the result and the distance from the search term.
type DistanceResult[T any] struct {
	Distances []*int
	Result    T
}

func NewDistanceTrees[T ComparableFuzzable](trees []*Tree[T]) *DistanceTrees[T] {
	return &DistanceTrees[T]{
		trees: trees,
		timer: defaultTimer,
	}
}

// Search searches the trees within this DistanceTrees instance for the given search term, returning
// results with their primary and secondary distances
func (wt *DistanceTrees[T]) Search(ctx context.Context, searchTerm string) ([]*DistanceResult[T], error) {
	results := make(map[T]*DistanceResult[T])

	for treeIndex := range wt.trees {
		if searchErr := wt.searchTree(ctx, treeIndex, searchTerm, results); searchErr != nil {
			return nil, fmt.Errorf("faield to search tree at index %d: %w", treeIndex, searchErr)
		}
	}

	var weightedResults []*DistanceResult[T]
	for item, weightedResult := range results {
		// Append the secondary distances, too
		weightedResult.Distances = append(weightedResult.Distances, item.GetSecondaryDistances()...)
		weightedResults = append(weightedResults, weightedResult)
	}

	wt.sortResults(ctx, weightedResults)

	return weightedResults, nil
}

// SetTimer sets the Timer implementation to be used by this tree to measure its behavior
func (wt *DistanceTrees[T]) SetTimer(timer Timer) {
	wt.timer = timer
}

// evaluate evaluates the given node against the given search term and, if applicable, calculates the weight for the given tree index.
// It returns a slice of the subsequent nodes, if any, to be examined.
func (wt *DistanceTrees[T]) evaluate(
	ctx context.Context,
	searchTerm string,
	treeIndex int,
	treeCount int,
	results map[T]*DistanceResult[T],
	node *Node[T],
) []*Node[T] {
	nodeSearchStart := time.Now()
	defer func() {
		_ = wt.timer.RecordNodeSearchIteration(ctx, time.Since(nodeSearchStart))
	}()

	// If this node can never contain the search term, skip it and its ancestors
	if !node.Contains(searchTerm) {
		return nil
	}

	levenshteinDistance := levenshtein.LevenshteinDistance(searchTerm, node.GetKeyTerm())

	for _, matchingNodeValue := range node.values {
		var distances []*int
		if existingResult, hasResult := results[matchingNodeValue]; hasResult {
			distances = existingResult.Distances
		} else {
			distances = make([]*int, treeCount+1)
			distanceResult := &DistanceResult[T]{
				Result:    matchingNodeValue,
				Distances: distances,
			}
			results[matchingNodeValue] = distanceResult
		}

		if distances[treeIndex] != nil {
			// the distance has already been calculated since this was a parent node to another node
			// that's been visited; don't re-calculate it
			continue
		}

		primaryDistanceFactor := matchingNodeValue.GetPrimaryDistanceFactor()
		if primaryDistanceFactor != nil {
			levenshteinDistance = int(float64(levenshteinDistance) * *primaryDistanceFactor)
		}
		weightedDistance := int(math.Pow(10, float64(treeIndex))) + levenshteinDistance
		distances[treeIndex] = &weightedDistance
	}

	// Continue crawling up the tree
	if parentNode := node.parent; parentNode != nil {
		return []*Node[T]{parentNode}
	}

	return nil
}

// searchTree searches the tree at the given index with the given search term, populating the results into the given results map.
func (wt *DistanceTrees[T]) searchTree(ctx context.Context, treeIndex int, searchTerm string, results map[T]*DistanceResult[T]) error {
	searchStart := time.Now()
	defer func() {
		_ = wt.timer.RecordTreeSearch(ctx, time.Since(searchStart))
	}()

	treeCount := len(wt.trees)
	tree := wt.trees[treeIndex]

	currentNodes := tree.GetLeafNodes()

	for {
		if ctxErr := context.Cause(ctx); ctxErr != nil {
			return ctxErr
		}

		var nextNodes []*Node[T]
		for _, currentNode := range currentNodes {
			nodeValues := currentNode.values
			if len(nodeValues) > 0 {
				// If at least one of the node's values' distance has been calculated for this index,
				// it can be assumed that all the node's distances have been calculated and
				// this does not need to run again.
				// Further, it can be assumed that this node's ancestors have been calculated elsewhere,
				// so break out completely from this traversal.
				if result0, hasResult0 := results[nodeValues[0]]; hasResult0 && result0.Distances[treeIndex] != nil {
					continue
				}

				nextEvaluationCandidates := wt.evaluate(ctx, searchTerm, treeIndex, treeCount, results, currentNode)
				if len(nextEvaluationCandidates) > 0 {
					nextNodes = append(nextNodes, nextEvaluationCandidates...)
				}
			} else if currentParent := currentNode.parent; currentParent != nil {
				// if the current node has no values, try proceeding onto its parent
				nextNodes = append(nextNodes, currentParent)
			}
		}

		if len(nextNodes) == 0 {
			break
		}

		currentNodes = nextNodes
	}

	return nil
}

// sortResults sorts the given DistanceResult objects according to their comparative distances.
func (wt *DistanceTrees[T]) sortResults(ctx context.Context, weightedResults []*DistanceResult[T]) {
	sortStart := time.Now()
	defer func() {
		_ = wt.timer.RecordSortTime(ctx, time.Since(sortStart))
	}()

	sort.Slice(weightedResults, func(i, j int) bool {
		for distanceIndex := 0; distanceIndex < len(weightedResults[i].Distances); distanceIndex++ {
			distance := weightedResults[i].Distances[distanceIndex]
			otherDistance := weightedResults[j].Distances[distanceIndex]
			// If i's distance here is nil and j's is non-nil, then this should actually be ranked as _greater_
			// because that means that i was not a Match for the search term. Conversely, if j is nil
			// and i is non-nil, then it should be considered 'less than' for purposes of sorting, as that means
			// i matched the search term.
			if distance == nil {
				if otherDistance != nil {
					return false
				}
			} else if otherDistance == nil {
				return true
			} else if *distance < *otherDistance {
				return true
			} else if *distance > *otherDistance {
				return false
			}
		}

		return false
	})
}
