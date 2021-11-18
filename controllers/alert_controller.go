package controllers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
)

const alertFinalizerName = "alert.finalizers.alerterator.nais.io"

type AlertReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=alerterator.nais.io,resources=alerts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=alerterator.nais.io,resources=alerts/status,verbs=get;update;patch

func (r *AlertReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.WithFields(log.Fields{
		"alert":         req.NamespacedName,
		"correlationId": uuid.New().String(),
	})
	logger.Info("Reconciling alert")

	var alert naisiov1.Alert
	err := r.Get(ctx, req.NamespacedName, &alert)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !alert.GetDeletionTimestamp().IsZero() {
		logger.Info("Alert resource is being deleted")
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&alert, alertFinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			logger.Debug("Deleting alert from Alertmanager")
			if err := r.deleteExternalResources(ctx, &alert); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return ctrl.Result{}, err
			}

			logger.Debug("Removing finalizer")
			controllerutil.RemoveFinalizer(&alert, alertFinalizerName)
			if err := r.Update(context.Background(), &alert); err != nil {
				return ctrl.Result{}, err
			}

			logger.Info("Finalizer processed")
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	// Add finalizer if not found
	if !controllerutil.ContainsFinalizer(&alert, alertFinalizerName) {
		logger.Debug("Finalizer not found; registering...")
		controllerutil.AddFinalizer(&alert, alertFinalizerName)
		if err := r.Update(context.Background(), &alert); err != nil {
			return ctrl.Result{}, err
		}

		logger.Info("Finalizer registered")
		return ctrl.Result{}, nil
	}

	logger.Debug("Updating Alertmanager config map")
	err = AddOrUpdateAlertmanagerConfigMap(ctx, r, &alert)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("while updating AlertManager.yml configMap: %s", err)
	}

	logger.Debug("Updating Alerterator rules config map")
	err = AddOrUpdateAlert(ctx, r, &alert)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("while adding rules to configMap: %s", err)
	}

	logger.Info("Done")
	return ctrl.Result{}, nil
}

func (r *AlertReconciler) deleteExternalResources(ctx context.Context, alert *naisiov1.Alert) error {
	err := DeleteRouteAndReceiverFromAlertManagerConfigMap(ctx, r, alert)
	if err != nil {
		return err
	}
	err = DeleteAlert(ctx, r, alert)
	if err != nil {
		return err
	}

	return nil
}

func (r *AlertReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&naisiov1.Alert{}).
		Complete(r)
}
