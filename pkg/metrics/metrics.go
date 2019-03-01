package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	AlertsUpdate = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "alerts_update",
		Namespace: "alerterator",
		Help:      "number of alert synchronization performed",
	})
	AlertsApplied = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "alerts_applied",
		Namespace: "alerterator",
		Help:      "number of nais.io.Alert resources that have been applied",
	})
	AlertsDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "alerts_deleted",
		Namespace: "alerterator",
		Help:      "number of nais.io.Alert resources that have been deleted",
	})
	AlertsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "alerts_processed",
		Namespace: "alerterator",
		Help:      "number of nais.io.Alert resources that have been processed",
	})
	AlertsFailed = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "alerts_failed",
		Namespace: "alerterator",
		Help:      "number of nais.io.Alert resources that failed processing",
	})
	AlertsFailedEvent = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "alerts_failed_event",
		Namespace: "alerterator",
		Help:      "number of events that have failed",
	})
)

func init() {
	prometheus.MustRegister(AlertsUpdate)
	prometheus.MustRegister(AlertsApplied)
	prometheus.MustRegister(AlertsDeleted)
	prometheus.MustRegister(AlertsProcessed)
	prometheus.MustRegister(AlertsFailed)
	prometheus.MustRegister(AlertsFailedEvent)
}

// Serve health and metric requests forever.
func Serve(addr, metrics, ready, alive string) {
	h := http.NewServeMux()
	h.Handle(metrics, promhttp.Handler())
	log.Infof("HTTP server started on %s", addr)
	log.Infof("Serving metrics on %s", metrics)
	log.Infof("Serving readiness check on %s", ready)
	log.Infof("Serving liveness check on %s", alive)
	log.Info(http.ListenAndServe(addr, h))
}
