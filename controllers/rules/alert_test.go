package rules

import (
	"testing"

	"alerterator/controllers/fixtures"
	"alerterator/utils"
	"github.com/stretchr/testify/assert"
)

func TestAlerts(t *testing.T) {
	t.Run("Validated that alert rules are created correctly", func(t *testing.T) {
		alertRules := createAlertRules(fixtures.AlertResource)
		assert.Len(t, alertRules, 1)

		alertRule := alertRules[0]
		assert.Equal(t, utils.GetCombinedName(fixtures.AlertResource), alertRule.Labels["alert"])

		alert := fixtures.AlertResource.Spec.Alerts[0]
		assert.Equal(t, alert.For, alertRule.For)
		assert.Equal(t, alert.Expr, alertRule.Expr)
		assert.Equal(t, alert.Alert, alertRule.Alert)
		assert.Equal(t, alert.Documentation, alertRule.Annotations["documentation"])
		assert.Equal(t, alert.Description, alertRule.Annotations["description"])
		assert.Equal(t, alert.Action, alertRule.Annotations["action"])
		assert.Equal(t, alert.SLA, alertRule.Annotations["sla"])
		assert.Equal(t, fixtures.AlertResource.Spec.Receivers.Slack.PrependText, alertRule.Annotations["prependText"])
		assert.Equal(t, alert.Severity, alertRule.Annotations["severity"])
	})

	t.Run("If severity is not set, default to danger", func(t *testing.T) {
		alertRules := createAlertRules(fixtures.MinimalAlertResource)
		assert.Len(t, alertRules, 1)

		alertRule := alertRules[0]
		assert.Equal(t, "danger", alertRule.Annotations["severity"])
	})
}
