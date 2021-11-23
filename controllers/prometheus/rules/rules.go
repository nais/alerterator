package rules

import (
	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"github.com/prometheus/common/model"

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

func createRules(name, slackPrependText, smsRecipients string, naisAlerts []naisiov1.Rule) ([]Rule, error) {
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
				"recipients":    smsRecipients,
			},
		}
		if rule.Annotations["severity"] == "" {
			rule.Annotations["severity"] = "danger"
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func AddOrUpdate(alert *naisiov1.Alert) (RuleGroups, error) {
	name := utils.GetCombinedName(alert)
	rules, err := createRules(name, alert.Spec.Receivers.Slack.PrependText, alert.Spec.Receivers.SMS.Recipients, alert.Spec.Alerts)
	if err != nil {
		return RuleGroups{}, err
	}
	ruleGroups := RuleGroups{
		Groups: []RuleGroup{
			{
				Name:  name,
				Rules: rules},
		},
	}

	return ruleGroups, nil
}
