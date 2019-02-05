package metrics

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Alerts = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "alerts",
		Namespace: "alerterator",
		Help:      "number of alert deployments performed",
	})
	HttpRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "http_requests",
		Namespace: "alerterator",
		Help:      "number of HTTP requests made to the health and liveness checks",
	})
	AlertsDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "alerts_deleted",
		Namespace: "alerterator",
		Help:      "number of nais.io.Alert resources that have been deleted",
	})
	AlertsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "applications_processed",
		Namespace: "alerterator",
		Help:      "number of nais.io.Alert resources that have been processed",
	})
	AlertsFailed = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "alerts_failed",
		Namespace: "alerterator",
		Help:      "number of nais.io.Alert resources that failed processing",
	})
	ResourcesGenerated = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "resources_generated",
		Namespace: "alerterator",
		Help:      "number of alerts-files that have been parsed/generated as a result of alert deployment",
	})
)

func init() {
	prometheus.MustRegister(Alerts)
	prometheus.MustRegister(HttpRequests)
	prometheus.MustRegister(AlertsDeleted)
	prometheus.MustRegister(AlertsProcessed)
	prometheus.MustRegister(AlertsFailed)
}

func isAlive(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Alive.")
	if err != nil {
		glog.Error("Failing when responding with Alive", err)
	}
	HttpRequests.Inc()
}

func isReady(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Ready.")
	if err != nil {
		glog.Error("Failing when responding with Ready", err)
	}
	HttpRequests.Inc()
}

// Serve health and metric requests forever.
func Serve(addr, metrics, ready, alive string) {
	h := http.NewServeMux()
	h.Handle(metrics, promhttp.Handler())
	h.HandleFunc(ready, isReady)
	h.HandleFunc(alive, isAlive)
	glog.Infof("HTTP server started on %s", addr)
	glog.Infof("Serving metrics on %s", metrics)
	glog.Infof("Serving readiness check on %s", ready)
	glog.Infof("Serving liveness check on %s", alive)
	glog.Info(http.ListenAndServe(addr, h))
}
