package controllers

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	log "sigs.k8s.io/controller-runtime/pkg/log"

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
	logger := log.FromContext(ctx)
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
			logger.Info("Deleting alert from Alertmanager")
			if err := r.deleteExternalResources(ctx, &alert); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return ctrl.Result{}, err
			}

			logger.Info("Removing finalizer")
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
		logger.Info("Finalizer not found; registering...")
		controllerutil.AddFinalizer(&alert, alertFinalizerName)
		if err := r.Update(context.Background(), &alert); err != nil {
			return ctrl.Result{}, err
		}

		logger.Info("Finalizer registered")
		return ctrl.Result{}, nil
	}

	logger.Info("Updating Alertmanager config map")
	err = AddOrUpdateAlertmanagerConfigMap(ctx, r, &alert)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("while updating Alertanager config map: %s", err)
	}

	logger.Info("Updating Prometheus rules config map")
	err = AddOrUpdateRules(ctx, r, &alert)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("while adding rules to config map: %s", err)
	}

	logger.Info("Done")
	return ctrl.Result{}, nil
}

func (r *AlertReconciler) deleteExternalResources(ctx context.Context, alert *naisiov1.Alert) error {
	err := DeleteFromAlertmanagerConfigMap(ctx, r, alert)
	if err != nil {
		return err
	}
	err = DeleteRules(ctx, r, alert)
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
