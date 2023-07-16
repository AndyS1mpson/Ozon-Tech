package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	histogram *prometheus.HistogramVec
	counter   *prometheus.CounterVec
}

var (
	prometheusMetrics *metrics
)

func init() {
	reg := prometheus.NewRegistry()
	prometheusMetrics = New(reg)
}

func New(reg prometheus.Registerer) *metrics {
	m := &metrics{
		histogram: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "loms",
				Name:      "histogram_response_time_seconds",
				Buckets:   []float64{},
			},
			[]string{
				"status",
				"method",
			},
		),
		counter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "loms",
				Name:      "counter_requests_total",
			},
			[]string{
				"type",
			},
		),
	}

	reg.MustRegister(m.histogram, m.counter)
	return m
}
