package updater

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/nais/alerterator/pkg/apis/alerterator/v1alpha1"
)

type slackConfig struct {
	Channel      string `mapstructure:"channel" yaml:"channel"`
	SendResolved bool   `mapstructure:"send_resolved" yaml:"send_resolve"`
	Title        string `mapstructure:"title" yaml:"title"`
	Text         string `mapstructure:"text" yaml:"text"`
	Username     string `mapstructure:"username" yaml:"username"`
}

type emailConfig struct {
	To          string `mapstructure:"to" yaml:"to"`
	SendResolve bool   `mapstructure:"send_resolve" yaml:"send_resolve"`
}

func getDefaultEmailConfig() emailConfig {
	return emailConfig{
		SendResolve: true,
	}
}

type receiverConfig struct {
	Name         string        `mapstructure:"name" yaml:"name"`
	SlackConfigs []slackConfig `mapstructure:"slack_configs" yaml:"slack_configs,omitempty"`
	EmailConfigs []emailConfig `mapstructure:"email" yaml:"email,omitempty"`
}

func getDefaultSlackConfig() slackConfig {
	return slackConfig{
		SendResolved: true,
		Title:        "{{ template \"nais-alert.title\" . }}", // TODO: Add prepend-text
		Text:         "{{ template \"nais-alert.text\" . }}",  // TODO: Add prepend-text
		Username:     fmt.Sprintf("Alertmanager in %s", os.Getenv("NAIS_CLUSTER_NAME")),
	}
}

func getReceiverByName(alert string, receivers []receiverConfig) (receiverConfig, int) {
	for i := 0; i < len(receivers); i++ {
		receiver := receivers[i]
		if receiver.Name == alert {
			return receiver, i
		}
	}
	return receiverConfig{}, -1
}

func AddOrUpdateReceivers(alert *v1alpha1.Alert, alertManager map[interface{}]interface{}) error {
	var receivers []receiverConfig
	err := mapstructure.Decode(alertManager["receivers"], &receivers)
	if err != nil {
		return fmt.Errorf("failed while decoding map structure: %s", err)
	}

	receiver, index := getReceiverByName(alert.Name, receivers)
	if receiver.Name != "" {
		if alert.Spec.Receivers.Slack.Channel != "" {
			slack := getDefaultSlackConfig()
			slack.Channel = alert.Spec.Receivers.Slack.Channel
			receiver.SlackConfigs = []slackConfig{slack}
		}
		if alert.Spec.Receivers.Email.To != "" {
			email := getDefaultEmailConfig()
			email.To = alert.Spec.Receivers.Email.To
			receiver.EmailConfigs = []emailConfig{email}
		}
		receivers[index] = receiver
	} else {
		receiver := receiverConfig{
			Name: alert.Name,
		}
		if alert.Spec.Receivers.Slack.Channel != "" {
			slack := getDefaultSlackConfig()
			slack.Channel = alert.Spec.Receivers.Slack.Channel
			receiver.SlackConfigs = append(receiver.SlackConfigs, slack)
		}
		if alert.Spec.Receivers.Email.To != "" {
			email := getDefaultEmailConfig()
			email.To = alert.Spec.Receivers.Email.To
			receiver.EmailConfigs = append(receiver.EmailConfigs, email)
		}
		receivers = append(receivers, receiver)
	}

	alertManager["receivers"] = receivers

	return nil
}
