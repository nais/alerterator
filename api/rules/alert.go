package rules

import (
	"fmt"

	"github.com/nais/alerterator/pkg/apis/alerterator/v1alpha1"
	"gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
)

type Groups struct {
	Groups []Group `yaml:"groups"`
}

type Group struct {
	Name  string  `yaml:"name"`
	Rules []Alert `yaml:"rules"`
}

type Alert struct {
	Alert       string            `yaml:"alert"`
	For         string            `yaml:"for"`
	Expr        string            `yaml:"expr"`
	Annotations map[string]string `yaml:"annotations"`
	Labels      map[string]string `yaml:"labels"`
}

func createAlertRules(alert *v1alpha1.Alert) (alertRules []Alert) {
	for i := range alert.Spec.Alerts {
		rule := alert.Spec.Alerts[i]
		alertRule := Alert{
			Alert: rule.Alert,
			Expr:  rule.Expr,
			For:   rule.For,
			Labels: map[string]string{
				"alert": alert.Name,
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

func addOrUpdateAlert(alert *v1alpha1.Alert, configMap *v1.ConfigMap) (*v1.ConfigMap, error) {
	alertRules := createAlertRules(alert)
	alertGroups := Groups{
		Groups: []Group{
			{
				Name:  alert.Name,
				Rules: alertRules},
		},
	}

	alertGroupYamlBytes, err := yaml.Marshal(alertGroups)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %v to yaml\n", alertGroups)
	}

	if configMap.Data == nil {
		configMap.Data = make(map[string]string)
	}

	configMap.Data[alert.Name+".yml"] = string(alertGroupYamlBytes)

	return configMap, nil
}
