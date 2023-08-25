package trie_test

import (
	"context"
	"github.com/coinbase/fuzzy-trie/pkg/trie"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"strings"
	"time"
)

var _ = Describe("Trie", func() {
	var ctx context.Context

	BeforeEach(func() {
		var cancelFn context.CancelFunc
		ctx, cancelFn = context.WithTimeout(context.Background(), 5*time.Second)
		DeferCleanup(cancelFn)
	})

	Context("LoadTree", func() {
		It("should return a trie", func() {
			cats := []*testItem{
				{
					text: "cat",
				},
				{
					text: "cat",
				},
			}
			cataracts := &testItem{
				text: "cataracts",
			}
			dog := &testItem{
				text: "dog",
			}

			tree, err := trie.LoadTree[*testItem](ctx, append(cats, cataracts, dog), func(ctx context.Context, item *testItem) (string, error) {
				return item.GetText(), nil
			})

			Expect(err).ToNot(HaveOccurred(), "loading the trie tree should not fail")

			getLeafNode := func(ctx context.Context, term string) *trie.Node[*testItem] {
				currentNodes := tree.GetLeafNodes()
				for {
					Expect(context.Cause(ctx)).ToNot(HaveOccurred(), "the context should not be cancelled")
					var parentNodes []*trie.Node[*testItem]
					for _, currentNode := range currentNodes {
						if currentNode.GetKeyTerm() == strings.ToUpper(term) {
							return currentNode
						}
						if parentNode := currentNode.GetParentNode(); parentNode != nil {
							parentNodes = append(parentNodes, parentNode)
						}
					}
					currentNodes = parentNodes
				}
			}

			catNode := getLeafNode(ctx, "cat")
			Expect(catNode).ToNot(BeNil(), "the cat node should exist")
			Expect(catNode.GetValues()).To(ConsistOf(cats), "the cat node should contain all of the cat values")
			Expect(catNode.GetKeyTerm()).To(Equal("CAT"), "the cat node should have the correct key term")

			cataractsNode := getLeafNode(ctx, "cataracts")
			Expect(cataractsNode).ToNot(BeNil(), "the cataracts node should exist")
			Expect(cataractsNode.GetValues()).To(ConsistOf(cataracts), "the cataracts node should contain all of the cataracts values")
			Expect(cataractsNode.GetKeyTerm()).To(Equal("CATARACTS"), "the cataracts node should have the correct key term")

			dogNode := getLeafNode(ctx, "dog")
			Expect(dogNode).ToNot(BeNil(), "the dog node should exist")
			Expect(dogNode.GetValues()).To(ConsistOf(dog), "the dog node should contain all of the dog values")
			Expect(dogNode.GetKeyTerm()).To(Equal("DOG"), "the dog node should have the correct key term")
		})
	})
})

type testItem struct {
	text string
}

func (t *testItem) GetText() string {
	return t.text
}
