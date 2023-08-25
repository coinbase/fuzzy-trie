package trie_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTrie(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Trie Suite")
}
