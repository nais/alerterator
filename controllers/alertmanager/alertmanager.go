package alertmanager

import (
	"context"
	"fmt"

	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nais/alerterator/controllers/alertmanager/inhibitions"
	"github.com/nais/alerterator/controllers/alertmanager/receivers"
	"github.com/nais/alerterator/controllers/alertmanager/routes"
	"github.com/nais/alerterator/controllers/configmap"
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

func AddOrUpdate(ctx context.Context, client client.Client, alert *naisiov1.Alert) error {
	oldConfig, err := configmap.Get(ctx, alertmanagerConfigMapName, client, alertmanagerConfigName)
	if err != nil {
		return err
	}
	newConfig, err := configmap.Get(ctx, alertmanagerTemplateConfigMapName, client, alertmanagerConfigName)
	if err != nil {
		return err
	}

	routes, err := routes.AddOrUpdate(alert, oldConfig.Route.Routes)
	if err != nil {
		return fmt.Errorf("failed while adding/updating routes: %s", err)
	}
	newConfig.Route.Routes = routes

	receivers, err := receivers.AddOrUpdate(alert, oldConfig.Receivers)
	if err != nil {
		return fmt.Errorf("failed while adding/updating receivers: %s", err)
	}
	newConfig.Receivers = receivers

	inhibitRules, err := inhibitions.AddOrUpdate(alert, oldConfig.InhibitRules)
	if err != nil {
		return fmt.Errorf("failed while adding/updating inhibitions: %s", err)
	}
	newConfig.InhibitRules = inhibitRules

	return configmap.Update(ctx, alertmanagerConfigMapName, newConfig, client, alertmanagerConfigName)
}

func Delete(ctx context.Context, client client.Client, alert *naisiov1.Alert) error {
	oldConfig, err := configmap.Get(ctx, alertmanagerConfigMapName, client, alertmanagerConfigName)
	if err != nil {
		return err
	}
	newConfig, err := configmap.Get(ctx, alertmanagerTemplateConfigMapName, client, alertmanagerConfigName)
	if err != nil {
		return err
	}

	routes := routes.Delete(alert, oldConfig.Route.Routes)
	newConfig.Route.Routes = routes

	receivers := receivers.Delete(alert, oldConfig.Receivers)
	newConfig.Receivers = receivers

	inhibitions := inhibitions.Delete(alert, oldConfig.InhibitRules)
	newConfig.InhibitRules = inhibitions

	return configmap.Update(ctx, alertmanagerConfigMapName, newConfig, client, alertmanagerConfigName)
}
