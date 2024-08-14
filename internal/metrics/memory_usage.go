package metrics

import (
	"log"
	"strconv"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/prometheus/client_golang/prometheus"
)

type gpuMemoryUsageCollector struct {
	gpuMemoryUsed  *prometheus.Desc
	gpuMemoryTotal *prometheus.Desc
}

func NewGPUMemoryUsageCollector() *gpuMemoryUsageCollector {
	return &gpuMemoryUsageCollector{
		gpuMemoryUsed: prometheus.NewDesc(
			"gpu_memory_used_bytes",
			"Current GPU memory usage in bytes",
			[]string{"name", "gpu", "gpu_idx"},
			nil,
		),
		gpuMemoryTotal: prometheus.NewDesc(
			"gpu_memory_total_bytes",
			"Total GPU memory in bytes",
			[]string{"name", "gpu", "gpu_idx"},
			nil,
		),
	}
}

func (collector *gpuMemoryUsageCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.gpuMemoryUsed
	ch <- collector.gpuMemoryTotal
}

func (collector *gpuMemoryUsageCollector) Collect(ch chan<- prometheus.Metric) {
	deviceCount, err := nvml.DeviceGetCount()
	if err != nvml.SUCCESS {
		log.Fatalf("Failed to get device count: %v", err)
	}

	for i := 0; i < deviceCount; i++ {
		device, err := nvml.DeviceGetHandleByIndex(i)
		if err != nvml.SUCCESS {
			log.Printf("Failed to get handle for device %d: %v", i, err)
			continue
		}
		uuid, ret := device.GetUUID()
		if ret != nvml.SUCCESS {
			log.Printf("Failed to get UUID for device %d: %v", i, ret)
			continue
		}
		name, err := device.GetName()
		if err != nvml.SUCCESS {
			log.Printf("Failed to get name for device %d: %v", i, err)
			continue
		}

		memoryInfo, ret := device.GetMemoryInfo()
		if ret == nvml.SUCCESS {
			ch <- prometheus.MustNewConstMetric(
				collector.gpuMemoryUsed,
				prometheus.GaugeValue,
				float64(memoryInfo.Used),
				name,
				uuid,
				strconv.Itoa(i),
			)
			ch <- prometheus.MustNewConstMetric(
				collector.gpuMemoryTotal,
				prometheus.GaugeValue,
				float64(memoryInfo.Total),
				name,
				uuid,
				strconv.Itoa(i),
			)
		}
	}
}
