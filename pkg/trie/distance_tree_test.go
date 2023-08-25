package trie_test

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/coinbase/fuzzy-trie/pkg/trie"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"hash/fnv"
	"strings"
	"time"
)

//go:embed animals.txt
var animalsText string

var _ = Describe("DistanceTree", func() {
	var ctx context.Context

	BeforeEach(func() {
		var cancelFn context.CancelFunc
		ctx, cancelFn = context.WithTimeout(context.Background(), 5*time.Second)
		DeferCleanup(cancelFn)
	})

	Context("Search", func() {
		var tree *trie.DistanceTrees[*testComparableFuzzable]

		BeforeEach(func() {
			animals := strings.Split(animalsText, "\n")
			animalsFuzzable := make([]*testComparableFuzzable, len(animals))
			for i, animal := range animals {
				animalsFuzzable[i] = newTestComparableFuzzable(animal)
			}

			animalsTree, err := trie.LoadTree[*testComparableFuzzable](ctx, animalsFuzzable, func(_ context.Context, item *testComparableFuzzable) (string, error) {
				return item.text, nil
			})
			Expect(err).ToNot(HaveOccurred(), "loading the animals tree should not fail")

			tree = trie.NewDistanceTrees[*testComparableFuzzable]([]*trie.Tree[*testComparableFuzzable]{animalsTree})
		})

		Context("substring matching", func() {
			It("returns the number of results in the expected order", func() {
				results, err := tree.Search(ctx, "cat")
				Expect(err).ToNot(HaveOccurred(), "searching the tree should not fail")
				Expect(results).To(HaveLen(21), "the correct number of results should be returned")
				// Sample to make sure the expected order is maintained
				Expect(results[0].Result.text).To(Equal("Cat"), "the 0th element should be correct")
				Expect(results[7].Result.text).To(Equal("Wildcat"), "the 7th element should be correct")
				Expect(results[16].Result.text).To(Equal("Domestic rabbit"), "the 16th element should be correct")
				Expect(results[20].Result.text).To(Equal("Domestic Bactrian camel"), "the 20st element should be correct")
			})
		})

		Context("fuzzy matching", func() {
			It("should return fuzzily-matched results", func() {
				results, err := tree.Search(ctx, "wol")
				Expect(err).ToNot(HaveOccurred(), "searching the tree should not fail")
				Expect(results).To(HaveLen(8), "the correct number of results should be returned")
				// Sample to make sure the expected order is maintained
				Expect(results[0].Result.text).To(Equal("Wolf"), "the 0th element should be correct")
				Expect(results[3].Result.text).To(Equal("Wolverine"), "the 3rd element should be correct")
				Expect(results[6].Result.text).To(Equal("New World quail"), "the 6th element should be correct")
			})
		})
	})
})

type testComparableFuzzable struct {
	text string
}

func newTestComparableFuzzable(text string) *testComparableFuzzable {
	return &testComparableFuzzable{
		text: text,
	}
}

func (t *testComparableFuzzable) GetPrimaryDistanceFactor() *float64 {
	return nil
}

func (t *testComparableFuzzable) GetSecondaryDistances() []*int {
	// In the event of a tie, use the length of the string as a secondary distance
	textLength := len(t.text)
	// If the length is still equal, then use the hash of the string
	hash := fnv.New32a()
	if _, hashErr := hash.Write([]byte(t.text)); hashErr != nil {
		panic(fmt.Sprintf("unable to generate hash for '%s': %v", t.text, hashErr))
	}
	hashSum := int(hash.Sum32())
	return []*int{&textLength, &hashSum}
}

func (t *testComparableFuzzable) SortingGroup() int {
	return 1
}
