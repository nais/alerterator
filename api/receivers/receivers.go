package receivers

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/nais/alerterator/pkg/apis/alerterator/v1"
	"os"
	"strings"
)

type slackConfig struct {
	Channel      string `mapstructure:"channel" yaml:"channel"`
	SendResolved bool   `mapstructure:"send_resolved" yaml:"send_resolved"`
	Title        string `mapstructure:"title" yaml:"title"`
	Text         string `mapstructure:"text" yaml:"text"`
	Color        string `mapstructure:"color" yaml:"color,omitempty"`
	Username     string `mapstructure:"username" yaml:"username"`
}

type emailConfig struct {
	To           string `mapstructure:"to" yaml:"to"`
	SendResolved bool   `mapstructure:"send_resolved" yaml:"send_resolved"`
}

type receiverConfig struct {
	Name         string        `mapstructure:"name" yaml:"name"`
	SlackConfigs []slackConfig `mapstructure:"slack_configs" yaml:"slack_configs,omitempty"`
	EmailConfigs []emailConfig `mapstructure:"email_configs" yaml:"email,omitempty"`
}

func getDefaultEmailConfig(to string) emailConfig {
	return emailConfig{
		To:           to,
		SendResolved: false,
	}
}

func getDefaultSlackConfig(channel string) slackConfig {
	if !strings.HasPrefix(channel, "#") {
		channel = "#" + channel
	}

	return slackConfig{
		Channel:      channel,
		SendResolved: true,
		Title:        "{{ template \"nais-alert.title\" . }}",
		Text:         "{{ template \"nais-alert.text\" . }}",
		Color:        "{{ template \"nais-alert.color\" . }}",
		Username:     fmt.Sprintf("Alertmanager in %s", os.Getenv("NAIS_CLUSTER_NAME")),
	}
}

func getReceiverIndexByName(alert string, receivers []receiverConfig) int {
	for i := 0; i < len(receivers); i++ {
		receiver := receivers[i]
		if receiver.Name == alert {
			return i
		}
	}
	return -1
}

func createReceiver(alert *v1.Alert) (receiver receiverConfig) {
	receiver.Name = alert.Name

	if alert.Spec.Receivers.Slack.Channel != "" {
		slack := getDefaultSlackConfig(alert.Spec.Receivers.Slack.Channel)
		receiver.SlackConfigs = append(receiver.SlackConfigs, slack)
	}

	if alert.Spec.Receivers.Email.To != "" {
		email := getDefaultEmailConfig(alert.Spec.Receivers.Email.To)
		if alert.Spec.Receivers.Email.SendResolved {
			email.SendResolved = true
		}
		receiver.EmailConfigs = append(receiver.EmailConfigs, email)
	}
	return
}

func AddOrUpdateReceiver(alert *v1.Alert, alertManager map[interface{}]interface{}) ([]receiverConfig, error) {
	var receivers []receiverConfig
	err := mapstructure.Decode(alertManager["receivers"], &receivers)
	if err != nil {
		return nil, fmt.Errorf("failed while decoding map structure: %s", err)
	}

	receiver := createReceiver(alert)
	index := getReceiverIndexByName(alert.Name, receivers)
	if index != -1 {
		receivers[index] = receiver
	} else {
		receivers = append(receivers, receiver)
	}

	return receivers, nil
}

func DeleteReceiver(alert *v1.Alert, alertManager map[interface{}]interface{}) error {
	var receivers []receiverConfig
	err := mapstructure.Decode(alertManager["receivers"], &receivers)
	if err != nil {
		return fmt.Errorf("failed while decoding map structure: %s", err)
	}

	index := getReceiverIndexByName(alert.Name, receivers)
	if index != -1 {
		receivers = append(receivers[:index], receivers[index+1:]...)
	}
	alertManager["receivers"] = receivers

	return nil
}
