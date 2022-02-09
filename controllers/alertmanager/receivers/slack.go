package receivers

import (
	"fmt"
	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	alertmanager "github.com/prometheus/alertmanager/config"
	"os"
	"strings"
)

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

func createSlackReceiver(slack nais_io_v1.Slack) *alertmanager.SlackConfig {
	slackConfig := getDefaultSlackConfig(slack.Channel)
	if slack.SendResolved != nil && !*slack.SendResolved {
		slackConfig.NotifierConfig.VSendResolved = false
	}
	if slack.Username != "" {
		slackConfig.Username = slack.Username
	}
	if slack.IconEmoji != "" {
		slackConfig.IconEmoji = slack.IconEmoji
	}
	if slack.IconUrl != "" {
		slackConfig.IconURL = slack.IconUrl
	}

	return &slackConfig
}
