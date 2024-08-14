package metrics

import (
	"log"
	"strconv"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/prometheus/client_golang/prometheus"
)

type eccErrorsCollector struct {
	eccErrors *prometheus.Desc
}

func NewECCErrorsCollector() *eccErrorsCollector {
	return &eccErrorsCollector{
		eccErrors: prometheus.NewDesc(
			"gpu_ecc_errors",
			"Number of ECC errors detected",
			[]string{"name", "gpu", "gpu_idx", "error_type"},
			nil,
		),
	}
}

func (collector *eccErrorsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.eccErrors
}

func (collector *eccErrorsCollector) Collect(ch chan<- prometheus.Metric) {
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

		eccCounts, err := device.GetTotalEccErrors(nvml.MEMORY_ERROR_TYPE_UNCORRECTED, nvml.VOLATILE_ECC)
		if err == nvml.SUCCESS {
			ch <- prometheus.MustNewConstMetric(
				collector.eccErrors,
				prometheus.GaugeValue,
				float64(eccCounts),
				name,
				uuid,
				strconv.Itoa(i),
				"uncorrected",
			)
		}
		eccCounts, err = device.GetTotalEccErrors(nvml.MEMORY_ERROR_TYPE_CORRECTED, nvml.VOLATILE_ECC)
		if err == nvml.SUCCESS {
			ch <- prometheus.MustNewConstMetric(
				collector.eccErrors,
				prometheus.GaugeValue,
				float64(eccCounts),
				name,
				uuid,
				strconv.Itoa(i),
				"corrected",
			)
		}
	}
}
