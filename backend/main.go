package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
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
	// Populate new CPU Metrics
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
	mu              sync.Mutex
	current_metrics systemMetrics
}

func (m *MetricsDashboard) spawn_state_manager() {
	// function to spawn the state manager goroutine
}

func (m *MetricsDashboard) serve_metrics() {
	// function to spawn the state manager goroutine
}

func (m *MetricsDashboard) spawn_fetchloops() {
	// function to spawn the state manager goroutine
}

// Function to start the entire infrastructure
func (m *MetricsDashboard) start(tick_rate int) {
	// TODO Graceful shutdown
	// TODO Maybe use a waitgroup

	// spawn comms channels

	// spawn state manager goroutine
	go m.spawn_state_manager()

	// spawn metric fetch loop goroutines
	go m.spawn_fetchloops()

	// spawn server goroutines
	go m.serve_metrics()

}

func main() {

	fmt.Print(BANNER)

	// TODO Declare shared state

	// TODO Make shared state accessible via a weeb server

	// TODO Send this off in a goroutine
	for {

		time.Sleep(time.Millisecond * 500)

		c, _ := cpu.Counts(false)
		cpu_usage_percentage, _ := cpu.Percent(0, true)

		log.Println("Usage Statistics =>")
		fmt.Printf("Memory Usage -> %f%%\n", v.UsedPercent)
		fmt.Printf("CPU Count -> %d\n", c)
		fmt.Println("CPU Usage -> ", cpu_usage_percentage)

	}

}
