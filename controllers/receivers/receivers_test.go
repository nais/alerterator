package receivers

import (
	"testing"

	"alerterator/utils"

	"alerterator/controllers/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestReceivers(t *testing.T) {
	t.Run("Validating that receivers are created correctly", func(t *testing.T) {
		receiver := createReceiver(fixtures.AlertResource)
		assert.Equal(t, utils.GetCombinedName(fixtures.AlertResource), receiver.Name)
		assert.Len(t, receiver.EmailConfigs, 1)
		assert.Len(t, receiver.SlackConfigs, 1)

		receivers := fixtures.AlertResource.Spec.Receivers
		assert.Equal(t, receivers.Email.To, receiver.EmailConfigs[0].To)
		assert.Equal(t, receivers.Email.SendResolved, receiver.EmailConfigs[0].SendResolved)

		assert.Equal(t, receivers.Slack.Channel, receiver.SlackConfigs[0].Channel)
		assert.Equal(t, receivers.Slack.PrependText, fixtures.AlertResource.Spec.Receivers.Slack.PrependText)
		assert.True(t, receiver.SlackConfigs[0].SendResolved)

		assert.True(t, receiver.WebhookConfigs[0].SendResolved)
	})

	t.Run("Valider at send_resolved for e-post blir beholdt", func(t *testing.T) {
		alert := fixtures.AlertResource
		alert.Spec.Receivers.Email.SendResolved = false
		receiver := createReceiver(alert)
		assert.False(t, receiver.EmailConfigs[0].SendResolved)
	})

	t.Run("Valider at send_resolved for sms blir beholdt", func(t *testing.T) {
		boolp := func(i bool) *bool {
			return &i
		}
		alert := fixtures.AlertResource
		alert.Spec.Receivers.SMS.SendResolved = boolp(false)
		receiver := createReceiver(alert)
		assert.False(t, receiver.WebhookConfigs[0].SendResolved)
	})

	t.Run("Valider at username og ikon for slack blir beholdt", func(t *testing.T) {
		boolp := func(i bool) *bool {
			return &i
		}

		alert := fixtures.AlertResource
		alert.Spec.Receivers.Slack.SendResolved = boolp(false)
		alert.Spec.Receivers.Slack.IconEmoji = ":fire:"
		alert.Spec.Receivers.Slack.IconUrl = "https://url"
		alert.Spec.Receivers.Slack.Username = "Username"
		receiver := createReceiver(alert)
		assert.Equal(t, "Username", receiver.SlackConfigs[0].Username)
		assert.Equal(t, ":fire:", receiver.SlackConfigs[0].IconEmoji)
		assert.Equal(t, "https://url", receiver.SlackConfigs[0].IconUrl)
		assert.False(t, receiver.SlackConfigs[0].SendResolved)
	})
}
