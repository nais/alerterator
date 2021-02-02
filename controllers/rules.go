package controllers

import (
	"alerterator/controllers/rules"
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"alerterator/utils"
	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
)

var configMapAlertsNamespacedName = types.NamespacedName{
	Namespace: "nais",
	Name:      "alerterator-rules",
}

func AddOrUpdateAlert(reconciler *AlertReconciler, ctx context.Context, alert *naisiov1.Alert) error {
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

func DeleteAlert(reconciler *AlertReconciler, ctx context.Context, alert *naisiov1.Alert) error {
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
