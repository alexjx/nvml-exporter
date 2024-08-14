package metrics

import (
	"log"
	"strconv"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/prometheus/client_golang/prometheus"
)

type pcieThroughputCollector struct {
	pcieThroughput *prometheus.Desc
}

func NewPCIeThroughputCollector() *pcieThroughputCollector {
	return &pcieThroughputCollector{
		pcieThroughput: prometheus.NewDesc(
			"gpu_pcie_throughput",
			"Current PCIe throughput in bytes per second",
			[]string{"name", "gpu", "gpu_idx", "direction"},
			nil,
		),
	}
}

func (collector *pcieThroughputCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.pcieThroughput
}

func (collector *pcieThroughputCollector) Collect(ch chan<- prometheus.Metric) {
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

		txBytes, ret := device.GetPcieThroughput(nvml.PCIE_UTIL_TX_BYTES)
		if ret == nvml.SUCCESS {
			ch <- prometheus.MustNewConstMetric(
				collector.pcieThroughput,
				prometheus.GaugeValue,
				float64(txBytes),
				name,
				uuid,
				strconv.Itoa(i),
				"tx",
			)
		}

		rxBytes, ret := device.GetPcieThroughput(nvml.PCIE_UTIL_RX_BYTES)
		if ret == nvml.SUCCESS {
			ch <- prometheus.MustNewConstMetric(
				collector.pcieThroughput,
				prometheus.GaugeValue,
				float64(rxBytes),
				name,
				uuid,
				strconv.Itoa(i),
				"rx",
			)
		}
	}
}
