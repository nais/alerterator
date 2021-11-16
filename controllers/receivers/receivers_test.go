package receivers

import (
	"testing"

	"github.com/nais/alerterator/utils"

	"github.com/nais/alerterator/controllers/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestReceivers(t *testing.T) {
	t.Run("Validating that receivers are created correctly", func(t *testing.T) {
		alert := fixtures.AlertResource()
		name := utils.GetCombinedName(alert)
		receiver := createReceiver(name, alert)
		assert.Equal(t, name, receiver.Name)
		assert.Len(t, receiver.EmailConfigs, 1)
		assert.Len(t, receiver.SlackConfigs, 1)

		alertReceivers := alert.Spec.Receivers
		assert.Equal(t, alertReceivers.Email.To, receiver.EmailConfigs[0].To)
		assert.Equal(t, alertReceivers.Email.SendResolved, receiver.EmailConfigs[0].SendResolved())

		assert.Equal(t, alertReceivers.Slack.Channel, receiver.SlackConfigs[0].Channel)
		assert.Equal(t, alertReceivers.Slack.PrependText, alert.Spec.Receivers.Slack.PrependText)
		assert.True(t, receiver.SlackConfigs[0].SendResolved())

		assert.True(t, receiver.WebhookConfigs[0].SendResolved())
	})

	t.Run("Valider at send_resolved for e-post blir beholdt", func(t *testing.T) {
		alert := fixtures.AlertResource()
		name := utils.GetCombinedName(alert)
		alert.Spec.Receivers.Email.SendResolved = true
		receiver := createReceiver(name, alert)
		assert.True(t, receiver.EmailConfigs[0].SendResolved())
	})

	t.Run("Valider at send_resolved for sms blir beholdt", func(t *testing.T) {
		boolp := func(i bool) *bool {
			return &i
		}
		alert := fixtures.AlertResource()
		name := utils.GetCombinedName(alert)
		alert.Spec.Receivers.SMS.SendResolved = boolp(false)
		receiver := createReceiver(name, alert)
		assert.False(t, receiver.WebhookConfigs[0].SendResolved())
	})

	t.Run("Valider at username og ikon for slack blir beholdt", func(t *testing.T) {
		boolp := func(i bool) *bool {
			return &i
		}

		alert := fixtures.AlertResource()
		name := utils.GetCombinedName(alert)
		alert.Spec.Receivers.Slack.SendResolved = boolp(false)
		alert.Spec.Receivers.Slack.IconEmoji = ":fire:"
		alert.Spec.Receivers.Slack.IconUrl = "https://url"
		alert.Spec.Receivers.Slack.Username = "Username"
		receiver := createReceiver(name, alert)
		assert.Equal(t, "Username", receiver.SlackConfigs[0].Username)
		assert.Equal(t, ":fire:", receiver.SlackConfigs[0].IconEmoji)
		assert.Equal(t, "https://url", receiver.SlackConfigs[0].IconURL)
		assert.False(t, receiver.SlackConfigs[0].SendResolved())
	})
}
