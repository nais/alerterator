package prometheus

import (
	"context"

	"github.com/nais/alerterator/controllers/configmap"
	"github.com/nais/alerterator/controllers/prometheus/rules"
	"github.com/nais/alerterator/utils"
	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var prometheusRulesConfigMap = types.NamespacedName{
	Namespace: "nais",
	Name:      "alerterator-rules",
}

func AddOrUpdateRules(ctx context.Context, client client.Client, alert *naisiov1.Alert) error {
	ruleGroups, err := rules.AddOrUpdate(alert)
	if err != nil {
		return err
	}

	return configmap.MarshalAndUpdateData(ctx, client, prometheusRulesConfigMap, utils.GetCombinedName(alert)+".yml", ruleGroups)
}

func DeleteRules(ctx context.Context, client client.Client, alert *naisiov1.Alert) error {
	return configmap.DeleteFileFromData(ctx, client, prometheusRulesConfigMap, utils.GetCombinedName(alert)+".yml")
}
