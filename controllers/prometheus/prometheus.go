package prometheus

import (
	"context"
	"fmt"

	"github.com/nais/alerterator/controllers/prometheus/rules"
	"github.com/nais/alerterator/utils"
	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var configMapAlertsNamespacedName = types.NamespacedName{
	Namespace: "nais",
	Name:      "alerterator-rules",
}

func AddOrUpdateRules(ctx context.Context, client client.Client, alert *naisiov1.Alert) error {
	var configMap v1.ConfigMap
	err := client.Get(ctx, configMapAlertsNamespacedName, &configMap)
	if err != nil {
		return fmt.Errorf("failing while retrieving %s configMap: %s", configMapAlertsNamespacedName.Name, err)
	}

	configMap, err = rules.AddOrUpdate(alert, configMap)
	if err != nil {
		return err
	}

	err = client.Update(ctx, &configMap)
	if err != nil {
		return fmt.Errorf("failed while updating %s configMaps: %s", configMapAlertsNamespacedName.Name, err)
	}

	return nil
}

func DeleteRules(ctx context.Context, client client.Client, alert *naisiov1.Alert) error {
	var configMap v1.ConfigMap
	err := client.Get(ctx, configMapAlertsNamespacedName, &configMap)
	if err != nil {
		return fmt.Errorf("failing while retrieving %s configMap: %s", configMapAlertsNamespacedName.Name, err)
	}
	delete(configMap.Data, utils.GetCombinedName(alert)+".yml")

	err = client.Update(ctx, &configMap)
	if err != nil {
		return fmt.Errorf("failed while updating %s configMaps: %s", configMapAlertsNamespacedName.Name, err)
	}

	return nil
}
