package metrics

import (
	"log"
	"strconv"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/prometheus/client_golang/prometheus"
)

type gpuTemperatureCollector struct {
	gpuTemp *prometheus.Desc
}

func NewGPUTemperatureCollector() *gpuTemperatureCollector {
	return &gpuTemperatureCollector{
		gpuTemp: prometheus.NewDesc(
			"gpu_temperature_celsius",
			"Current GPU temperature in degrees Celsius",
			[]string{"name", "gpu", "gpu_idx"},
			nil,
		),
	}
}

func (collector *gpuTemperatureCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.gpuTemp
}

func (collector *gpuTemperatureCollector) Collect(ch chan<- prometheus.Metric) {
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

		temp, err := device.GetTemperature(nvml.TEMPERATURE_GPU)
		if err == nvml.SUCCESS {
			ch <- prometheus.MustNewConstMetric(
				collector.gpuTemp,
				prometheus.GaugeValue,
				float64(temp),
				name,
				uuid,
				strconv.Itoa(i),
			)
		}
	}
}
