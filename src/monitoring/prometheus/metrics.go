package prometheus

import "github.com/prometheus/client_golang/prometheus"

// Define future metrics below
var (
	containerCPUUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "docker",
			Name:      "container_cpu_usage_seconds_total",
			Help:      "Total CPU time consumed by containers in seconds",
		},
		[]string{"container_id"},
	)
	containerMemUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "docker",
			Name:      "container_memory_usage_bytes",
			Help:      "Current memory usage of containers in bytes",
		},
		[]string{"container_id"},
	)
)

// RegisterPrometheusMetrics registers Prometheus metrics for monitoring
func RegisterPrometheusMetrics() {
	prometheus.MustRegister(containerCPUUsage)
	prometheus.MustRegister(containerMemUsage)
	// Register additional metrics here if needed
}
