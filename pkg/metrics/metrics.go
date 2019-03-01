package metrics

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Alerts = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "alerts",
		Namespace: "alerterator",
		Help:      "number of alert synchronization performed",
	})
	HttpRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "http_requests",
		Namespace: "alerterator",
		Help:      "number of HTTP requests made to the health and liveness checks",
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
)

func init() {
	prometheus.MustRegister(Alerts)
	prometheus.MustRegister(HttpRequests)
	prometheus.MustRegister(AlertsApplied)
	prometheus.MustRegister(AlertsDeleted)
	prometheus.MustRegister(AlertsProcessed)
	prometheus.MustRegister(AlertsFailed)
}

func isAlive(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Alive.")
	if err != nil {
		log.Error("Failing when responding with Alive", err)
	}
	HttpRequests.Inc()
}

func isReady(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Ready.")
	if err != nil {
		log.Error("Failing when responding with Ready", err)
	}
	HttpRequests.Inc()
}

// Serve health and metric requests forever.
func Serve(addr, metrics, ready, alive string) {
	h := http.NewServeMux()
	h.Handle(metrics, promhttp.Handler())
	h.HandleFunc(ready, isReady)
	h.HandleFunc(alive, isAlive)
	log.Infof("HTTP server started on %s", addr)
	log.Infof("Serving metrics on %s", metrics)
	log.Infof("Serving readiness check on %s", ready)
	log.Infof("Serving liveness check on %s", alive)
	log.Info(http.ListenAndServe(addr, h))
}
