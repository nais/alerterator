package receivers

import (
	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	alertmanager "github.com/prometheus/alertmanager/config"
	"github.com/prometheus/common/config"
	"net/url"
)

// getDefaultWebhook returns a empty webhookConfig. HttpConfig needs to be an empty object to turn off the default
// httpConfig which uses proxy-settings
func getDefaultWebhook() alertmanager.WebhookConfig {
	return alertmanager.WebhookConfig{
		NotifierConfig: alertmanager.NotifierConfig{
			VSendResolved: true,
		},
		HTTPConfig: &config.HTTPClientConfig{},
	}
}

// getDefaultSMSConfig returns a webhookConfig that has an endpoint that will send alerts via SMS to the recipients
// in the alert-request.
func getDefaultSMSConfig() alertmanager.WebhookConfig {
	webhook := getDefaultWebhook()
	url, _ := url.Parse("http://smsmanager/sms")
	webhook.URL = &alertmanager.URL{URL: url}
	return webhook
}

func createSMSReceiver(sms nais_io_v1.SMS) *alertmanager.WebhookConfig {
	smsConfig := getDefaultSMSConfig()
	if sms.SendResolved != nil && !*sms.SendResolved {
		smsConfig.NotifierConfig.VSendResolved = false
	}

	return &smsConfig
}

func createWebhookReceiver(webhook nais_io_v1.Webhook) (*alertmanager.WebhookConfig, error) {
	webhookConfig := getDefaultWebhook()
	if webhook.SendResolved != nil && !*webhook.SendResolved {
		webhookConfig.NotifierConfig.VSendResolved = false
	}
	webhookConfig.MaxAlerts = uint64(webhook.MaxAlerts)
	url, err := url.Parse(webhook.URL)
	if err != nil {
		return nil, err
	}
	webhookConfig.URL = &alertmanager.URL{URL: url}
	if webhook.HttpConfig.ProxyUrl != "" {
		proxyUrl, err := url.Parse(webhook.HttpConfig.ProxyUrl)
		if err != nil {
			return nil, err
		}
		webhookConfig.HTTPConfig.ProxyURL = config.URL{URL: proxyUrl}
	}
	if webhook.HttpConfig.TLSConfig.InsecureSkipVerify {
		webhookConfig.HTTPConfig.TLSConfig.InsecureSkipVerify = true
	}

	return &webhookConfig, nil
}
