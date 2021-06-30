package receivers

import (
	"fmt"
	"os"
	"strings"

	"alerterator/utils"
	"github.com/mitchellh/mapstructure"
	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"github.com/spf13/viper"
)

type slackConfig struct {
	Channel      string `mapstructure:"channel" yaml:"channel"`
	SendResolved bool   `mapstructure:"send_resolved" yaml:"send_resolved"`
	Title        string `mapstructure:"title" yaml:"title"`
	Text         string `mapstructure:"text" yaml:"text"`
	Color        string `mapstructure:"color" yaml:"color,omitempty"`
	Username     string `mapstructure:"username" yaml:"username"`
	IconUrl      string `mapstructure:"icon_url" yaml:"icon_url,omitempty"`
	IconEmoji    string `mapstructure:"icon_emoji" yaml:"icon_emoji,omitempty"`
}

type emailConfig struct {
	To           string `mapstructure:"to" yaml:"to"`
	SendResolved bool   `mapstructure:"send_resolved" yaml:"send_resolved"`
}

type webhookConfig struct {
	URL          string   `mapstructure:"url" yaml:"url"`
	SendResolved bool     `mapstructure:"send_resolved" yaml:"send_resolved"`
	HttpConfig   struct{} `mapstructure:"http_config" yaml:"http_config"`
}

type pushoverConfig struct {
	SendResolved bool   `mapstructure:"send_resolved" yaml:"send_resolved"`
	UserKey      string `mapstructure:"user_key" yaml:"user_key"`
	Token        string `mapstructure:"token" yaml:"token"`
	Title        string `mapstructure:"title" yaml:"title"`
	Message      string `mapstructure:"message" yaml:"message"`
	Priority     string `mapstructure:"priority" yaml:"priority"`
	Retry        string `mapstructure:"retry" yaml:"retry"`
	Expire       string `mapstructure:"expire" yaml:"expire"`
}

type receiverConfig struct {
	Name            string           `mapstructure:"name" yaml:"name"`
	SlackConfigs    []slackConfig    `mapstructure:"slack_configs" yaml:"slack_configs,omitempty"`
	EmailConfigs    []emailConfig    `mapstructure:"email_configs" yaml:"email_configs,omitempty"`
	WebhookConfigs  []webhookConfig  `mapstructure:"webhook_configs" yaml:"webhook_configs,omitempty"`
	PushoverConfigs []pushoverConfig `mapstructure:"pushover_configs" yaml:"pushover_configs,omitempty"`
}

func getDefaultEmailConfig(to string) emailConfig {
	return emailConfig{
		To:           to,
		SendResolved: false,
	}
}

// getDefaultSMSConfig returns a webhookConfig that has an endpoint that will send alerts via SMS to the recipients
// in the alert-request.
//
// HttpConfig needs to be an empty object to turn off the default httpConfig which uses proxy-settings
func getDefaultSMSConfig() webhookConfig {
	return webhookConfig{
		URL:          "http://smsmanager/sms",
		SendResolved: true,
		HttpConfig:   struct{}{},
	}
}

func getDefaultPushoverConfig(userKey string) pushoverConfig {
	return pushoverConfig{
		SendResolved: true,
		UserKey:      userKey,
		Token:        viper.GetString("pushover_token"),
		Title:        "{{ template \"nais-pushover.title\" . }}",
		Message:      "{{ template \"nais-pushover.text\" . }}",
		Priority:     "{{ template \"nais-pushover.priority\" }}",
		Retry:        "1m",
		Expire:       "1h",
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

func createReceiver(alert *naisiov1.Alert) (receiver receiverConfig) {
	receivers := alert.Spec.Receivers
	receiver.Name = utils.GetCombinedName(alert)

	if receivers.Slack.Channel != "" {
		slack := getDefaultSlackConfig(receivers.Slack.Channel)
		if receivers.Slack.SendResolved != nil && !*receivers.Slack.SendResolved {
			slack.SendResolved = false
		}
		if receivers.Slack.Username != "" {
			slack.Username = receivers.Slack.Username
		}
		if receivers.Slack.IconEmoji != "" {
			slack.IconEmoji = receivers.Slack.IconEmoji
		}
		if receivers.Slack.IconUrl != "" {
			slack.IconUrl = receivers.Slack.IconUrl
		}
		receiver.SlackConfigs = append(receiver.SlackConfigs, slack)
	}

	if receivers.Email.To != "" {
		email := getDefaultEmailConfig(receivers.Email.To)
		if receivers.Email.SendResolved {
			email.SendResolved = true
		}
		receiver.EmailConfigs = append(receiver.EmailConfigs, email)
	}

	if receivers.SMS.Recipients != "" {
		sms := getDefaultSMSConfig()
		if !receivers.SMS.SendResolved {
			sms.SendResolved = false
		}
		receiver.WebhookConfigs = append(receiver.WebhookConfigs, sms)
	}

	return
}

func AddOrUpdateReceiver(alert *naisiov1.Alert, alertManager map[interface{}]interface{}) ([]receiverConfig, error) {
	var receivers []receiverConfig
	err := mapstructure.Decode(alertManager["receivers"], &receivers)
	if err != nil {
		return nil, fmt.Errorf("failed while decoding map structure: %s", err)
	}

	receiver := createReceiver(alert)
	index := getReceiverIndexByName(utils.GetCombinedName(alert), receivers)
	if index != -1 {
		receivers[index] = receiver
	} else {
		receivers = append(receivers, receiver)
	}

	return receivers, nil
}

func DeleteReceiver(alert *naisiov1.Alert, alertManager map[interface{}]interface{}) error {
	var receivers []receiverConfig
	err := mapstructure.Decode(alertManager["receivers"], &receivers)
	if err != nil {
		return fmt.Errorf("failed while decoding map structure: %s", err)
	}

	index := getReceiverIndexByName(utils.GetCombinedName(alert), receivers)
	if index != -1 {
		receivers = append(receivers[:index], receivers[index+1:]...)
	}
	alertManager["receivers"] = receivers

	return nil
}
