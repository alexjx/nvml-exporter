package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

func RegisterMetrics() {
	// Register your metrics here
	prometheus.MustRegister(
		NewClockCollector(),
		NewECCErrorsCollector(),
		NewFanSpeedCollector(),
		NewGPUUsageCollector(),
		NewGPUMemoryUsageCollector(),
		NewPCIeThroughputCollector(),
		NewPowerUsageCollector(),
		NewGPUTemperatureCollector(),
		NewThrottleReasonsCollector(),
	)
}
