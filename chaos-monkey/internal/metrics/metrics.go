package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ActionsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "chaos_actions_total",
		Help: "Total number of chaos actions performed",
	}, []string{"action"})

	ActionErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "chaos_action_errors_total",
		Help: "Total number of chaos action errors",
	}, []string{"action"})

	ActionDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "chaos_action_duration_seconds",
		Help:    "Duration of chaos actions",
		Buckets: prometheus.DefBuckets,
	}, []string{"action"})

	LastActionTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "chaos_last_action_timestamp",
		Help: "Unix timestamp of the last chaos action by type",
	}, []string{"action"})
)
