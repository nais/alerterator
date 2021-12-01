package alertmanager

import (
	"context"
	"fmt"

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

	routes, err := routes.AddOrUpdate(alert, oldConfig.Route.Routes)
	if err != nil {
		return fmt.Errorf("failed while adding/updating routes: %s", err)
	}
	mergeRoutes(newConfig.Route.Routes, routes)

	receivers, err := receivers.AddOrUpdate(alert, oldConfig.Receivers)
	if err != nil {
		return fmt.Errorf("failed while adding/updating receivers: %s", err)
	}
	mergeReceivers(newConfig.Receivers, receivers)

	inhibitRules, err := inhibitions.AddOrUpdate(alert, oldConfig.InhibitRules)
	if err != nil {
		return fmt.Errorf("failed while adding/updating inhibitions: %s", err)
	}
	newConfig.InhibitRules = inhibitRules

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

	routes := routes.Delete(alert, oldConfig.Route.Routes)
	newConfig.Route.Routes = routes

	receivers := receivers.Delete(alert, oldConfig.Receivers)
	newConfig.Receivers = receivers

	inhibitions := inhibitions.Delete(alert, oldConfig.InhibitRules)
	newConfig.InhibitRules = inhibitions

	return configmap.MarshalAndUpdateData(ctx, client, alertmanagerConfigMapName, alertmanagerConfigName, newConfig)
}
