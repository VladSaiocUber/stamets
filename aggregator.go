package stamets

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

// AggregatePTAResults recursively walks a directory, reads every file,
// and then unparses and aggregates all potential PTA metrics in the contents.
func AggregatePTAResults(dir string) []PTAMetrics {
	results := make([]PTAMetrics, 0)

	wg, c := &sync.WaitGroup{}, make(chan struct{}, 10)

	dir, _ = filepath.Abs(dir)
	filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}

		bs, err := os.ReadFile(p)
		if err != nil && err != io.EOF {
			return nil
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-c }()
			c <- struct{}{}
			results = append(results, UnparsePTAResultsFromReader(bytes.NewReader(bs))...)
		}()
		return nil
	})

	wg.Wait()

	return results
}

// AggregateCallGraphResults recursively walks a directory, reads every file,
// and then unparses and aggregates all potential call graph metrics in the contents.
func AggregateCallGraphResults(dir string) []CallGraphMetrics {
	results := make([]CallGraphMetrics, 0)

	wg, c := &sync.WaitGroup{}, make(chan struct{}, 10)

	filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}

		bs, err := os.ReadFile(p)
		if err != nil {
			return nil
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-c }()
			c <- struct{}{}
			results = append(results, UnparseCallGraphMetricsFromReader(bytes.NewReader(bs))...)
		}()
		return nil
	})

	wg.Wait()

	return results
}
