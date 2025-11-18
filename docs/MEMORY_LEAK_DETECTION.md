# Memory Leak Detection Guide

Complete guide for detecting and analyzing memory leaks in REST API and gRPC services.

## Table of Contents

- [Memory Leak Detection Guide](#memory-leak-detection-guide)
  - [Table of Contents](#table-of-contents)
  - [Method 1: Using pprof (Recommended)](#method-1-using-pprof-recommended)
    - [1.1 Run the Application](#11-run-the-application)
    - [1.2 Access pprof Endpoints](#12-access-pprof-endpoints)
    - [1.3 Capture Heap Profile (Before Load Test)](#13-capture-heap-profile-before-load-test)
    - [1.4 Generate Load Test](#14-generate-load-test)
    - [1.5 Capture Heap Profile (After Load Test)](#15-capture-heap-profile-after-load-test)
    - [1.6 Analyze Memory Profile](#16-analyze-memory-profile)
      - [Option A: Interactive Analysis](#option-a-interactive-analysis)
      - [Option B: Compare Profiles](#option-b-compare-profiles)
      - [Option C: Generate Visual Report](#option-c-generate-visual-report)
    - [1.7 Continuous Profiling](#17-continuous-profiling)
  - [Method 2: Check Goroutine Leaks](#method-2-check-goroutine-leaks)
    - [2.1 Check Goroutine Count](#21-check-goroutine-count)
    - [2.2 Analyze Goroutine Profile](#22-analyze-goroutine-profile)
  - [Method 3: Runtime Memory Stats](#method-3-runtime-memory-stats)
    - [3.1 Create Monitoring Endpoint](#31-create-monitoring-endpoint)
    - [3.2 Monitor Memory Usage](#32-monitor-memory-usage)
  - [Method 4: Using External Tools](#method-4-using-external-tools)
    - [4.1 Continuous Profiler (Pyroscope)](#41-continuous-profiler-pyroscope)
    - [4.2 Prometheus + Grafana](#42-prometheus--grafana)
  - [Identifying Memory Leaks](#identifying-memory-leaks)
    - [Signs of Memory Leak](#signs-of-memory-leak)
    - [Common Causes](#common-causes)
    - [Example Analysis Output](#example-analysis-output)
  - [Best Practices](#best-practices)
  - [Quick Commands Reference](#quick-commands-reference)
  - [Troubleshooting](#troubleshooting)
    - [Cannot access pprof](#cannot-access-pprof)
    - [Graphviz error](#graphviz-error)
    - [Memory remains high after load test](#memory-remains-high-after-load-test)
  - [Additional Resources](#additional-resources)

## Method 1: Using pprof (Recommended)

### 1.1 Run the Application

```bash
make run-restapi
```

The application will run pprof endpoints on the same port as your REST API.

### 1.2 Access pprof Endpoints

**Important:** All pprof endpoints require authentication via `Authotization` header with the static token configured in `env.json`.

Using curl (replace `<PORT>` with your application port, default is usually 8080):

```bash
# Set your token from env.json
export PPROF_TOKEN="your-secret-token"

# Heap profile (memory allocation)
curl -H "Authotization: ${PPROF_TOKEN}" \
  http://localhost:<PORT>/debug/pprof/heap

# Goroutine profile
curl -H "Authotization: ${PPROF_TOKEN}" \
  http://localhost:<PORT>/debug/pprof/goroutine

# All available profiles
curl -H "Authotization: ${PPROF_TOKEN}" \
  http://localhost:<PORT>/debug/pprof/
```

**Notes:**

- pprof endpoints are available on the same port as your REST API, not on port 6060
- Authentication is required for security - configure `static_token` in your app's pprof config
- Header name is `Authotization` (matches code implementation)

### 1.3 Capture Heap Profile (Before Load Test)

```bash
# Capture baseline heap profile (authentication required)
curl -H "Authotization: your-secret-token" \
  http://localhost:<PORT>/debug/pprof/heap > heap_before.prof
```

**Note:** Replace `your-secret-token` with the actual token from `env.json` pprof config.

### 1.4 Generate Load Test

Use tools like `hey`, `ab`, or `wrk` to generate load:

```bash
# Install hey (if not already installed)
go install github.com/rakyll/hey@latest

# Generate 10000 requests with 100 concurrent connections
hey -n 10000 -c 100 http://localhost:<PORT>/api/v1/healthcheck
```

### 1.5 Capture Heap Profile (After Load Test)

```bash
# Wait a few seconds for GC to run
sleep 5

# Capture heap profile after load test (authentication required)
curl -H "Authotization: your-secret-token" \
  http://localhost:<PORT>/debug/pprof/heap > heap_after.prof
```

### 1.6 Analyze Memory Profile

#### Option A: Interactive Analysis

```bash
# Analyze with pprof interactive mode
go tool pprof heap_after.prof

# Commands in interactive mode:
# top    - Show top memory consumers
# list   - Show source code of function
# web    - Generate graph visualization (requires graphviz)
# pdf    - Generate PDF report (requires graphviz)
```

#### Option B: Compare Profiles

```bash
# Compare before and after
go tool pprof -base=heap_before.prof heap_after.prof

# In interactive mode, use:
# top    - Show functions with biggest memory increase
# list <function_name> - Show source code
```

#### Option C: Generate Visual Report

```bash
# Install graphviz (if not already installed)
brew install graphviz  # macOS
# or
sudo apt-get install graphviz  # Linux

# Generate graph
go tool pprof -png heap_after.prof > heap_graph.png
go tool pprof -pdf heap_after.prof > heap_report.pdf

# Generate comparison graph
go tool pprof -base=heap_before.prof -png heap_after.prof > heap_diff.png
```

### 1.7 Continuous Profiling

For monitoring real-time (authentication required):

```bash
# Set environment variable for auth header
export PPROF_TOKEN="your-secret-token"

# Open pprof web UI with custom fetcher
go tool pprof -http=:8080 \
  -fetch_timeout=30s \
  "http://localhost:<PORT>/debug/pprof/heap?token=${PPROF_TOKEN}"

# Alternative: Use curl to fetch then analyze
curl -H "Authotization: ${PPROF_TOKEN}" \
  http://localhost:<PORT>/debug/pprof/heap > heap.prof
go tool pprof -http=:8080 heap.prof
```

**Note:** go tool pprof doesn't support custom headers directly, so use curl to fetch profiles first or configure token in URL if your implementation supports it.

## Method 2: Check Goroutine Leaks

Goroutine leaks can also cause memory leaks.

### 2.1 Check Goroutine Count

```bash
# Before load test (authentication required)
curl -H "Authotization: your-secret-token" \
  http://localhost:<PORT>/debug/pprof/goroutine?debug=1 > goroutine_before.txt

# After load test and wait a few seconds
sleep 10
curl -H "Authotization: your-secret-token" \
  http://localhost:<PORT>/debug/pprof/goroutine?debug=1 > goroutine_after.txt

# Compare
diff goroutine_before.txt goroutine_after.txt
```

### 2.2 Analyze Goroutine Profile

```bash
# Capture goroutine profile (authentication required)
curl -H "Authotization: your-secret-token" \
  http://localhost:<PORT>/debug/pprof/goroutine > goroutine.prof

# Analyze
go tool pprof goroutine.prof

# In interactive mode:
# top    - Show goroutine creators
# traces - Show goroutine call stacks
```

## Method 3: Runtime Memory Stats

### 3.1 Create Monitoring Endpoint

Add an endpoint for monitoring memory stats (optional):

```go
// internal/transport/monitoring/restapi_monitoring.go
package transportmonitoring

import (
    "encoding/json"
    "net/http"
    "runtime"
)

type TransportMonitoringRestApi struct{}

func NewTransportRestApi() *TransportMonitoringRestApi {
    return &TransportMonitoringRestApi{}
}

func (t *TransportMonitoringRestApi) GetMemoryStats(w http.ResponseWriter, r *http.Request) {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    stats := map[string]interface{}{
        "alloc_mb":       bToMb(m.Alloc),
        "total_alloc_mb": bToMb(m.TotalAlloc),
        "sys_mb":         bToMb(m.Sys),
        "num_gc":         m.NumGC,
        "goroutines":     runtime.NumGoroutine(),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}
```

### 3.2 Monitor Memory Usage

```bash
# Watch memory stats every 2 seconds
watch -n 2 'curl -s http://localhost:<PORT>/api/v1/monitoring/memory'
```

## Method 4: Using External Tools

### 4.1 Continuous Profiler (Pyroscope)

```bash
# Install Pyroscope
docker run -it -p 4040:4040 pyroscope/pyroscope:latest server

# In your application, add pyroscope client
```

### 4.2 Prometheus + Grafana

Monitor memory metrics with prometheus:

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Expose prometheus metrics
http.Handle("/metrics", promhttp.Handler())
```

## Identifying Memory Leaks

### Signs of Memory Leak

1. **Heap memory continues to increase** after load test and doesn't decrease after GC
2. **Goroutine count increases** and doesn't return to baseline
3. **Object allocation** in certain functions keeps growing without being released
4. **RSS (Resident Set Size)** continues to increase in system monitor

### Common Causes

1. **Goroutine leaks:**
   - Goroutines not cleaned up
   - Missing context cancellation
   - Channels not closed

2. **Memory not released:**
   - Global variables that keep growing
   - Cache without eviction policy
   - Database connections not closed
   - HTTP response body not closed

3. **Circular references:**
   - Objects that reference each other and cannot be garbage collected

### Example Analysis Output

```bash
$ go tool pprof heap_after.prof

(pprof) top
Showing nodes accounting for 512.01MB, 100% of 512.01MB total
      flat  flat%   sum%        cum   cum%
  256.01MB 50.00% 50.00%   256.01MB 50.00%  database/sql.(*DB).addDep
  128.00MB 25.00% 75.00%   128.00MB 25.00%  internal/module/user/repository.(*repository).cache
  128.00MB 25.00% 100%      128.00MB 25.00%  net/http.(*persistConn).readLoop
```

In the example above:

- 50% memory used by database connections (possible connection leak)
- 25% memory for cache (might need eviction policy)
- 25% memory for HTTP connections (possibly keep-alive connections)

## Best Practices

1. **Run GC before capturing profile:**

   ```bash
   curl -H "Authotization: your-secret-token" \
     http://localhost:<PORT>/debug/pprof/heap?gc=1 > heap.prof
   ```

2. **Test under realistic load** - use load that resembles production

3. **Wait for warmup** - wait for application to stabilize before capturing baseline

4. **Compare multiple snapshots** - don't just compare 2 snapshots

5. **Check after idle period** - memory should decrease after load stops

6. **Monitor goroutines** - goroutine leaks can cause memory leaks

7. **Profile production safely** - pprof has minimal overhead but still monitor

## Quick Commands Reference

```bash
# Set authentication token
export PPROF_TOKEN="your-secret-token"

# Heap memory profile (fetch with auth, then analyze)
curl -H "Authotization: ${PPROF_TOKEN}" \
  http://localhost:<PORT>/debug/pprof/heap > heap.prof
go tool pprof heap.prof

# Goroutine profile
curl -H "Authotization: ${PPROF_TOKEN}" \
  http://localhost:<PORT>/debug/pprof/goroutine > goroutine.prof
go tool pprof goroutine.prof

# Allocs (all allocations)
curl -H "Authotization: ${PPROF_TOKEN}" \
  http://localhost:<PORT>/debug/pprof/allocs > allocs.prof
go tool pprof allocs.prof

# Compare two profiles
go tool pprof -base=before.prof after.prof

# Web UI (use pre-fetched profile)
go tool pprof -http=:8080 heap.prof

# Generate PNG
go tool pprof -png heap.prof > output.png

# Capture with GC first
curl -H "Authotization: ${PPROF_TOKEN}" \
  http://localhost:<PORT>/debug/pprof/heap?gc=1 > heap.prof
```

## Troubleshooting

### Cannot access pprof

Ensure the import has been added:

```go
import _ "net/http/pprof"
```

### Graphviz error

Install graphviz:

```bash
# macOS
brew install graphviz

# Ubuntu/Debian
sudo apt-get install graphviz
```

### Memory remains high after load test

Normal behavior if:

- Go runtime caching memory for reuse
- GOGC environment variable affects GC frequency
- See `HeapInuse` vs `HeapSys` in memstats

Force GC for testing:

```bash
curl http://localhost:<PORT>/debug/pprof/heap?gc=1
```

## Additional Resources

- [Go pprof documentation](https://pkg.go.dev/runtime/pprof)
- [Profiling Go Programs](https://go.dev/blog/pprof)
- [Go Memory Management](https://go.dev/doc/diagnostics)
