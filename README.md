# Fuzzy Trie

A [trie](https://en.wikipedia.org/wiki/Trie)-based implementation of fuzzy searching.

This uses a [Levenshtein distance](https://en.wikipedia.org/wiki/Levenshtein_distance) calculation to evaluate the closeness of an item in the tree to the search term.

## Usage

The items to be loaded into the tree must implement the `Fuzzable` interface defined in this library.

Once you have implemented that interface for your items to be loaded into the tree, you can do the following:

```
ctx := context.Background()
fuzzables := getMyItems() // returns structs that implement Fuzzable - not Fuzzable, as this must also satisfy ComparableFuzzable
tree := trie.LoadTree(ctx, fuzzables, func(ctx context.Context, item *myFuzzableImpl) (string, error) {
    return item.trieNodeKey, nil
})
searchableTree := trie.NewDistanceTrees([]*trie.Tree[*myFuzzableImpl]{tree})
searchResults, err := searchableTree.Search(ctx, "cat")
```

### Multi-Dimensional Trees

If you wish to evaluate _closeness_ to your search term with multiple fields on your items (e.g., perhaps supporting search by a family name and then using relevance of given name as a tie-breaker), you can provide multiple trees to the `DistanceTree` to execute such functionality:

```
ctx := context.Background()
people := getPeople() // struct defined as { familyName string, givenName string }
familyNameTree := trie.LoadTree(ctx, people, func(ctx context.Context, item *person) (string, error) {
    return person.familyName, nil
})
givenNameTree := trie.LoadTree(ctx, people, func(ctx context.Context, item *person) (string, error) {
    return person.givenName, nil
})
searchableTree := trie.NewDistanceTrees([]*trie.Tree[*person]{familyNameTree, givenNameTree})
searchResults, err := searchableTree.Search(ctx, "jo")
```

The above tree will first find results that have a family name close to 'jo' and, for cases where multiple people are equally close to that search term, will then evaluate the closeness of the person's given name to 'jo' and return the results in that order.

### Measurement

If you wish to measure the performance of this tree within your application, you can supply an implementation of the `trie.Timer` interface provided in this library and use the `SetTimer` method on the `DistanceTrees` struct to inject your implementation.

## Benchmarking

Refer to [benchmarking.md](./internal/benchmark/benchmarking.md) for more information.