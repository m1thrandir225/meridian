package cache

import (
	"sync/atomic"
)

type CacheMetrics struct {
	hits   int64
	misses int64
}

func (m *CacheMetrics) IncrementHits() {
	atomic.AddInt64(&m.hits, 1)
}

func (m *CacheMetrics) IncrementMisses() {
	atomic.AddInt64(&m.misses, 1)
}

func (m *CacheMetrics) GetHits() int64 {
	return atomic.LoadInt64(&m.hits)
}

func (m *CacheMetrics) GetMisses() int64 {
	return atomic.LoadInt64(&m.misses)
}

func (m *CacheMetrics) GetHitRate() float64 {
	hits := m.GetHits()
	misses := m.GetMisses()
	total := hits + misses
	if total == 0 {
		return 0
	}
	return float64(hits) / float64(total)
}
