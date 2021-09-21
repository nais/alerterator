package rules

import (
	"testing"

	"github.com/nais/alerterator/controllers/fixtures"
	"github.com/nais/alerterator/utils"
	"github.com/stretchr/testify/assert"
)

func TestAlerts(t *testing.T) {
	t.Run("Validated that alert rules are created correctly", func(t *testing.T) {
		naisAlert := fixtures.AlertResource()
		alertRules := CreateAlertRules(naisAlert)
		assert.Len(t, alertRules, 1)

		alertRule := alertRules[0]
		assert.Equal(t, utils.GetCombinedName(naisAlert), alertRule.Labels["alert"])

		alert := naisAlert.Spec.Alerts[0]
		assert.Equal(t, alert.For, alertRule.For)
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
		alertRules := CreateAlertRules(alert)
		assert.Len(t, alertRules, 1)

		alertRule := alertRules[0]
		assert.Equal(t, "danger", alertRule.Annotations["severity"])
	})
}
