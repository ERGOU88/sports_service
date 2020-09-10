package pprof

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
)

type Handle func() []byte

var handlers map[string]Handle

func init() {
	handlers = make(map[string]Handle)
	handlers["mem"] = memStats
	handlers["go"] = goStats
}

func Start(bind string) error {
	httpServeMux := http.NewServeMux()
	httpServeMux.HandleFunc("/stat", handle)
	fmt.Println("start  ")
	go func() {
		if err := http.ListenAndServe(bind, httpServeMux); err != nil {
			fmt.Printf("http.ListenAdServe(\"%s\") error(%v)", bind, err)
			os.Exit(1)
		}
	}()
	return nil
}

// memory stats
func memStats() []byte {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	// general
	res := map[string]interface{}{}
	res["alloc"] = m.Alloc
	res["total_alloc"] = m.TotalAlloc
	res["sys"] = m.Sys
	res["lookups"] = m.Lookups
	res["mallocs"] = m.Mallocs
	res["frees"] = m.Frees
	// heap
	res["heap_alloc"] = m.HeapAlloc
	res["heap_sys"] = m.HeapSys
	res["heap_idle"] = m.HeapIdle
	res["heap_inuse"] = m.HeapInuse
	res["heap_released"] = m.HeapReleased
	res["heap_objects"] = m.HeapObjects
	// low-level fixed-size struct alloctor
	res["stack_inuse"] = m.StackInuse
	res["stack_sys"] = m.StackSys
	res["mspan_inuse"] = m.MSpanInuse
	res["mspan_sys"] = m.MSpanSys
	res["mcache_inuse"] = m.MCacheInuse
	res["mcache_sys"] = m.MCacheSys
	res["buckhash_sys"] = m.BuckHashSys
	// GC
	res["next_gc"] = m.NextGC
	res["last_gc"] = m.LastGC
	res["pause_total_ns"] = m.PauseTotalNs
	//res["pause_ns"] = m.PauseNs
	res["num_gc"] = m.NumGC
	res["enable_gc"] = m.EnableGC
	res["gc_cpu"] = m.GCCPUFraction
	res["debug_gc"] = m.DebugGC
	//res["by_size"] = m.BySize
	return jsonRes(res)
}

// golang stats
func goStats() []byte {
	res := map[string]interface{}{}
	res["compiler"] = runtime.Compiler
	res["arch"] = runtime.GOARCH
	res["os"] = runtime.GOOS
	res["max_procs"] = runtime.GOMAXPROCS(-1)
	res["root"] = runtime.GOROOT()
	res["cgo_call"] = runtime.NumCgoCall()
	res["goroutine_num"] = runtime.NumGoroutine()
	res["version"] = runtime.Version()
	return jsonRes(res)
}

func jsonRes(res interface{}) []byte {
	byteJson, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		fmt.Printf("json.MarshalIndent(\"%v\", \"\", \"    \") error(%v)", res, err)
		return nil
	}
	return byteJson
}

// StatHandle get stat info by http
func handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	params := r.URL.Query()
	types := params.Get("type")
	h, ok := handlers[types]

	if !ok {
		http.Error(w, "Not Found", 404)
		return
	}

	res := h()
	if res != nil {
		if _, err := w.Write(res); err != nil {
			fmt.Printf("w.Write(\"%s\") error(%v)", string(res), err)
		}
	}
}
