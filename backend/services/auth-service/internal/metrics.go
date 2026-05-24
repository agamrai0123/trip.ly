package internal

import "github.com/prometheus/client_golang/prometheus"

// NewRegistry creates a new Prometheus registry with default process + Go collectors.
func NewRegistry() *prometheus.Registry {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)
	return reg
}
