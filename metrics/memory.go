package metrics

import (
	"runtime"
)

type MemoryReport struct {
	Goroutines     int     `json:"goroutines"`
	HeapAllocated  float64 `json:"heapAllocated"`
	HeapInUse      float64 `json:"heapInUse"`
	StackAllocated float64 `json:"stackAllocated"`
	StackInUse     float64 `json:"stackInUse"`
	TotalAllocated float64 `json:"totalAllocated"`
	TotalInUse     float64 `json:"totalInUse"`
	MemoryUnit     string  `json:"memoryUnit"`
}

const megabyte = 1048576.0

func GetMemoryReport() MemoryReport {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return MemoryReport{
		Goroutines:     runtime.NumGoroutine(),
		HeapAllocated:  float64(m.HeapSys-m.HeapReleased) / megabyte,
		HeapInUse:      float64(m.HeapInuse) / megabyte,
		StackAllocated: float64(m.StackSys) / megabyte,
		StackInUse:     float64(m.StackInuse) / megabyte,
		TotalAllocated: float64(m.HeapSys+m.StackSys-m.HeapReleased) / megabyte,
		TotalInUse:     float64(m.HeapInuse+m.StackInuse) / megabyte,
		MemoryUnit:     "MB",
	}
}
