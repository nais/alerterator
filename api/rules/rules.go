package rules

import (
	"fmt"
	"github.com/nais/alerterator/pkg/apis/alerterator/v1"
	"github.com/nais/alerterator/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const ConfigMapAlerts = "alerterator-rules"

func AddOrUpdateAlert(configMapInterface corev1.ConfigMapInterface, alert *v1.Alert) error {
	configMap, err := configMapInterface.Get(ConfigMapAlerts, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failing while retrieving %s configMap: %s", ConfigMapAlerts, err)
	}

	configMap, err = addOrUpdateAlert(alert, configMap)
	if err != nil {
		return err
	}

	_, err = configMapInterface.Update(configMap)
	if err != nil {
		return fmt.Errorf("failed while updating %s configMaps: %s", ConfigMapAlerts, err)
	}

	return nil
}

func DeleteAlert(configMapInterface corev1.ConfigMapInterface, alert *v1.Alert) error {
	configMap, err := configMapInterface.Get(ConfigMapAlerts, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failing while retrieving %s configMap: %s", ConfigMapAlerts, err)
	}
	delete(configMap.Data, utils.GetCombinedName(alert)+".yml")

	_, err = configMapInterface.Update(configMap)
	if err != nil {
		return fmt.Errorf("failed while updating %s configMaps: %s", ConfigMapAlerts, err)
	}

	return nil
}
