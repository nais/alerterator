package updater

import (
	"fmt"
	"github.com/nais/alerterator/pkg/apis/alerterator/v1alpha1"
	"gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
)

type AlertGroups struct {
	Groups []AlertGroup
}

type AlertGroup struct {
	Name  string
	Rules []AlertRule
}

type AlertRule struct {
	For         string            `json:"for"`
	Expr        string            `json:"expr"`
	Annotations map[string]string `json:"annotations"`
	Labels      map[string]string `json:"labels"`
}

func createAlertRules(alert *v1alpha1.Alert) (alertRules []AlertRule) {
	for i := range alert.Spec.Alerts {
		rule := alert.Spec.Alerts[i]
		alertRule := AlertRule{
			Expr: rule.Expr,
			For:  rule.For,
			Labels: map[string]string{
				"team": alert.GetTeamName(),
			},
			Annotations: map[string]string{
				"action":        rule.Action,
				"description":   rule.Description,
				"documentation": rule.Documentation,
				"prependText":   alert.Spec.Receivers.Slack.PrependText,
				"sla":           rule.SLA,
			},
		}
		alertRules = append(alertRules, alertRule)
	}
	return
}

func AddOrUpdateAlerts(alert *v1alpha1.Alert, configMap *v1.ConfigMap) (*v1.ConfigMap, error) {
	alertRules := createAlertRules(alert)
	alertGroup := AlertGroup{Name: alert.Name, Rules: alertRules}
	alertGroups := AlertGroups{Groups: []AlertGroup{alertGroup}}

	alertGroupYamlBytes, err := yaml.Marshal(alertGroups)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %v to yaml\n", alertGroup)
	}

	if configMap.Data == nil {
		configMap.Data = make(map[string]string)
	}

	configMap.Data[alert.Name+".yml"] = string(alertGroupYamlBytes)

	return configMap, nil
}

func DeleteAlert(alertName string, configMap *v1.ConfigMap) *v1.ConfigMap {
	delete(configMap.Data, alertName+".yml")
	return configMap
}
