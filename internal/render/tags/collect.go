package hashtag

import (
	"bytes"
	"sort"
	"sync"
)

type Collector struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

func NewCollector() *Collector {
	return &Collector{seen: make(map[string]struct{})}
}

func (c *Collector) Add(tag []byte) {
	if c == nil || len(tag) == 0 {
		return
	}

	key := string(bytes.Clone(tag))

	c.mu.Lock()
	c.seen[key] = struct{}{}
	c.mu.Unlock()
}

func (c *Collector) Tags() []string {
	if c == nil {
		return nil
	}

	c.mu.Lock()
	out := make([]string, 0, len(c.seen))
	for k := range c.seen {
		out = append(out, k)
	}
	c.mu.Unlock()

	sort.Strings(out)
	return out
}
