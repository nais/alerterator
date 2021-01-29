package controllers

import (
	"alerterator/utils"
	"context"
	"fmt"
	logrus "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
)

const alertFinalizerName = "alert.finalizers.alerterator.nais.io"

// AlertReconciler reconciles a Alert object
type AlertReconciler struct {
	client.Client
	Log    logrus.Entry
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=alerterator.nais.io,resources=alerts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=alerterator.nais.io,resources=alerts/status,verbs=get;update;patch

func (r *AlertReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithField("alert", req.NamespacedName)

	var alert *naisiov1.Alert
	err := r.Get(ctx, req.NamespacedName, alert)
	if err != nil {
		log.Error(err, "Unable to fetch Alert")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// examine DeletionTimestamp to determine if object is under deletion
	if alert.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !utils.ContainsString(alert.ObjectMeta.Finalizers, alertFinalizerName) {
			alert.ObjectMeta.Finalizers = append(alert.ObjectMeta.Finalizers, alertFinalizerName)
			if err := r.Update(context.Background(), alert); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if utils.ContainsString(alert.ObjectMeta.Finalizers, alertFinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if err := r.deleteExternalResources(alert); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			alert.ObjectMeta.Finalizers = utils.RemoveString(alert.ObjectMeta.Finalizers, alertFinalizerName)
			if err := r.Update(context.Background(), alert); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	// your logic here
	err = AddOrUpdateAlertmanagerConfigMap(r, ctx, alert)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("while updating AlertManager.yml configMap: %s", err)
	}

	err = AddOrUpdateAlert(alert, r, ctx)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("while adding rules to configMap: %s", err)
	}

	return ctrl.Result{}, nil
}

func (r *AlertReconciler) deleteExternalResources(alert *naisiov1.Alert) error {
	ctx := context.Background()
	err := DeleteRouteAndReceiverFromAlertManagerConfigMap(r, ctx, alert)
	if err != nil {
		return err
	}
	err = DeleteAlert(alert, r, ctx)
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
