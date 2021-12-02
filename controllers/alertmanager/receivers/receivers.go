package receivers

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/nais/alerterator/utils"
	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	alertmanager "github.com/prometheus/alertmanager/config"
	"github.com/prometheus/common/config"
)

func getDefaultEmailConfig(to string) alertmanager.EmailConfig {
	return alertmanager.EmailConfig{
		To: to,
		NotifierConfig: alertmanager.NotifierConfig{
			VSendResolved: false,
		},
	}
}

// getDefaultSMSConfig returns a webhookConfig that has an endpoint that will send alerts via SMS to the recipients
// in the alert-request.
//
// HttpConfig needs to be an empty object to turn off the default httpConfig which uses proxy-settings
func getDefaultSMSConfig() alertmanager.WebhookConfig {
	url, _ := url.Parse("http://smsmanager/sms")
	return alertmanager.WebhookConfig{
		URL: &alertmanager.URL{
			URL: url,
		},
		NotifierConfig: alertmanager.NotifierConfig{
			VSendResolved: true,
		},
		HTTPConfig: &config.HTTPClientConfig{},
	}
}

func getDefaultSlackConfig(channel string) alertmanager.SlackConfig {
	if !strings.HasPrefix(channel, "#") {
		channel = "#" + channel
	}

	return alertmanager.SlackConfig{
		Channel: channel,
		NotifierConfig: alertmanager.NotifierConfig{
			VSendResolved: true,
		},
		Title:    "{{ template \"nais-alert.title\" . }}",
		Text:     "{{ template \"nais-alert.text\" . }}",
		Color:    "{{ template \"nais-alert.color\" . }}",
		Username: fmt.Sprintf("Alertmanager in %s", os.Getenv("NAIS_CLUSTER_NAME")),
	}
}

func getReceiverIndexByName(name string, receivers []*alertmanager.Receiver) int {
	for i := range receivers {
		if receivers[i].Name == name {
			return i
		}
	}
	return -1
}

func createReceiver(name string, alert *naisiov1.Alert) *alertmanager.Receiver {
	receivers := alert.Spec.Receivers
	receiver := alertmanager.Receiver{
		Name: name,
	}

	if receivers.Slack.Channel != "" {
		slack := getDefaultSlackConfig(receivers.Slack.Channel)
		if receivers.Slack.SendResolved != nil && !*receivers.Slack.SendResolved {
			slack.NotifierConfig.VSendResolved = false
		}
		if receivers.Slack.Username != "" {
			slack.Username = receivers.Slack.Username
		}
		if receivers.Slack.IconEmoji != "" {
			slack.IconEmoji = receivers.Slack.IconEmoji
		}
		if receivers.Slack.IconUrl != "" {
			slack.IconURL = receivers.Slack.IconUrl
		}
		receiver.SlackConfigs = append(receiver.SlackConfigs, &slack)
	}

	if receivers.Email.To != "" {
		email := getDefaultEmailConfig(receivers.Email.To)
		if receivers.Email.SendResolved {
			email.NotifierConfig.VSendResolved = true
		}
		receiver.EmailConfigs = append(receiver.EmailConfigs, &email)
	}

	if receivers.SMS.Recipients != "" {
		sms := getDefaultSMSConfig()
		if receivers.SMS.SendResolved != nil && !*receivers.SMS.SendResolved {
			sms.NotifierConfig.VSendResolved = false
		}
		receiver.WebhookConfigs = append(receiver.WebhookConfigs, &sms)
	}

	return &receiver
}

func AddOrUpdate(alert *naisiov1.Alert, receivers []*alertmanager.Receiver) ([]*alertmanager.Receiver, error) {
	name := utils.GetCombinedName(alert)
	receiver := createReceiver(name, alert)

	if i := getReceiverIndexByName(name, receivers); i != -1 {
		receivers[i] = receiver
	} else {
		receivers = append(receivers, receiver)
	}

	return receivers, nil
}

func Delete(alert *naisiov1.Alert, receivers []*alertmanager.Receiver) []*alertmanager.Receiver {
	name := utils.GetCombinedName(alert)
	if i := getReceiverIndexByName(name, receivers); i != -1 {
		receivers = append(receivers[:i], receivers[i+1:]...)
	}

	return receivers
}
