package api

import (
	"fmt"
	"github.com/nais/alerterator/api/updater"
	"github.com/nais/alerterator/pkg/apis/alerterator/v1alpha1"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	configMapAlerts       = "app-rules"
	configMapAlertmanager = "nais-prometheus-prometheus-alertmanager"
)

func AddOrUpdateReceivers(alert *v1alpha1.Alert, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	if configMap.Data == nil {
		return nil, fmt.Errorf("alertmanager is not properly set up, missing alertmanager.yml")
	}

	alertManager := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(configMap.Data["alertmanager.yml"]), alertManager)
	if err != nil {
		return nil, fmt.Errorf("failed while unmarshling alertmanager.yml: %s", err)
	}

	err = updater.AddOrUpdateRoutes(alert, alertManager)
	if err != nil {
		return nil, err
	}

	err = updater.AddOrUpdateReceivers(alert, alertManager)
	if err != nil {
		return nil, err
	}

	data, err := yaml.Marshal(&alertManager)
	configMap.Data["alertmanager.yml"] = string(data)

	return configMap, nil
}

func UpdateAlertManagerConfigMap(configMapInterface v1.ConfigMapInterface, alert *v1alpha1.Alert) error {
	configMap, err := configMapInterface.Get(configMapAlertmanager, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failing while retrieving %s configMap: %s", configMapAlertmanager, err)
	}

	configMap, err = AddOrUpdateReceivers(alert, configMap)
	if err != nil {
		return fmt.Errorf("failed while adding/updating receivers: %s", err)
	}

	_, err = configMapInterface.Update(configMap)
	if err != nil {
		return fmt.Errorf("failed while updating %s configMaps: %s", configMapAlertmanager, err)
	}

	return nil
}

func UpdateAppRulesConfigMap(configMapInterface v1.ConfigMapInterface, alert *v1alpha1.Alert) error {
	configMap, err := configMapInterface.Get(configMapAlerts, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failing while retrieving %s configMap: %s", configMapAlerts, err)
	}

	configMap, err = updater.AddOrUpdateAlerts(alert, configMap)
	if err != nil {
		return err
	}

	_, err = configMapInterface.Update(configMap)
	if err != nil {
		return fmt.Errorf("failed while updating %s configMaps: %s", configMapAlerts, err)
	}

	return nil
}
