package metrics

import (
	"runtime"
	"strings"
)

var sizeKB float64 = 1 << (10 * 1)
var sizeMB float64 = 1 << (10 * 2)
var sizeGB float64 = 1 << (10 * 3)

//MemoryReport Contrato de Memory Check
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

//GetMemoryReport Obtém dados de Memory Check
func GetMemoryReport(u string) MemoryReport {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	var unit float64

	switch strings.ToUpper(u) {
	case "GB":
		unit = sizeGB
		break
	case "KB":
		unit = sizeKB
	default:
		u = "MB"
		unit = sizeMB
	}

	return MemoryReport{
		Goroutines:     runtime.NumGoroutine(),
		HeapAllocated:  float64(m.HeapSys-m.HeapReleased) / unit,
		HeapInUse:      float64(m.HeapInuse) / unit,
		StackAllocated: float64(m.StackSys) / unit,
		StackInUse:     float64(m.StackInuse) / unit,
		TotalAllocated: float64(m.HeapSys+m.StackSys-m.HeapReleased) / unit,
		TotalInUse:     float64(m.HeapInuse+m.StackInuse) / unit,
		MemoryUnit:     strings.ToUpper(u),
	}
}
