package rules

import (
	"testing"

	"github.com/nais/alerterator/controllers/fixtures"
	"github.com/nais/alerterator/utils"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
)

func TestRules(t *testing.T) {
	t.Run("Validated that alert rules are created correctly", func(t *testing.T) {
		naisAlert := fixtures.AlertResource()
		name := utils.GetCombinedName(naisAlert)
		alertRules, err := createRules(name, naisAlert.Spec.Receivers.Slack.PrependText, naisAlert.Spec.Receivers.SMS.Recipients, naisAlert.Spec.Alerts)
		assert.NoError(t, err)
		assert.Len(t, alertRules, 1)

		alertRule := alertRules[0]
		assert.Equal(t, utils.GetCombinedName(naisAlert), alertRule.Labels["alert"])

		alert := naisAlert.Spec.Alerts[0]

		forDuration, err := model.ParseDuration(alert.For)
		assert.NoError(t, err)

		assert.Equal(t, forDuration, alertRule.For)
		assert.Equal(t, alert.Expr, alertRule.Expr)
		assert.Equal(t, alert.Alert, alertRule.Alert)
		assert.Equal(t, alert.Documentation, alertRule.Annotations["documentation"])
		assert.Equal(t, alert.Description, alertRule.Annotations["description"])
		assert.Equal(t, alert.Action, alertRule.Annotations["action"])
		assert.Equal(t, alert.SLA, alertRule.Annotations["sla"])
		assert.Equal(t, naisAlert.Spec.Receivers.Slack.PrependText, alertRule.Annotations["prependText"])
		assert.Equal(t, alert.Severity, alertRule.Annotations["severity"])
	})

	t.Run("If severity is not set, default to danger", func(t *testing.T) {
		alert := fixtures.MinimalAlertResource()
		name := utils.GetCombinedName(alert)
		alertRules, err := createRules(name, alert.Spec.Receivers.Slack.PrependText, alert.Spec.Receivers.SMS.Recipients, alert.Spec.Alerts)
		assert.NoError(t, err)
		assert.Len(t, alertRules, 1)

		alertRule := alertRules[0]
		assert.Equal(t, "danger", alertRule.Annotations["severity"])
	})

	t.Run("Ensure that for==0 if not set", func(t *testing.T) {
		naisAlert := fixtures.MinimalAlertResource()
		naisAlert.Spec.Alerts[0].For = ""
		name := utils.GetCombinedName(naisAlert)
		rules, err := createRules(name, naisAlert.Spec.Receivers.Slack.PrependText, naisAlert.Spec.Receivers.SMS.Recipients, naisAlert.Spec.Alerts)
		assert.NoError(t, err)
		assert.Len(t, rules, 1)

		rule := rules[0]
		assert.Equal(t, model.Duration(0), rule.For)
	})
}

func TestConfigMapUpdater(t *testing.T) {
	t.Run("Test that alerts get added", func(t *testing.T) {
		alert := fixtures.AlertResource()
		configMap, err := AddOrUpdate(alert, *fixtures.ConfigMapBeforeAlerts())
		assert.NoError(t, err)
		assert.Equal(t, fixtures.ExpectedConfigMapAfterAlerts().Data["aura.yml"], configMap.Data["aura.yml"])
	})
}
