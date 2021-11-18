package rules

import (
	"fmt"

	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"github.com/prometheus/common/model"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"

	"github.com/nais/alerterator/utils"
)

// RuleGroups is a set of rule groups that are typically exposed in a file.
type RuleGroups struct {
	Groups []RuleGroup `yaml:"groups"`
}

// RuleGroup is a list of sequentially evaluated recording and alerting rules.
type RuleGroup struct {
	Name     string         `yaml:"name"`
	Interval model.Duration `yaml:"interval,omitempty"`
	Limit    int            `yaml:"limit,omitempty"`
	Rules    []Rule         `yaml:"rules"`
}

// Rule describes an alerting or recording rule.
type Rule struct {
	Record      string            `yaml:"record,omitempty"`
	Alert       string            `yaml:"alert,omitempty"`
	Expr        string            `yaml:"expr"`
	For         model.Duration    `yaml:"for,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

func createAlertRules(name, slackPrependText, smsRecipients string, naisAlerts []naisiov1.Rule) ([]Rule, error) {
	var rules []Rule

	for _, ar := range naisAlerts {
		if len(ar.For) == 0 {
			ar.For = "0"
		}

		forDuration, err := model.ParseDuration(ar.For)
		if err != nil {
			return nil, err
		}
		rule := Rule{
			Alert: ar.Alert,
			Expr:  ar.Expr,
			For:   forDuration,
			Labels: map[string]string{
				"alert": name,
			},
			Annotations: map[string]string{
				"action":        ar.Action,
				"description":   ar.Description,
				"documentation": ar.Documentation,
				"prependText":   slackPrependText,
				"sla":           ar.SLA,
				"severity":      ar.Severity,
				"priority":      ar.Priority,
				"recipients":    smsRecipients,
			},
		}
		if rule.Annotations["severity"] == "" {
			rule.Annotations["severity"] = "danger"
		}
		if rule.Annotations["priority"] == "" {
			rule.Annotations["priority"] = "0"
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func AddOrUpdateAlert(alert *naisiov1.Alert, configMap corev1.ConfigMap) (corev1.ConfigMap, error) {
	name := utils.GetCombinedName(alert)
	alertRules, err := createAlertRules(name, alert.Spec.Receivers.Slack.PrependText, alert.Spec.Receivers.SMS.Recipients, alert.Spec.Alerts)
	if err != nil {
		return corev1.ConfigMap{}, err
	}
	alertGroups := RuleGroups{
		Groups: []RuleGroup{
			{
				Name:  name,
				Rules: alertRules},
		},
	}

	alertGroupYamlBytes, err := yaml.Marshal(alertGroups)
	if err != nil {
		return corev1.ConfigMap{}, fmt.Errorf("failed to marshal %v to yaml", alertGroups)
	}

	if configMap.Data == nil {
		configMap.Data = make(map[string]string)
	}

	configMap.Data[utils.GetCombinedName(alert)+".yml"] = string(alertGroupYamlBytes)

	return configMap, nil
}
