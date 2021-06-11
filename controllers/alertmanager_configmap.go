package controllers

import (
	"context"
	"fmt"
	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"alerterator/controllers/inhibitions"
	"alerterator/controllers/receivers"
	"alerterator/controllers/routes"
	"gopkg.in/yaml.v2"
)

const alertmanagerConfigName = "alertmanager.yml"

var alertmanagerConfigMapName = types.NamespacedName{
	Namespace: "nais",
	Name:      "nais-prometheus-prometheus-alertmanager",
}
var alertmanagerTemplateConfigMapName = types.NamespacedName{
	Namespace: "nais",
	Name:      "alertmanager-template-config",
}

func getConfig(ctx context.Context, namespacedName types.NamespacedName, alertReconciler *AlertReconciler) (map[interface{}]interface{}, error) {
	var configMap v1.ConfigMap
	err := alertReconciler.Get(ctx, namespacedName, &configMap)
	if err != nil {
		return nil, fmt.Errorf("failing while retrieving %s configMap: %s", namespacedName.Name, err)
	}

	if configMap.Data == nil {
		return nil, fmt.Errorf("alertmanager is not properly set up, data is empty")
	}

	config := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(configMap.Data[alertmanagerConfigName]), config)
	if err != nil {
		return nil, fmt.Errorf("failed while unmarshling %s: %s", alertmanagerConfigName, err)
	}

	return config, nil
}

func updateConfigMap(ctx context.Context, namespacedName types.NamespacedName, config map[interface{}]interface{}, alertReconciler *AlertReconciler) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed while marshaling: %s", err)
	}

	var configMap v1.ConfigMap
	err = alertReconciler.Get(ctx, namespacedName, &configMap)
	if err != nil {
		return fmt.Errorf("failing while retrieving %s configMap: %s", alertmanagerConfigMapName, err)
	}

	configMap.Data[alertmanagerConfigName] = string(data)
	err = alertReconciler.Update(ctx, &configMap)
	if err != nil {
		return fmt.Errorf("failed while updating %s: %s", alertmanagerConfigMapName, err)
	}

	return nil
}

func AddOrUpdateAlertmanagerConfigMap(ctx context.Context, alertReconciler *AlertReconciler, alert *naisiov1.Alert) error {
	currentConfig, err := getConfig(ctx, alertmanagerConfigMapName, alertReconciler)
	if err != nil {
		return err
	}
	latestConfig, err := getConfig(ctx, alertmanagerTemplateConfigMapName, alertReconciler)
	if err != nil {
		return err
	}

	updatedRoutes, err := routes.AddOrUpdateRoute(alert, currentConfig, latestConfig)
	if err != nil {
		return fmt.Errorf("failed while adding/updating routes: %s", err)
	}
	latestConfig["route"] = updatedRoutes

	updatedReceivers, err := receivers.AddOrUpdateReceiver(alert, currentConfig)
	if err != nil {
		return fmt.Errorf("failed while adding/updating receivers: %s", err)
	}
	latestConfig["receivers"] = updatedReceivers

	updatedInhibitRules, err := inhibitions.AddOrUpdateInhibition(alert, currentConfig)
	if err != nil {
		return fmt.Errorf("failed while adding/updating inhibitions: %s", err)
	}
	latestConfig["inhibit_rules"] = updatedInhibitRules

	return updateConfigMap(ctx, alertmanagerConfigMapName, latestConfig, alertReconciler)
}

func DeleteRouteAndReceiverFromAlertManagerConfigMap(ctx context.Context, alertReconciler *AlertReconciler, alert *naisiov1.Alert) error {
	config, err := getConfig(ctx, alertmanagerConfigMapName, alertReconciler)
	if err != nil {
		return err
	}

	err = routes.DeleteRoute(alert, config)
	if err != nil {
		return fmt.Errorf("failed while deleting route: %s", err)
	}

	err = receivers.DeleteReceiver(alert, config)
	if err != nil {
		return fmt.Errorf("failed while deleting receivers: %s", err)
	}

	err = inhibitions.DeleteInhibition(alert, config)
	if err != nil {
		return fmt.Errorf("failed while deleting receivers: %s", err)
	}

	return updateConfigMap(ctx, alertmanagerConfigMapName, config, alertReconciler)
}
