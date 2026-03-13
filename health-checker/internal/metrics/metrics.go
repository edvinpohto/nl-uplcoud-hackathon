package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ProbeSuccess = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "uptime_probe_success",
		Help: "1 if the last probe was successful, 0 otherwise",
	})

	ResponseTime = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "probe_response_seconds",
		Help:    "HTTP response time of the victim-app health endpoint",
		Buckets: prometheus.DefBuckets,
	})

	ProbeErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "probe_errors_total",
		Help: "Total number of failed probes",
	})

	ProbeTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "probe_total",
		Help: "Total number of probes performed",
	})

	ConsecutiveFailures = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "probe_consecutive_failures",
		Help: "Current number of consecutive probe failures",
	})
)
