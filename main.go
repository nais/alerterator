package main

import (
	"context"
	"flag"
	"github.com/nais/alerterator/controllers/alertmanager"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/zapr"
	alertv1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/nais/alerterator/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	// TODO: Register custom metrics with the global prometheus registry
	_ = clientgoscheme.AddToScheme(scheme)

	_ = alertv1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	zapLogger, err := cfg.Build()

	if err != nil {
		setupLog.Error(err, "Unable to set up controller logger")
		os.Exit(1)
	}

	ctrl.SetLogger(zapr.NewLogger(zapLogger))

	kconfig, err := ctrl.GetConfig()
	simpleClient, err := client.New(kconfig, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		setupLog.Error(err, "Unable to create go client")
		os.Exit(1)
	}
	err = alertmanager.EnsureConfigExists(context.Background(), simpleClient, setupLog)
	if err != nil {
		setupLog.Error(err, "Unable to set up config")
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "ade053be.nais.io",
	})
	if err != nil {
		setupLog.Error(err, "Unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.AlertReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "Unable to create controller")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "Problem running manager: ")
		os.Exit(1)
	}
}
