package alerterator

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/nais/alerterator/pkg/apis/alerterator/v1alpha1"
	clientV1Alpha1 "github.com/nais/alerterator/pkg/client/clientset/versioned"
	informers "github.com/nais/alerterator/pkg/client/informers/externalversions/alerterator/v1alpha1"
	"github.com/nais/alerterator/pkg/metrics"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
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
		})

	return &alerterator
}

// Creates a Kubernetes event.
func (n *Alerterator) reportEvent(event *corev1.Event) (*corev1.Event, error) {
	glog.Info(event)
	return n.ClientSet.CoreV1().Events(event.Namespace).Create(event)
}

// Reports an error through the error log, a Kubernetes event, and possibly logs a failure in event creation.
func (n *Alerterator) reportError(source string, err error, alert *v1alpha1.Alert) {
	glog.Error(err)
	ev := alert.CreateEvent(source, err.Error(), "Warning", "default")
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

	// TODO: Retrieve configMap, and update alerts

	// At this point, the deployment is complete. All that is left is to register the application hash and cache it,
	// so that the deployment does not happen again. Thus, we update the metrics before the end of the function.
	// metrics.ResourcesGenerated.Add(float64(len(resources)))
	metrics.Alerts.Inc()

	alert.SetLastSyncedHash(hash)
	glog.Infof("%s: setting new hash %s", alert.Name, hash)

	alert.NilFix()
	_, err = n.AppClient.AlerteratorV1alpha1().Alerts().Update(alert)
	if err != nil {
		return fmt.Errorf("while storing alert sync metadata: %s", err)
	}

	_, err = n.reportEvent(alert.CreateEvent("synchronize", fmt.Sprintf("successfully synchronized alert resources (hash = %s)", hash), "Normal", "default"))
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
	n.update(nil, alert)
}

func (n *Alerterator) Run(stop <-chan struct{}) {
	glog.Info("Starting alert synchronization")
	if !cache.WaitForCacheSync(stop, n.AlertInformerSynced) {
		glog.Error("timed out waiting for cache sync")
		return
	}
}