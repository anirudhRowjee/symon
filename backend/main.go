package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// Define the standard tick rate of 500 ms
const TICK_RATE = time.Millisecond * 500
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

type metricsDashboard struct {

	// --- CPU Information
	CPU_counts           int       `json:"cpu_count"`            // number of CPUs present
	CPU_cores            int       `json:"cpu_cores"`            // number of cores
	CPU_usage_percentage []float64 `json:"cpu_usage_percentage"` // usage statistics

	// --- Memory Information
	MEM_usage_percentage float64 `json:"mem_usage_percentage"` // usage statistics
	MEM_total_available  int     `json:"mem_total_available"`  // all the available memory
	MEM_total_used       int     `json:"mem_total_used"`       // how much of it is used
}

func main() {

	fmt.Println(BANNER)

	// TODO Declare shared state

	// TODO Make shared state accessible via a weeb server

	// TODO Send this off in a goroutine
	for {

		time.Sleep(time.Millisecond * 500)

		c, _ := cpu.Counts(false)
		v, _ := mem.VirtualMemory()
		cpu_usage_percentage, _ := cpu.Percent(0, true)

		log.Println("Usage Statistics =>")
		fmt.Printf("Memory Usage -> %f%%\n", v.UsedPercent)
		fmt.Printf("CPU Count -> %d\n", c)
		fmt.Println("CPU Usage -> ", cpu_usage_percentage)

	}

}
