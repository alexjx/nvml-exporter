package metrics

import (
	"log"
	"strconv"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/prometheus/client_golang/prometheus"
)

type gpuUsageCollector struct {
	gpuUsage *prometheus.Desc
}

func NewGPUUsageCollector() *gpuUsageCollector {
	return &gpuUsageCollector{
		gpuUsage: prometheus.NewDesc(
			"gpu_usage_percentage",
			"Current GPU usage as a percentage",
			[]string{"name", "gpu", "gpu_idx"},
			nil,
		),
	}
}

func (collector *gpuUsageCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.gpuUsage
}

// Collect is called by the Prometheus registry when collecting metrics
func (collector *gpuUsageCollector) Collect(ch chan<- prometheus.Metric) {
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

		usage, err := device.GetUtilizationRates()
		if err == nvml.SUCCESS {
			ch <- prometheus.MustNewConstMetric(
				collector.gpuUsage,
				prometheus.GaugeValue,
				float64(usage.Gpu),
				name,
				uuid,
				strconv.Itoa(i),
			)
		}
	}
}
