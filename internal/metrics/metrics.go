package metrics

import (
	"fmt"
	"sync"
)

var mu sync.Mutex
var counters = map[string]int{}

// Incr simulates DogStatsD `increment` style metrics
func Incr(name string) {
	mu.Lock()
	defer mu.Unlock()
	counters[name]++
	fmt.Printf("📈 %s: %d\n", name, counters[name])
}

// GetAll returns current metrics (optional: for HTTP /metrics)
func GetAll() map[string]int {
	mu.Lock()
	defer mu.Unlock()
	snapshot := make(map[string]int)
	for k, v := range counters {
		snapshot[k] = v
	}
	return snapshot
}
