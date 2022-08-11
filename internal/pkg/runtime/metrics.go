package runtime

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	"notifications/internal/pkg/metrics"
)

const (
	maxGCPausesTimings = 10
)

func CollectGoMetrics(ctx context.Context, metric metrics.Metrics, id string) {
	collect := goCollector(metric, id)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		collect()
		time.Sleep(time.Second)
	}
}

func goCollector(metrics metrics.Metrics, id string) func() {
	prefix := fmt.Sprintf(`go.%s.`, id)

	lastNumGC := int64(0)

	var statGC debug.GCStats
	var statMem runtime.MemStats

	return func() {

		// goroutines and threads
		metrics.Gauge(prefix+`goroutines`, runtime.NumGoroutine())
		n, _ := runtime.ThreadCreateProfile(nil)
		metrics.Gauge(prefix+`threads`, n)

		// garbage collector stats
		debug.ReadGCStats(&statGC)
		pausesCount := int(statGC.NumGC - lastNumGC)
		if pausesCount > len(statGC.Pause) {
			pausesCount = len(statGC.Pause) // 256*2+3
		}
		for i := 0; i < pausesCount && i < maxGCPausesTimings; i++ {
			metrics.Timing(
				prefix+`gc_pause_microseconds`,
				float64(statGC.Pause[i]*time.Microsecond)/float64(time.Millisecond),
			)
		}
		lastNumGC = statGC.NumGC

		// memory stats
		runtime.ReadMemStats(&statMem)
		metrics.Gauge(prefix+`mem_alloc_bytes`, statMem.Alloc)
		metrics.Gauge(prefix+`mem_alloc_bytes_total`, statMem.TotalAlloc)
		metrics.Gauge(prefix+`mem_sys_bytes`, statMem.Sys)
		metrics.Gauge(prefix+`mem_heap_alloc_bytes`, statMem.HeapAlloc)
	}
}
