package metrics

import (
	"log"
	"strconv"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/prometheus/client_golang/prometheus"
)

type powerUsageCollector struct {
	powerUsed *prometheus.Desc
	powerCap  *prometheus.Desc
}

func NewPowerUsageCollector() *powerUsageCollector {
	return &powerUsageCollector{
		powerUsed: prometheus.NewDesc(
			"gpu_power_usage_watts",
			"Current GPU power usage in watts",
			[]string{"name", "gpu", "gpu_idx"},
			nil,
		),
		powerCap: prometheus.NewDesc(
			"gpu_power_cap_watts",
			"Maximum GPU power usage in watts",
			[]string{"name", "gpu", "gpu_idx"},
			nil,
		),
	}
}

func (collector *powerUsageCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.powerUsed
	ch <- collector.powerCap
}

// Collect is called by the Prometheus registry when collecting metrics
func (collector *powerUsageCollector) Collect(ch chan<- prometheus.Metric) {
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

		power, err := device.GetPowerUsage()
		if err == nvml.SUCCESS {
			ch <- prometheus.MustNewConstMetric(
				collector.powerUsed,
				prometheus.GaugeValue,
				float64(power)/1000,
				name,
				uuid,
				strconv.Itoa(i),
			)
		}

		powerCap, err := device.GetPowerManagementLimit()
		if err == nvml.SUCCESS {
			ch <- prometheus.MustNewConstMetric(
				collector.powerCap,
				prometheus.GaugeValue,
				float64(powerCap)/1000,
				name,
				uuid,
				strconv.Itoa(i),
			)
		}
	}
}
