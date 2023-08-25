package single_test

import (
	"context"
	"testing"
	"time"
)

func BenchmarkLoadTree100(b *testing.B) {
	benchmarkLoadTree(100, b)
}

func BenchmarkLoadTree1000(b *testing.B) {
	benchmarkLoadTree(1000, b)
}

func BenchmarkLoadTree10000(b *testing.B) {
	benchmarkLoadTree(10_000, b)
}

func BenchmarkLoadTree100000(b *testing.B) {
	benchmarkLoadTree(100_000, b)
}

func BenchmarkLoadTree1000000(b *testing.B) {
	benchmarkLoadTree(1_000_000, b)
}

func BenchmarkLoadTree3000000(b *testing.B) {
	benchmarkLoadTree(3_000_000, b)
}

func benchmarkLoadTree(dataCount int, b *testing.B) {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFn()

	dataSubset := getTestData()[:dataCount]

	traceWriter := getTraceWriter(b)
	defer func() {
		_ = traceWriter.Close()
	}()

	traceStop := beginTrace(traceWriter)
	defer traceStop()

	b.ResetTimer()

	if _, loadErr := loadDeveloperNameTree(ctx, dataSubset); loadErr != nil {
		panic(loadErr)
	}
}
