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

func main() {

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
