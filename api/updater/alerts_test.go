package updater

import (
	"testing"

	"github.com/nais/alerterator/pkg/apis/alerterator/v1alpha1"
	"github.com/stretchr/testify/assert"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAlerts(t *testing.T) {
	alertResource := v1alpha1.Alert{
		ObjectMeta: metav1.ObjectMeta{
			Name: "TestAlerts",
			Labels: map[string]string{
				"team": "test",
			},
		},
		Spec: v1alpha1.AlertSpec{
			Receivers: v1alpha1.Receivers{
				Slack: v1alpha1.Slack{
					Channel:     "#example",
					PrependText: "<!here>",
				},
				Email: v1alpha1.Email{
					To: "test@example.com",
				},
			},
			Alerts: []v1alpha1.Rule{
				{
					Alert:         "alertName",
					For:           "2m",
					Expr:          "some Prometheus expression",
					Documentation: "some documentation, or link to documentation",
					Action:        "what to do when triggered?",
					Description:   "this is a description of the alert",
					SLA:           "we need to fix this ASAP",
				},
			},
		},
	}

	t.Run("Validerer at AlertRules blir opprettet riktig", func(t *testing.T) {
		alertRules := createAlertRules(&alertResource)
		assert.Len(t, alertRules, 1)

		alertRule := alertRules[0]
		assert.Equal(t, alertResource.GetTeamName(), alertRule.Labels["team"])

		alert := alertResource.Spec.Alerts[0]
		assert.Equal(t, alert.For, alertRule.For)
		assert.Equal(t, alert.Expr, alertRule.Expr)
		assert.Equal(t, alert.Alert, alertRule.Alert)
		assert.Equal(t, alert.Documentation, alertRule.Annotations["documentation"])
		assert.Equal(t, alert.Description, alertRule.Annotations["description"])
		assert.Equal(t, alert.Action, alertRule.Annotations["action"])
		assert.Equal(t, alert.SLA, alertRule.Annotations["sla"])
		assert.Equal(t, alertResource.Spec.Receivers.Slack.PrependText, alertRule.Annotations["prependText"])
	})
}
