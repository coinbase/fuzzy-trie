package single_test

import (
	"bufio"
	"compress/bzip2"
	"context"
	"fmt"
	"github.com/coinbase/fuzzy-trie/pkg/trie"
	"io"
	"net/http"
	"os"
	"runtime/trace"
	"strings"
	"testing"
	"time"
)

var testData []*testDatum

func getTestData() []*testDatum {
	if testData != nil {
		return testData
	}

	fmt.Println("Loading benchmarking data...")

	ctx, cancelFn := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFn()

	// Does the data already exist locally? If so, use that and save a download
	var bz2DataReader io.ReadCloser
	if fileRef, err := os.Open("ghProjectInfo2013-Feb.txt.bz2"); err == nil {
		bz2DataReader = fileRef
	} else if !os.IsNotExist(err) {
		panic(fmt.Sprintf("failed to open local data: %v", err))
	} else {
		fmt.Println("Downloading benchmarking data")
		// Download the data and copy it locally for re-use
		dataRequest, dataRequestErr := http.NewRequestWithContext(ctx, http.MethodGet, "http://flossdata.syr.edu/data/gh/2013/2013-Feb/ghProjectInfo2013-Feb.txt.bz2", nil)
		if dataRequestErr != nil {
			panic(fmt.Sprintf("failed to build request for data: %v", dataRequestErr))
		}

		dataResponse, dataResponseErr := http.DefaultClient.Do(dataRequest)
		if dataResponseErr != nil {
			panic(fmt.Sprintf("failed to execute GET for data: %v", dataResponseErr))
		} else if dataResponse.StatusCode != http.StatusOK {
			panic(fmt.Sprintf("unexpected status code reading data: %d", dataResponse.StatusCode))
		}

		bz2DataWriter, openErr := os.Create("ghProjectInfo2013-Feb.txt.bz2")
		if openErr != nil {
			panic(fmt.Sprintf("failed to open local data: %v", openErr))
		}

		if _, copyErr := io.Copy(bz2DataWriter, dataResponse.Body); copyErr != nil {
			panic(fmt.Sprintf("failed to copy data: %v", copyErr))
		} else if closeErr := bz2DataWriter.Close(); closeErr != nil {
			panic(fmt.Sprintf("failed to close local data writer: %v", closeErr))
		}

		var localOpenErr error
		bz2DataReader, localOpenErr = os.Open("ghProjectInfo2013-Feb.txt.bz2")
		if localOpenErr != nil {
			panic(fmt.Sprintf("failed to open local data after downloading: %v", localOpenErr))
		}
	}

	defer func() {
		_ = bz2DataReader.Close()
	}()

	dataReader := bufio.NewReader(bzip2.NewReader(bz2DataReader))
	lineNumber := 0
	for {
		if ctxErr := context.Cause(ctx); ctxErr != nil {
			panic(ctxErr)
		}

		line, readErr := dataReader.ReadString('\n')
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			panic(fmt.Sprintf("failed to read line: %v", readErr))
		}

		lineNumber++

		if strings.HasPrefix(line, "#") || len(strings.TrimSpace(line)) == 0 {
			// it's a comment; skip it or just some formatting to group data; ignore it
			continue
		}

		splitLine := strings.Split(line, "\t")
		if len(splitLine) != 3 {
			panic(fmt.Sprintf("unexpected line format at line %d: %s", lineNumber, line))
		}

		projectName := splitLine[0]
		developerName := splitLine[1]

		testData = append(testData, &testDatum{
			projectName:   projectName,
			developerName: developerName,
		})
	}

	fmt.Printf("Loaded %d benchmark data points.\n", len(testData))

	return testData
}

// beginTrace begins tracing if it is enabled.
// It always returns a no-arg function that should be invoked to end the trace.
func beginTrace(traceOutputWriter io.Writer) func() {
	if !isTraceEnabled() {
		return func() {}
	}

	if startErr := trace.Start(traceOutputWriter); startErr != nil {
		panic(fmt.Sprintf("failed to start trace: %v", startErr))
	}
	return trace.Stop
}

// getTraceWriter builds an io.Writer that can be used to write the trace data for the given benchmark testing
func getTraceWriter(b *testing.B) io.WriteCloser {
	if !isTraceEnabled() {
		return &noOpWriteCloser{}
	}

	traceFilename := fmt.Sprintf("trace_%s.out", b.Name())
	if removeErr := os.Remove(traceFilename); removeErr != nil && !os.IsNotExist(removeErr) {
		panic(fmt.Sprintf("failed to remove trace file '%s': %v", traceFilename, removeErr))
	}

	fileRef, err := os.Create(traceFilename)
	if err != nil {
		panic(fmt.Sprintf("failed to create trace file '%s': %v", traceFilename, err))
	}

	return fileRef
}

// isTraceEnabled determines if tracing is enabled.
func isTraceEnabled() bool {
	return os.Getenv("TRACE_ENABLED") == "true"
}

func loadDeveloperNameTree(ctx context.Context, testData []*testDatum) (*trie.Tree[*testDatum], error) {
	return trie.LoadTree(ctx, testData, func(ctx context.Context, item *testDatum) (string, error) {
		return item.developerName, nil
	})
}

func loadProjectNameTree(ctx context.Context, testData []*testDatum) (*trie.Tree[*testDatum], error) {
	return trie.LoadTree(ctx, testData, func(ctx context.Context, item *testDatum) (string, error) {
		return item.projectName, nil
	})
}

type testDatum struct {
	projectName   string
	developerName string
}

func (t *testDatum) GetPrimaryDistanceFactor() *float64 {
	return nil
}

func (t *testDatum) GetSecondaryDistances() []*int {
	return nil
}

func (t *testDatum) SortingGroup() int {
	return 1
}

// noOpWriteCloser is a no-op implementation of io.WriteCloser
type noOpWriteCloser struct {
}

func (noOpWriteCloser) Write(_ []byte) (n int, err error) {
	return 0, nil
}

func (noOpWriteCloser) Close() error {
	return nil
}
