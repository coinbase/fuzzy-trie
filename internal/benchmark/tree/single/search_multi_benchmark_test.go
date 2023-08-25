package single_test

import (
	"context"
	"fmt"
	"github.com/coinbase/fuzzy-trie/pkg/trie"
	"testing"
	"time"
)

const singleCharPhrase = "e"
const multiCharPhrase = "fed"
const lowResultCountPhrase = "fuz"

// Tests that exercise searching with a high-result-count, short search phrase
func BenchmarkMultiTree100(b *testing.B) {
	benchmarkMultiTree(100, singleCharPhrase, b)
}

func BenchmarkMultiTree1000(b *testing.B) {
	benchmarkMultiTree(1000, singleCharPhrase, b)
}

func BenchmarkMultiTree10000(b *testing.B) {
	benchmarkMultiTree(10_000, singleCharPhrase, b)
}

func BenchmarkMultiTree100000(b *testing.B) {
	benchmarkMultiTree(100_000, singleCharPhrase, b)
}

func BenchmarkMultiTree1000000(b *testing.B) {
	benchmarkMultiTree(1_000_000, singleCharPhrase, b)
}

func BenchmarkMultiTree3000000(b *testing.B) {
	benchmarkMultiTree(3_000_000, singleCharPhrase, b)
}

// Tests that execute with a multi-character search that has a high number of matches
func BenchmarkMultiTreeMultiCharPhrase100(b *testing.B) {
	benchmarkMultiTree(100, multiCharPhrase, b)
}

func BenchmarkMultiTreeMultiCharPhrase1000(b *testing.B) {
	benchmarkMultiTree(1000, multiCharPhrase, b)
}

func BenchmarkMultiTreeMultiCharPhrase10000(b *testing.B) {
	benchmarkMultiTree(10_000, multiCharPhrase, b)
}

func BenchmarkMultiTreeMultiCharPhrase100000(b *testing.B) {
	benchmarkMultiTree(100_000, multiCharPhrase, b)
}

func BenchmarkMultiTreeMultiCharPhrase1000000(b *testing.B) {
	benchmarkMultiTree(1_000_000, multiCharPhrase, b)
}

func BenchmarkMultiTreeMultiCharPhrase3000000(b *testing.B) {
	benchmarkMultiTree(3_000_000, multiCharPhrase, b)
}

// Tests that return few results
func BenchmarkMultiTreeLowResultCountPhrase100(b *testing.B) {
	benchmarkMultiTree(100, lowResultCountPhrase, b)
}

func BenchmarkMultiTreeLowResultCountPhrase1000(b *testing.B) {
	benchmarkMultiTree(1000, lowResultCountPhrase, b)
}

func BenchmarkMultiTreeLowResultCountPhrase10000(b *testing.B) {
	benchmarkMultiTree(10_000, lowResultCountPhrase, b)
}

func BenchmarkMultiTreeLowResultCountPhrase100000(b *testing.B) {
	benchmarkMultiTree(100_000, lowResultCountPhrase, b)
}

func BenchmarkMultiTreeLowResultCountPhrase1000000(b *testing.B) {
	benchmarkMultiTree(1_000_000, lowResultCountPhrase, b)
}

func BenchmarkMultiTreeLowResultCountPhrase3000000(b *testing.B) {
	benchmarkMultiTree(3_000_000, lowResultCountPhrase, b)
}

func benchmarkMultiTree(dataCount int, searchPhrase string, b *testing.B) {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFn()

	dataSubset := getTestData()[:dataCount]

	developerNameTree, err := loadDeveloperNameTree(ctx, dataSubset)
	if err != nil {
		panic(fmt.Sprintf("failed to load developer name tree: %v", err))
	}

	projectNameTree, err := loadProjectNameTree(ctx, dataSubset)
	if err != nil {
		panic(fmt.Sprintf("failed to load project name tree: %v", err))
	}

	chanTimer := NewChanTimer(ctx)

	distanceTree := trie.NewDistanceTrees([]*trie.Tree[*testDatum]{developerNameTree, projectNameTree})
	distanceTree.SetTimer(chanTimer)

	traceWriter := getTraceWriter(b)
	defer func() {
		_ = traceWriter.Close()
	}()

	traceStop := beginTrace(traceWriter)
	defer traceStop()

	b.ResetTimer()

	results, searchErr := distanceTree.Search(ctx, searchPhrase)
	if searchErr != nil {
		panic(fmt.Sprintf("failed to search: %v", searchErr))
	}

	b.StopTimer()

	chanTimer.Collect()

	b.ReportMetric(float64(len(results)), "search_results")
	b.ReportMetric(float64(chanTimer.GetNodeSearchIterationCount()), "search_count_total")
	b.ReportMetric(float64(chanTimer.GetAverageNodeSearchIterationDuration().Milliseconds()), "search_node_duration_average_millis")
	b.ReportMetric(float64(chanTimer.GetTotalNodeSearchIterationDuration().Milliseconds()), "search_node_duration_total_millis")
	b.ReportMetric(float64(chanTimer.GetTotalSortDuration().Milliseconds()), "sort_duration_total_millis")
	b.ReportMetric(float64(chanTimer.GetTotalTreeSearchDuration().Milliseconds()), "search_tree_duration_total_millis")

	b.ReportMetric(float64(len(developerNameTree.GetLeafNodes())), "leaf_node_name_dev_count")
	b.ReportMetric(float64(len(projectNameTree.GetLeafNodes())), "leaf_node_name_project_count")
}
