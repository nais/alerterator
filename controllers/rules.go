package controllers

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/nais/alerterator/controllers/rules"

	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"

	"github.com/nais/alerterator/utils"
)

var configMapAlertsNamespacedName = types.NamespacedName{
	Namespace: "nais",
	Name:      "alerterator-rules",
}

func AddOrUpdateAlert(ctx context.Context, reconciler *AlertReconciler, alert *naisiov1.Alert) error {
	var configMap v1.ConfigMap
	err := reconciler.Get(ctx, configMapAlertsNamespacedName, &configMap)
	if err != nil {
		return fmt.Errorf("failing while retrieving %s configMap: %s", configMapAlertsNamespacedName.Name, err)
	}

	configMap, err = rules.AddOrUpdateAlert(alert, configMap)
	if err != nil {
		return err
	}

	err = reconciler.Update(ctx, &configMap)
	if err != nil {
		return fmt.Errorf("failed while updating %s configMaps: %s", configMapAlertsNamespacedName.Name, err)
	}

	return nil
}

func DeleteAlert(ctx context.Context, reconciler *AlertReconciler, alert *naisiov1.Alert) error {
	var configMap v1.ConfigMap
	err := reconciler.Get(ctx, configMapAlertsNamespacedName, &configMap)
	if err != nil {
		return fmt.Errorf("failing while retrieving %s configMap: %s", configMapAlertsNamespacedName.Name, err)
	}
	delete(configMap.Data, utils.GetCombinedName(alert)+".yml")

	err = reconciler.Update(ctx, &configMap)
	if err != nil {
		return fmt.Errorf("failed while updating %s configMaps: %s", configMapAlertsNamespacedName.Name, err)
	}

	return nil
}
