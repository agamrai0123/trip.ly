// Package database provides a pgx/v5 connection pool factory and golang-migrate runner.
package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

// PgxPoolCollector is a prometheus.Collector that reports pgxpool.Pool statistics.
// Register it with your service's Prometheus registry to expose active DB connection counts.
type PgxPoolCollector struct {
	pool        *pgxpool.Pool
	serviceName string
	acquired    *prometheus.Desc
	idle        *prometheus.Desc
	total       *prometheus.Desc
	maxConns    *prometheus.Desc
}

// NewPoolCollector returns a PgxPoolCollector for the given pool and service name.
// The serviceName is embedded as a constant label so the metric can be filtered in Grafana.
func NewPoolCollector(pool *pgxpool.Pool, serviceName string) *PgxPoolCollector {
	constLabels := prometheus.Labels{"service": serviceName}
	return &PgxPoolCollector{
		pool:        pool,
		serviceName: serviceName,
		acquired: prometheus.NewDesc(
			"db_pool_acquired_connections",
			"Number of currently acquired (in-use) database connections.",
			nil, constLabels,
		),
		idle: prometheus.NewDesc(
			"db_pool_idle_connections",
			"Number of idle (ready) database connections in the pool.",
			nil, constLabels,
		),
		total: prometheus.NewDesc(
			"db_pool_total_connections",
			"Total number of open database connections (acquired + idle + constructing).",
			nil, constLabels,
		),
		maxConns: prometheus.NewDesc(
			"db_pool_max_connections",
			"Maximum number of database connections the pool is allowed to open.",
			nil, constLabels,
		),
	}
}

// Describe sends descriptor metadata to the Prometheus registry.
func (c *PgxPoolCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.acquired
	ch <- c.idle
	ch <- c.total
	ch <- c.maxConns
}

// Collect reads current pool statistics and emits gauge metrics.
func (c *PgxPoolCollector) Collect(ch chan<- prometheus.Metric) {
	stat := c.pool.Stat()
	ch <- prometheus.MustNewConstMetric(c.acquired, prometheus.GaugeValue, float64(stat.AcquiredConns()))
	ch <- prometheus.MustNewConstMetric(c.idle, prometheus.GaugeValue, float64(stat.IdleConns()))
	ch <- prometheus.MustNewConstMetric(c.total, prometheus.GaugeValue, float64(stat.TotalConns()))
	ch <- prometheus.MustNewConstMetric(c.maxConns, prometheus.GaugeValue, float64(stat.MaxConns()))
}
