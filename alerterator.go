package alerterator

import (
	"fmt"
	"github.com/nais/alerterator/api"

	"github.com/golang/glog"
	"github.com/nais/alerterator/pkg/apis/alerterator/v1alpha1"
	clientV1Alpha1 "github.com/nais/alerterator/pkg/client/clientset/versioned"
	informers "github.com/nais/alerterator/pkg/client/informers/externalversions/alerterator/v1alpha1"
	"github.com/nais/alerterator/pkg/metrics"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

const (
	configMapNamespace = "nais"
)

// Alerterator is a singleton that holds Kubernetes client instances.
type Alerterator struct {
	ClientSet           kubernetes.Interface
	AppClient           *clientV1Alpha1.Clientset
	AlertInformer       informers.AlertInformer
	AlertInformerSynced cache.InformerSynced
}

func NewAlerterator(clientSet kubernetes.Interface, appClient *clientV1Alpha1.Clientset, alertInformer informers.AlertInformer) *Alerterator {
	alerterator := Alerterator{
		ClientSet:           clientSet,
		AppClient:           appClient,
		AlertInformer:       alertInformer,
		AlertInformerSynced: alertInformer.Informer().HasSynced}

	alertInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(newPod interface{}) {
				alerterator.add(newPod)
			},
			UpdateFunc: func(oldPod, newPod interface{}) {
				alerterator.update(oldPod, newPod)
			},
			DeleteFunc: func(delPod interface{}) {
				alerterator.delete(delPod)
			},
		})

	return &alerterator
}

// Creates a Kubernetes event.
func (n *Alerterator) reportEvent(event *corev1.Event) (*corev1.Event, error) {
	return n.ClientSet.CoreV1().Events(event.Namespace).Create(event)
}

// Reports an error through the error log, a Kubernetes event, and possibly logs a failure in event creation.
func (n *Alerterator) reportError(source string, err error, alert *v1alpha1.Alert) {
	glog.Error(err)
	ev := alert.CreateEvent(source, err.Error(), "Warning")
	_, err = n.reportEvent(ev)
	if err != nil {
		glog.Errorf("While creating an event for this error, another error occurred: %s", err)
	}
}

func (n *Alerterator) synchronize(previous, alert *v1alpha1.Alert) error {
	hash, err := alert.Hash()
	if err != nil {
		return fmt.Errorf("while hashing alert spec: %s", err)
	}
	if alert.LastSyncedHash() == hash {
		glog.Infof("%s: no changes", alert.Name)
		return nil
	}
	// Kubernetes events needs a namespace when created, and it needs to be the same as the alerts.
	// Alerts are cluster-wide, so we just add the 'default'-namespace as an easy fix
	alert.Namespace = "default"

	err = api.UpdateAlertManagerConfigMap(n.ClientSet.CoreV1().ConfigMaps(configMapNamespace), alert)
	if err != nil {
		return fmt.Errorf("while updating AlertManager.yml configMap: %s", err)
	}

	err = api.UpdateAppRulesConfigMap(n.ClientSet.CoreV1().ConfigMaps(configMapNamespace), alert)
	if err != nil {
		return fmt.Errorf("while adding rules to configMap: %s", err)
	}

	metrics.Alerts.Inc()

	alert.SetLastSyncedHash(hash)
	glog.Infof("%s: setting new hash %s", alert.Name, hash)

	alert.NilFix()
	_, err = n.AppClient.AlerteratorV1alpha1().Alerts().Update(alert)
	if err != nil {
		return fmt.Errorf("while storing alert sync metadata: %s", err)
	}

	_, err = n.reportEvent(alert.CreateEvent("synchronize", fmt.Sprintf("successfully synchronized alert resources (hash = %s)", hash), "Normal"))
	if err != nil {
		glog.Errorf("While creating an event for this error, another error occurred: %s", err)
	}

	return nil
}

func (n *Alerterator) update(old, new interface{}) {
	var alert, previous *v1alpha1.Alert
	if old != nil {
		previous = old.(*v1alpha1.Alert)
	}
	if new != nil {
		alert = new.(*v1alpha1.Alert)
	}

	metrics.AlertsProcessed.Inc()
	glog.Infof("%s: synchronizing alert", alert.Name)

	if err := n.synchronize(previous, alert); err != nil {
		metrics.AlertsFailed.Inc()
		glog.Errorf("%s: error %s", alert.Name, err)
		n.reportError("synchronize", err, alert)
	} else {
		glog.Infof("%s: synchronized successfully", alert.Name)
	}

	glog.Infof("%s: finished synchronizing", alert.Name)
}

func (n *Alerterator) add(alert interface{}) {
	glog.Info("Applying new alert")
	metrics.AlertsApplied.Inc()
	n.update(nil, alert)
}

func (n *Alerterator) delete(alert interface{}) {
	glog.Infof("%s: deleted", alert.(*v1alpha1.Alert).Name)
	metrics.AlertsDeleted.Inc()
}

func (n *Alerterator) Run(stop <-chan struct{}) {
	glog.Info("Starting alert synchronization")
	if !cache.WaitForCacheSync(stop, n.AlertInformerSynced) {
		glog.Error("timed out waiting for cache sync")
		return
	}
}
