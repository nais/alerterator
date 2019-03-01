package alerterator

import (
	"fmt"

	"github.com/nais/alerterator/api"
	"github.com/nais/alerterator/pkg/apis/alerterator/v1alpha1"
	clientV1Alpha1 "github.com/nais/alerterator/pkg/client/clientset/versioned"
	informers "github.com/nais/alerterator/pkg/client/informers/externalversions/alerterator/v1alpha1"
	"github.com/nais/alerterator/pkg/metrics"
	log "github.com/sirupsen/logrus"
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
	log.Error(err)
	ev := alert.CreateEvent(source, err.Error(), "Warning")
	_, err = n.reportEvent(ev)
	if err != nil {
		log.Errorf("While creating an event for this error, another error occurred: %s", err)
	}
}

func (n *Alerterator) synchronize(previous, alert *v1alpha1.Alert) error {
	hash, err := alert.Hash()
	if err != nil {
		return fmt.Errorf("while hashing alert spec: %s", err)
	}
	if alert.LastSyncedHash() == hash {
		log.Infof("%s: no changes", alert.Name)
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

	alert.SetLastSyncedHash(hash)
	log.Infof("%s: setting new hash %s", alert.Name, hash)
	metrics.AlertsProcessed.Inc()

	alert.NilFix()
	_, err = n.AppClient.AlerteratorV1alpha1().Alerts().Update(alert)
	if err != nil {
		return fmt.Errorf("while storing alert sync metadata: %s", err)
	}

	_, err = n.reportEvent(alert.CreateEvent("synchronize", fmt.Sprintf("successfully synchronized alert resources (hash = %s)", hash), "Normal"))
	if err != nil {
		log.Errorf("While creating an event for this error, another error occurred: %s", err)
		metrics.AlertsFailedEvent.Inc()
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

	log.Infof("%s: synchronizing alert", alert.Name)

	if err := n.synchronize(previous, alert); err != nil {
		metrics.AlertsFailed.Inc()
		log.Errorf("%s: error %s", alert.Name, err)
		n.reportError("synchronize", err, alert)
	} else {
		log.Infof("%s: synchronized successfully", alert.Name)
		metrics.AlertsUpdate.Inc()
	}

	log.Infof("%s: finished synchronizing", alert.Name)
}

func (n *Alerterator) add(newAlert interface{}) {
	log.Info("Applying new alert")
	alert := newAlert.(*v1alpha1.Alert)

	log.Infof("%s: adding alert", alert.Name)

	if err := n.synchronize(nil, alert); err != nil {
		metrics.AlertsFailed.Inc()
		log.Errorf("%s: error %s", alert.Name, err)
		n.reportError("adding", err, alert)
	} else {
		log.Infof("%s: adding successfully", alert.Name)
		metrics.AlertsApplied.Inc()
	}

	log.Infof("%s: finished adding", alert.Name)
}

func (n *Alerterator) delete(delete interface{}) {
	alert := delete.(*v1alpha1.Alert)

	// Kubernetes events needs a namespace when created, and it needs to be the same as the alerts.
	// Alerts are cluster-wide, so we just add the 'default'-namespace as an easy fix
	alert.Namespace = "default"

	err := api.DeleteReceiversFromAlertManagerConfigMap(n.ClientSet.CoreV1().ConfigMaps(configMapNamespace), alert)
	if err != nil {
		metrics.AlertsFailed.Inc()
		log.Errorf("while deleting %s from AlertManager.yml configMap: %s", alert.Name, err)
		return
	}

	err = api.DeleteAlertFromConfigMap(n.ClientSet.CoreV1().ConfigMaps(configMapNamespace), alert)
	if err != nil {
		metrics.AlertsFailed.Inc()
		log.Errorf("while deleting rules for %s from the configMap: %s", alert.Name, err)
		return
	}

	log.Infof("%s: deleted", alert.Name)
	metrics.AlertsDeleted.Inc()

	_, err = n.reportEvent(alert.CreateEvent("synchronize", fmt.Sprintf("successfully deleted alert resources (name = %s)", alert.Name), "Normal"))
	if err != nil {
		log.Errorf("While creating an event for this error, another error occurred: %s", err)
		metrics.AlertsFailedEvent.Inc()
	}
}

func (n *Alerterator) Run(stop <-chan struct{}) {
	log.Info("Starting alert synchronization")
	if !cache.WaitForCacheSync(stop, n.AlertInformerSynced) {
		log.Error("timed out waiting for cache sync")
		return
	}
}
