package updater

import (
	"github.com/nais/alerterator/api/fixtures"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAlerts(t *testing.T) {
	t.Run("Validerer at AlertRules blir opprettet riktig", func(t *testing.T) {
		alertRules := createAlertRules(fixtures.AlertResource)
		assert.Len(t, alertRules, 1)

		alertRule := alertRules[0]
		assert.Equal(t, fixtures.AlertResource.GetTeamName(), alertRule.Labels["team"])

		alert := fixtures.AlertResource.Spec.Alerts[0]
		assert.Equal(t, alert.For, alertRule.For)
		assert.Equal(t, alert.Expr, alertRule.Expr)
		assert.Equal(t, alert.Alert, alertRule.Alert)
		assert.Equal(t, alert.Documentation, alertRule.Annotations["documentation"])
		assert.Equal(t, alert.Description, alertRule.Annotations["description"])
		assert.Equal(t, alert.Action, alertRule.Annotations["action"])
		assert.Equal(t, alert.SLA, alertRule.Annotations["sla"])
		assert.Equal(t, fixtures.AlertResource.Spec.Receivers.Slack.PrependText, alertRule.Annotations["prependText"])
	})
}
