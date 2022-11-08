package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// Define the standard tick rate of 500 ms
const TICK_RATE = time.Millisecond * 750
const BANNER = `
  /$$$$$$            /$$      /$$                                     
 /$$__  $$          | $$$    /$$$                                     
| $$  \__/ /$$   /$$| $$$$  /$$$$  /$$$$$$  /$$$$$$$                  
|  $$$$$$ | $$  | $$| $$ $$/$$ $$ /$$__  $$| $$__  $$                 
 \____  $$| $$  | $$| $$  $$$| $$| $$  \ $$| $$  \ $$                 
 /$$  \ $$| $$  | $$| $$\  $ | $$| $$  | $$| $$  | $$                 
|  $$$$$$/|  $$$$$$$| $$ \/  | $$|  $$$$$$/| $$  | $$                 
 \______/  \____  $$|__/     |__/ \______/ |__/  |__/                 
           /$$  | $$                                                  
          |  $$$$$$/                                                  
           \______/                                                   
 /$$                           /$$                                 /$$
| $$                          | $$                                | $$
| $$$$$$$   /$$$$$$   /$$$$$$$| $$   /$$  /$$$$$$  /$$$$$$$   /$$$$$$$
| $$__  $$ |____  $$ /$$_____/| $$  /$$/ /$$__  $$| $$__  $$ /$$__  $$
| $$  \ $$  /$$$$$$$| $$      | $$$$$$/ | $$$$$$$$| $$  \ $$| $$  | $$
| $$  | $$ /$$__  $$| $$      | $$_  $$ | $$_____/| $$  | $$| $$  | $$
| $$$$$$$/|  $$$$$$$|  $$$$$$$| $$ \  $$|  $$$$$$$| $$  | $$|  $$$$$$$
|_______/  \_______/ \_______/|__/  \__/ \_______/|__/  |__/ \_______/
`

// Define the shape of data being shared between the backend
// and the frontend

type CpuMetrics struct {
	// --- CPU Information
	CPU_counts           int       `json:"cpu_count"`            // number of CPUs present
	CPU_cores            int       `json:"cpu_cores"`            // number of cores
	CPU_usage_percentage []float64 `json:"cpu_usage_percentage"` // usage statistics
}

func (c *CpuMetrics) new() {
	// Populate new Memory Metrics
	cpuStatStruct, _ := cpu.Info()
	cpunum, _ := cpu.Counts(false)
	cpucores := len(cpuStatStruct)
	cpu_usage_percentage, _ := cpu.Percent(0, true)

	c.CPU_cores = cpucores
	c.CPU_counts = cpunum
	c.CPU_usage_percentage = cpu_usage_percentage
}

type MemoryMetrics struct {
	// --- Memory Information
	MEM_usage_percentage float64 `json:"mem_usage_percentage"` // usage statistics
	MEM_total_available  int     `json:"mem_total_available"`  // all the available memory
	MEM_total_used       int     `json:"mem_total_used"`       // how much of it is used
}

func (m *MemoryMetrics) new() {
	// Populate new Memory Metrics
	v, _ := mem.VirtualMemory()
	m.MEM_usage_percentage = v.UsedPercent
	m.MEM_total_used = int(v.Used)
	m.MEM_total_available = int(v.Available)
}

type systemMetrics struct {
	// timestamped metrics to export
	Current_timestamp      string        `json:"metrics_timestamp"`
	Current_CPU_Metrics    CpuMetrics    `json:"cpu_metrics"`
	Current_Memory_Metrics MemoryMetrics `json:"memory_metrics"`
}

// Function to do timestamped udpate
// No pointers here because we want to copy, not read reference
func (s *systemMetrics) update(cpu CpuMetrics, memory MemoryMetrics) {

	// current time
	time := time.Now().Local().Format("Mon Jan 2 15:04:05 MST 2006")
	s.Current_timestamp = time

	// Format the timestamp
	s.Current_CPU_Metrics = cpu
	s.Current_Memory_Metrics = memory
}

// Container Struct for Metrics with Timestamp
type MetricsDashboard struct {
	mutex                 sync.Mutex // Concurrency Control
	waitgroup             sync.WaitGroup
	latest_system_metrics systemMetrics // State
	cpu_metrics_chan      chan CpuMetrics
	memory_metrics_chan   chan MemoryMetrics
	shutdown_chan         chan int
	tick_rate             time.Duration
}

func (m *MetricsDashboard) init(sleep time.Duration) {
	m.mutex = sync.Mutex{}
	m.waitgroup = sync.WaitGroup{}
	m.latest_system_metrics = systemMetrics{}
	m.cpu_metrics_chan = make(chan CpuMetrics)
	m.memory_metrics_chan = make(chan MemoryMetrics)
	m.shutdown_chan = make(chan int)
	m.tick_rate = sleep
}

func (m *MetricsDashboard) spawn_state_manager() {
	// function to spawn the state manager goroutine
	// Run this infinite loop
	for {
		select {
		case latest_cpu_metrics_copy := <-m.cpu_metrics_chan:
			// update CPU metrics
			m.mutex.Lock()
			m.latest_system_metrics.Current_CPU_Metrics = latest_cpu_metrics_copy
			m.mutex.Unlock()
		case latest_memory_metrics_copy := <-m.memory_metrics_chan:
			// update Memory metrics
			m.mutex.Lock()
			m.latest_system_metrics.Current_Memory_Metrics = latest_memory_metrics_copy
			m.mutex.Unlock()
			// TODO Implement with errorgroup context https://www.fullstory.com/blog/why-errgroup-withcontext-in-golang-server-handlers/
		}
	}
}

func (m *MetricsDashboard) serve_metrics() {
	// function to spawn the Metrics server
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {

		// read the metric
		m.mutex.Lock()
		response := m.latest_system_metrics
		m.mutex.Unlock()

		// serialize
		w.Header().Set("Content-Type", "application/json")
		json_output, err := json.Marshal(response)
		if err != nil {
			w.Write([]byte("Could not marshal"))
		}
		w.Write(json_output)
	})
	http.ListenAndServe("127.0.0.1:1337", mux)
}

func (m *MetricsDashboard) spawn_fetchloops() {
	// function to spawn the state manager goroutine

	// Memory fetch loop
	go func() {
		for {

			// Load Metrics at this time
			metrics := MemoryMetrics{}
			metrics.new()

			// send this into the channel
			m.memory_metrics_chan <- metrics

			// Sleep
			time.Sleep(m.tick_rate)
		}
	}()

	// CPU fetch loop
	go func() {
		for {

			// Load Metrics at this time
			metrics := CpuMetrics{}
			metrics.new()

			// send this into the channel
			m.cpu_metrics_chan <- metrics

			// Sleep
			time.Sleep(m.tick_rate)
		}
	}()
}

// Function to start the entire infrastructure
func (m *MetricsDashboard) start() {
	// TODO Graceful shutdown
	// TODO Maybe use a waitgroup

	// spawn state manager goroutine
	go m.spawn_state_manager()

	// spawn metric fetch loop goroutines
	m.spawn_fetchloops()

	// spawn server goroutines
	m.serve_metrics()

}

func main() {

	fmt.Print(BANNER)

	m := MetricsDashboard{}
	m.init(TICK_RATE)

	m.start()
}
