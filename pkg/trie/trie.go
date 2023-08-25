package trie

import (
	"context"
	"fmt"
	"slices"
)

// Node defines a node participating in a trie tree
type Node[T any] struct {
	keyRune  *rune
	parent   *Node[T]
	children map[rune]*Node[T]
	values   []T
}

// Tree defines a trie tree that only allows bottom-up traversal.
// This is structured as such because it is assumed that, if a search term does not sufficiently exist in a particular node,
// then it will not sufficiently exist in parent nodes (e.g., if searching for 'o' and the current node is "BTC", then
// the parent nodes "BT" and "B" will never satisfy the search term, either).
// This allows the traversal of the tree to terminate and discard consideration of entire ancestries of nodes in the trie
// tree.
type Tree[T any] struct {
	leafNodes []*Node[T]
}

// GetLeafNodes gets all the leaf nodes of the tree.
func (t *Tree[T]) GetLeafNodes() []*Node[T] {
	return t.leafNodes
}

// KeyTermExtractor is a function used to, while loading a Node tree, extract the tree placement term from a given value.
type KeyTermExtractor[T any] func(context.Context, T) (string, error)

// LoadTree builds a Trie tree from the given items, using the given termExtractor to extract the tree placement term from each item.
func LoadTree[T any](ctx context.Context, items []T, termExtractor KeyTermExtractor[T]) (*Tree[T], error) {
	rootNode := newTrieNode[T](nil, nil)
	for _, item := range items {
		itemKeyTerm, err := termExtractor(ctx, item)
		if err != nil {
			return nil, fmt.Errorf("failed to extract term from item: %w", err)
		}
		rootNode.addItem([]rune(normalizeTerm(itemKeyTerm)), item)
	}

	leafNodes := rootNode.getLeafNodes()

	return &Tree[T]{
		leafNodes: leafNodes,
	}, nil
}

func newTrieNode[T any](parentNode *Node[T], keyRune *rune) *Node[T] {
	return &Node[T]{
		parent:   parentNode,
		children: make(map[rune]*Node[T]),
		keyRune:  keyRune,
	}
}

// Contains determines if the key term of this node contains the characters of
// the given phrase in the order given. The key term does not need to contiguously contain
// the phrase's characters - e.g., if the phrase is 'cal', this will return true if the key term is 'catapult'.
func (n *Node[T]) Contains(phrase string) bool {
	normalizedPhrase := normalizeTerm(phrase)
	normalizedRunes := []rune(normalizedPhrase)

	currentNode := n
	runeIndex := len(normalizedRunes) - 1
	for {
		if currentNode == nil || currentNode.keyRune == nil || runeIndex < 0 {
			break
		}

		currentRune := normalizedRunes[runeIndex]

		// If the current node has the current rune being examined, proceed onto the next rune
		if *currentNode.keyRune == currentRune {
			runeIndex--
		}

		// Regardless, continue traversing up the tree
		currentNode = currentNode.parent
	}

	// were all the runes found in the order given, even if not contiguous?
	return runeIndex < 0
}

// GetParentNode gets this node's parent node; if this is the root node, this will return nil.
func (n *Node[T]) GetParentNode() *Node[T] {
	return n.parent
}

// GetKeyTerm gets the string value for which this node represents a word in the tree
func (n *Node[T]) GetKeyTerm() string {
	if n.keyRune == nil {
		return ""
	}

	var keyRunes = []rune{*n.keyRune}
	parentNode := n.parent
	for {
		if parentNode == nil || parentNode.keyRune == nil {
			break
		}

		keyRunes = append(keyRunes, *parentNode.keyRune)
		parentNode = parentNode.parent
	}
	slices.Reverse(keyRunes)
	return string(keyRunes)
}

// GetValues gets the values stored within this node.
func (n *Node[T]) GetValues() []T {
	return n.values
}

func (n *Node[T]) addItem(runes []rune, item T) {
	if len(runes) == 0 {
		n.values = append(n.values, item)
		return
	}

	firstRune := runes[0]
	if _, hasChild := n.children[firstRune]; !hasChild {
		n.children[firstRune] = newTrieNode[T](n, &runes[0])
	}

	n.children[firstRune].addItem(runes[1:], item)
}

// getLeafNodes gets all leaf nodes that exist beneath this node
func (n *Node[T]) getLeafNodes() []*Node[T] {
	if len(n.children) == 0 {
		return []*Node[T]{n}
	}

	var leafNodes []*Node[T]

	var candidateNodes []*Node[T]
	for _, childNode := range n.children {
		candidateNodes = append(candidateNodes, childNode)
	}

	// Don't use recursion just in case it's a very deep tree
	for {
		if len(candidateNodes) == 0 {
			return leafNodes
		}

		var nextCandidates []*Node[T]
		for _, childNode := range candidateNodes {
			if len(childNode.children) == 0 {
				leafNodes = append(leafNodes, childNode)
			} else {
				for _, grandChildNode := range childNode.children {
					nextCandidates = append(nextCandidates, grandChildNode)
				}
			}
		}

		candidateNodes = nextCandidates
	}
}
