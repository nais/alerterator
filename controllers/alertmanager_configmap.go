package controllers

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"alerterator/controllers/inhibitions"
	"alerterator/controllers/receivers"
	"alerterator/controllers/routes"
	"github.com/nais/liberator/pkg/apis/nais.io/v1"
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

func getConfig(namespacedName types.NamespacedName, alertReconciler *AlertReconciler, ctx context.Context) (map[interface{}]interface{}, error) {
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

func updateConfigMap(namespacedName types.NamespacedName, config map[interface{}]interface{}, alertReconciler *AlertReconciler, ctx context.Context) error {
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

func AddOrUpdateAlertmanagerConfigMap(alertReconciler *AlertReconciler, ctx context.Context, alert *nais_io_v1.Alert) error {
	currentConfig, err := getConfig(alertmanagerConfigMapName, alertReconciler, ctx)
	if err != nil {
		return err
	}
	latestConfig, err := getConfig(alertmanagerTemplateConfigMapName, alertReconciler, ctx)
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

	return updateConfigMap(alertmanagerConfigMapName, latestConfig, alertReconciler, ctx)
}

func DeleteRouteAndReceiverFromAlertManagerConfigMap(alertReconciler *AlertReconciler, ctx context.Context, alert *nais_io_v1.Alert) error {
	config, err := getConfig(alertmanagerConfigMapName, alertReconciler, ctx)
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

	return updateConfigMap(alertmanagerConfigMapName, config, alertReconciler, ctx)
}
