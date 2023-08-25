package single_test

import (
	"context"
	"github.com/coinbase/fuzzy-trie/pkg/trie"
	"testing"
	"time"
)

func BenchmarkSearchSingleTree100(b *testing.B) {
	benchmarkSearchSingleTree(100, b)
}

func BenchmarkSearchSingleTree1000(b *testing.B) {
	benchmarkSearchSingleTree(1000, b)
}

func BenchmarkSearchSingleTree10000(b *testing.B) {
	benchmarkSearchSingleTree(10_000, b)
}

func BenchmarkSearchSingleTree100000(b *testing.B) {
	benchmarkSearchSingleTree(100_000, b)
}

func BenchmarkSearchSingleTree1000000(b *testing.B) {
	benchmarkSearchSingleTree(1_000_000, b)
}

func BenchmarkSearchSingleTree3000000(b *testing.B) {
	benchmarkSearchSingleTree(3_000_000, b)
}

func benchmarkSearchSingleTree(dataCount int, b *testing.B) {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFn()

	tree, loadErr := loadDeveloperNameTree(ctx, getTestData()[:dataCount])
	if loadErr != nil {
		panic(loadErr)
	}

	chanTimer := NewChanTimer(ctx)

	distanceTree := trie.NewDistanceTrees([]*trie.Tree[*testDatum]{tree})
	distanceTree.SetTimer(chanTimer)

	traceWriter := getTraceWriter(b)
	defer func() {
		_ = traceWriter.Close()
	}()

	traceStop := beginTrace(traceWriter)
	defer traceStop()

	b.ResetTimer()

	results, searchErr := distanceTree.Search(ctx, "e")

	if searchErr != nil {
		panic(searchErr)
	}

	b.StopTimer()

	chanTimer.Collect()

	b.ReportMetric(float64(len(results)), "search_results")
	b.ReportMetric(float64(chanTimer.GetNodeSearchIterationCount()), "search_count_total")
	b.ReportMetric(float64(chanTimer.GetAverageNodeSearchIterationDuration().Milliseconds()), "search_node_duration_average_millis")
	b.ReportMetric(float64(chanTimer.GetTotalNodeSearchIterationDuration().Milliseconds()), "search_node_duration_total_millis")
	b.ReportMetric(float64(chanTimer.GetTotalSortDuration().Milliseconds()), "sort_duration_total_millis")
	b.ReportMetric(float64(chanTimer.GetTotalTreeSearchDuration().Milliseconds()), "search_tree_duration_total_millis")

	b.ReportMetric(float64(len(tree.GetLeafNodes())), "leaf_node_name_dev_count")
}
