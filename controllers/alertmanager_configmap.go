package controllers

import (
	"context"
	"fmt"

	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/nais/alerterator/controllers/inhibitions"
	"github.com/nais/alerterator/controllers/receivers"
	"github.com/nais/alerterator/controllers/routes"
	alertmanager "github.com/prometheus/alertmanager/config"
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

func getConfig(ctx context.Context, namespacedName types.NamespacedName, alertReconciler *AlertReconciler) (*alertmanager.Config, error) {
	var configMap v1.ConfigMap
	err := alertReconciler.Get(ctx, namespacedName, &configMap)
	if err != nil {
		return nil, fmt.Errorf("failing while retrieving %s configMap: %s", namespacedName.Name, err)
	}

	if configMap.Data == nil {
		return nil, fmt.Errorf("alertmanager is not properly set up, data is empty")
	}

	config := alertmanager.Config{}
	err = yaml.Unmarshal([]byte(configMap.Data[alertmanagerConfigName]), config)
	if err != nil {
		return nil, fmt.Errorf("failed while unmarshling %s: %s", alertmanagerConfigName, err)
	}

	return &config, nil
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
	oldConfig, err := getConfig(ctx, alertmanagerConfigMapName, alertReconciler)
	if err != nil {
		return err
	}
	newConfig, err := getConfig(ctx, alertmanagerTemplateConfigMapName, alertReconciler)
	if err != nil {
		return err
	}

	routes, err := routes.AddOrUpdateRoute(alert, oldConfig.Route.Routes)
	if err != nil {
		return fmt.Errorf("failed while adding/updating routes: %s", err)
	}
	newConfig.Route.Routes = routes

	updatedReceivers, err := receivers.AddOrUpdateReceiver(alert, oldConfig)
	if err != nil {
		return fmt.Errorf("failed while adding/updating receivers: %s", err)
	}
	newConfig["receivers"] = updatedReceivers

	updatedInhibitRules, err := inhibitions.AddOrUpdateInhibition(alert, oldConfig)
	if err != nil {
		return fmt.Errorf("failed while adding/updating inhibitions: %s", err)
	}
	newConfig["inhibit_rules"] = updatedInhibitRules

	return updateConfigMap(ctx, alertmanagerConfigMapName, newConfig, alertReconciler)
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
