package single_test

import (
	"context"
	"time"
)

type ChanTimer struct {
	nodeSearchIterationTimeChan chan time.Duration
	treeSearchTimeChan          chan time.Duration
	sortTimeChan                chan time.Duration
	collectChan                 chan bool

	totalNodeSearchIterationTime time.Duration
	nodeSearchIterationTimeCount int64

	totalTreeSearchTime time.Duration
	treeSearchCount     int64

	totalSortTime time.Duration
	sortTimeCount int64
}

func NewChanTimer(ctx context.Context) *ChanTimer {
	nodeSearchIterationTimeChan := make(chan time.Duration, 1000)
	treeSearchTimeChan := make(chan time.Duration, 1000)
	sortTimeChan := make(chan time.Duration)
	collectChan := make(chan bool)

	timer := &ChanTimer{
		nodeSearchIterationTimeChan: nodeSearchIterationTimeChan,
		treeSearchTimeChan:          treeSearchTimeChan,
		sortTimeChan:                sortTimeChan,
		collectChan:                 collectChan,
	}

	go func() {
		for {
			select {
			case nodeSearchTime := <-nodeSearchIterationTimeChan:
				timer.nodeSearchIterationTimeCount++
				timer.totalNodeSearchIterationTime += nodeSearchTime
			case sortTime := <-sortTimeChan:
				timer.sortTimeCount++
				timer.totalSortTime += sortTime
			case searchTime := <-treeSearchTimeChan:
				timer.treeSearchCount++
				timer.totalTreeSearchTime += searchTime
			case _, _ = <-collectChan:
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	return timer
}

// Collect closes the collection of metrics and calculates the total metrics
func (c *ChanTimer) Collect() {
	// Drain the sort times
	close(c.sortTimeChan)
	for sortTimeDuration := range c.sortTimeChan {
		c.sortTimeCount++
		c.totalSortTime += sortTimeDuration
	}

	// Drain the remaining search times
	close(c.nodeSearchIterationTimeChan)
	for searchTimeDuration := range c.nodeSearchIterationTimeChan {
		c.nodeSearchIterationTimeCount++
		c.totalNodeSearchIterationTime += searchTimeDuration
	}

	// Drain the remaining tree searches
	close(c.treeSearchTimeChan)
	for searchTimeDuration := range c.treeSearchTimeChan {
		c.treeSearchCount++
		c.totalTreeSearchTime += searchTimeDuration
	}

	// Tell the goroutine to stop running
	c.collectChan <- true
	close(c.collectChan)
}

// GetAverageNodeSearchIterationDuration gets the average amount of time spent evaluating a node as a match for a given search term.
func (c *ChanTimer) GetAverageNodeSearchIterationDuration() time.Duration {
	if c.nodeSearchIterationTimeCount == 0 {
		return time.Duration(0)
	}
	return c.totalNodeSearchIterationTime / time.Duration(c.nodeSearchIterationTimeCount)
}

// GetNodeSearchIterationCount gets the number of times individual nodes were evaluated for a match against a search term.
func (c *ChanTimer) GetNodeSearchIterationCount() int64 {
	return c.nodeSearchIterationTimeCount
}

// GetTotalNodeSearchIterationDuration gets the total amount of time spent evaluating individual nodes' candidacy for matching a given search term.
func (c *ChanTimer) GetTotalNodeSearchIterationDuration() time.Duration {
	return c.totalNodeSearchIterationTime
}

// GetTotalSortDuration gets the total amount of time spent sorting the results after evalutaing nodes
func (c *ChanTimer) GetTotalSortDuration() time.Duration {
	return c.totalSortTime
}

// GetTotalTreeSearchDuration gets the total time spent searching trees for data
func (c *ChanTimer) GetTotalTreeSearchDuration() time.Duration {
	return c.totalTreeSearchTime
}

func (c *ChanTimer) RecordNodeSearchIteration(_ context.Context, searchDuration time.Duration) error {
	c.nodeSearchIterationTimeChan <- searchDuration
	return nil
}

func (c *ChanTimer) RecordSortTime(_ context.Context, sortDuration time.Duration) error {
	c.sortTimeChan <- sortDuration
	return nil
}

func (c *ChanTimer) RecordTreeSearch(_ context.Context, searchDuration time.Duration) error {
	c.treeSearchTimeChan <- searchDuration
	return nil
}
