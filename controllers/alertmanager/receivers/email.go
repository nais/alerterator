package receivers

import (
	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	alertmanager "github.com/prometheus/alertmanager/config"
)

func getDefaultEmailConfig(to string) alertmanager.EmailConfig {
	return alertmanager.EmailConfig{
		To: to,
		NotifierConfig: alertmanager.NotifierConfig{
			VSendResolved: false,
		},
	}
}

func createEmailReceiver(email nais_io_v1.Email) *alertmanager.EmailConfig {
	emailConfig := getDefaultEmailConfig(email.To)
	if email.SendResolved {
		emailConfig.NotifierConfig.VSendResolved = true
	}

	return &emailConfig
}
