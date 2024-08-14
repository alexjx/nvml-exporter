package metrics

import (
	"log"
	"strconv"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/prometheus/client_golang/prometheus"
)

type fanSpeedCollector struct {
	fanSpeed *prometheus.Desc
}

func NewFanSpeedCollector() *fanSpeedCollector {
	return &fanSpeedCollector{
		fanSpeed: prometheus.NewDesc(
			"gpu_fan_speed_percentage",
			"Current GPU fan speed as a percentage of maximum speed",
			[]string{"name", "gpu", "gpu_idx"},
			nil,
		),
	}
}

func (collector *fanSpeedCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.fanSpeed
}

func (collector *fanSpeedCollector) Collect(ch chan<- prometheus.Metric) {
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

		speed, ret := device.GetFanSpeed()
		if ret == nvml.SUCCESS {
			ch <- prometheus.MustNewConstMetric(
				collector.fanSpeed,
				prometheus.GaugeValue,
				float64(speed),
				name,
				uuid,
				strconv.Itoa(i),
			)
		}
	}
}
