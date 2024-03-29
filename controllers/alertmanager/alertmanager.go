package alertmanager

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"

	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	alertmanager "github.com/prometheus/alertmanager/config"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nais/alerterator/controllers/alertmanager/inhibitions"
	"github.com/nais/alerterator/controllers/alertmanager/receivers"
	"github.com/nais/alerterator/controllers/alertmanager/routes"
	"github.com/nais/alerterator/controllers/configmap"
	"github.com/nais/alerterator/controllers/overrides"
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

func mergeRoutes(a, b []*alertmanager.Route) []*alertmanager.Route {
	m := make(map[string]bool)

	for _, r := range a {
		m[r.Receiver] = true
	}

	for _, r := range b {
		if !m[r.Receiver] {
			a = append(a, r)
			m[r.Receiver] = true
		}
	}
	return a
}

func mergeReceivers(a, b []*alertmanager.Receiver) []*alertmanager.Receiver {
	m := make(map[string]bool)

	for _, r := range a {
		m[r.Name] = true
	}

	for _, r := range b {
		if !m[r.Name] {
			a = append(a, r)
			m[r.Name] = true
		}
	}
	return a
}

func addOrUpdate(alert *naisiov1.Alert, oldConfig, newConfig *overrides.Config) (*overrides.Config, error) {
	routes, err := routes.AddOrUpdate(alert, oldConfig.Route.Routes)
	if err != nil {
		return nil, fmt.Errorf("failed while adding/updating routes: %s", err)
	}
	newConfig.Route.Routes = mergeRoutes(newConfig.Route.Routes, routes)

	receivers, err := receivers.AddOrUpdate(alert, oldConfig.Receivers)
	if err != nil {
		return nil, fmt.Errorf("failed while adding/updating receivers: %s", err)
	}
	newConfig.Receivers = mergeReceivers(newConfig.Receivers, receivers)

	inhibitRules, err := inhibitions.AddOrUpdate(alert, oldConfig.InhibitRules)
	if err != nil {
		return nil, fmt.Errorf("failed while adding/updating inhibitions: %s", err)
	}
	newConfig.InhibitRules = inhibitRules

	return newConfig, nil
}

func AddOrUpdate(ctx context.Context, client client.Client, alert *naisiov1.Alert) error {
	var oldConfig *overrides.Config
	err := configmap.GetAndUnmarshal(ctx, client, alertmanagerConfigMapName, alertmanagerConfigName, &oldConfig)
	if err != nil {
		return err
	}
	var newConfig *overrides.Config
	err = configmap.GetAndUnmarshal(ctx, client, alertmanagerTemplateConfigMapName, alertmanagerConfigName, &newConfig)
	if err != nil {
		return err
	}

	newConfig, err = addOrUpdate(alert, oldConfig, newConfig)
	if err != nil {
		return err
	}

	return configmap.MarshalAndUpdateData(ctx, client, alertmanagerConfigMapName, alertmanagerConfigName, newConfig)
}

func Delete(ctx context.Context, client client.Client, alert *naisiov1.Alert) error {
	var oldConfig overrides.Config
	err := configmap.GetAndUnmarshal(ctx, client, alertmanagerConfigMapName, alertmanagerConfigName, &oldConfig)
	if err != nil {
		return err
	}
	var newConfig overrides.Config
	err = configmap.GetAndUnmarshal(ctx, client, alertmanagerTemplateConfigMapName, alertmanagerConfigName, &newConfig)
	if err != nil {
		return err
	}

	newConfig.Route.Routes = routes.Delete(alert, oldConfig.Route.Routes)
	newConfig.Receivers = receivers.Delete(alert, oldConfig.Receivers)
	newConfig.InhibitRules = inhibitions.Delete(alert, oldConfig.InhibitRules)

	return configmap.MarshalAndUpdateData(ctx, client, alertmanagerConfigMapName, alertmanagerConfigName, newConfig)
}

func EnsureConfigExists(ctx context.Context, client client.Client, logger logr.Logger) error {
	configMap, err := configmap.Get(ctx, client, alertmanagerConfigMapName)
	if err != nil {
		return err
	}
	exists := configMap.Data[alertmanagerConfigName]
	if exists == "" {
		logger.Info(fmt.Sprintf("Configmap %v was missing %v, creating it based on %v",
			alertmanagerConfigMapName.Name, alertmanagerConfigName, alertmanagerTemplateConfigMapName.Name))
		var newConfig *overrides.Config
		err = configmap.GetAndUnmarshal(ctx, client, alertmanagerTemplateConfigMapName, alertmanagerConfigName, &newConfig)
		if err != nil {
			return err
		}

		return configmap.MarshalAndUpdateData(ctx, client, alertmanagerConfigMapName, alertmanagerConfigName, newConfig)
	}

	return nil
}
