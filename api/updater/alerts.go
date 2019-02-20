package updater

import (
	"fmt"
	"github.com/nais/alerterator/pkg/apis/alerterator/v1alpha1"
	"gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
)

type AlertGroups struct {
	Groups []AlertGroup `yaml:"groups"`
}

type AlertGroup struct {
	Name  string      `yaml:"name"`
	Rules []AlertRule `yaml:"rules"`
}

type AlertRule struct {
	Alert       string            `yaml:"alert"`
	For         string            `yaml:"for"`
	Expr        string            `yaml:"expr"`
	Annotations map[string]string `yaml:"annotations"`
	Labels      map[string]string `yaml:"labels"`
}

func createAlertRules(alert *v1alpha1.Alert) (alertRules []AlertRule) {
	for i := range alert.Spec.Alerts {
		rule := alert.Spec.Alerts[i]
		alertRule := AlertRule{
			Alert: rule.Alert,
			Expr:  rule.Expr,
			For:   rule.For,
			Labels: map[string]string{
				"team": alert.GetTeamName(),
			},
			Annotations: map[string]string{
				"action":        rule.Action,
				"description":   rule.Description,
				"documentation": rule.Documentation,
				"prependText":   alert.Spec.Receivers.Slack.PrependText,
				"sla":           rule.SLA,
				"severity":      rule.Severity,
			},
		}
		if alertRule.Annotations["severity"] == "" {
			alertRule.Annotations["severity"] = "danger"
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
