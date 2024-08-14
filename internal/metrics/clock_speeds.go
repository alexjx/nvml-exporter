package metrics

import (
	"log"
	"strconv"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/prometheus/client_golang/prometheus"
)

type clockCollector struct {
	gpuClock *prometheus.Desc
	memClock *prometheus.Desc
}

func NewClockCollector() *clockCollector {
	return &clockCollector{
		gpuClock: prometheus.NewDesc(
			"gpu_clock_speed_mhz",
			"Current GPU clock speed in MHz",
			[]string{"name", "gpu", "gpu_idx"},
			nil,
		),
		memClock: prometheus.NewDesc(
			"gpu_memory_clock_speed_mhz",
			"Current GPU memory clock speed in MHz",
			[]string{"name", "gpu", "gpu_idx"},
			nil,
		),
	}
}

func (c *clockCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.gpuClock
	ch <- c.memClock
}

func (collector *clockCollector) Collect(ch chan<- prometheus.Metric) {
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

		clock, err := device.GetClockInfo(nvml.CLOCK_GRAPHICS)
		if err == nvml.SUCCESS {
			ch <- prometheus.MustNewConstMetric(collector.gpuClock, prometheus.GaugeValue, float64(clock), name, uuid, strconv.Itoa(i))
		}

		memClockSpeed, err := device.GetClockInfo(nvml.CLOCK_MEM)
		if err == nvml.SUCCESS {
			ch <- prometheus.MustNewConstMetric(collector.memClock, prometheus.GaugeValue, float64(memClockSpeed), name, uuid, strconv.Itoa(i))
		}
	}
}
