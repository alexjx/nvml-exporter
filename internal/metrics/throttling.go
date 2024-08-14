package metrics

import (
	"log"
	"strconv"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/prometheus/client_golang/prometheus"
)

type throttleReasonsCollector struct {
	throttleReasons *prometheus.Desc
}

func NewThrottleReasonsCollector() *throttleReasonsCollector {
	return &throttleReasonsCollector{
		throttleReasons: prometheus.NewDesc(
			"gpu_throttle_reasons",
			"Reasons for GPU throttling (1 if active, 0 if not)",
			[]string{"name", "gpu", "gpu_idx", "reason"},
			nil,
		),
	}
}

func (collector *throttleReasonsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.throttleReasons
}

func (collector *throttleReasonsCollector) Collect(ch chan<- prometheus.Metric) {
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

		reasons, err := device.GetCurrentClocksThrottleReasons()
		if err == nvml.SUCCESS {
			collector.collectThrottleReason(ch, name, uuid, strconv.Itoa(i), "idle", reasons&nvml.ClocksThrottleReasonGpuIdle != 0)
			collector.collectThrottleReason(ch, name, uuid, strconv.Itoa(i), "power", reasons&nvml.ClocksThrottleReasonSwPowerCap != 0)
			collector.collectThrottleReason(ch, name, uuid, strconv.Itoa(i), "thermal", reasons&nvml.ClocksThrottleReasonHwThermalSlowdown != 0)
			// Add more reasons as needed
		}
	}
}

func (collector *throttleReasonsCollector) collectThrottleReason(ch chan<- prometheus.Metric, name, uuid, gpuIdx, reason string, active bool) {
	ch <- prometheus.MustNewConstMetric(
		collector.throttleReasons,
		prometheus.GaugeValue,
		boolToFloat(active),
		name,
		uuid,
		gpuIdx,
		reason,
	)
}

func boolToFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
