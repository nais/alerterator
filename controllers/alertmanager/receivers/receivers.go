package receivers

import (
	"github.com/nais/alerterator/utils"
	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	alertmanager "github.com/prometheus/alertmanager/config"
)

func getReceiverIndexByName(name string, receivers []*alertmanager.Receiver) int {
	for i := range receivers {
		if receivers[i].Name == name {
			return i
		}
	}
	return -1
}

func createReceiver(name string, alert *naisiov1.Alert) (*alertmanager.Receiver, error) {
	receivers := alert.Spec.Receivers
	receiver := alertmanager.Receiver{
		Name: name,
	}

	if receivers.Slack.Channel != "" {
		slackReceiver := createSlackReceiver(receivers.Slack)
		receiver.SlackConfigs = append(receiver.SlackConfigs, slackReceiver)
	}
	if receivers.Email.To != "" {
		emailReceiver := createEmailReceiver(receivers.Email)
		receiver.EmailConfigs = append(receiver.EmailConfigs, emailReceiver)
	}
	if receivers.SMS.Recipients != "" {
		smsReceiver := createSMSReceiver(receivers.SMS)
		receiver.WebhookConfigs = append(receiver.WebhookConfigs, smsReceiver)
	}
	if receivers.Webhook.URL != "" {
		webhookReceiver, err := createWebhookReceiver(receivers.Webhook)
		if err != nil {
			return nil, err
		}
		receiver.WebhookConfigs = append(receiver.WebhookConfigs, webhookReceiver)
	}

	return &receiver, nil
}

func AddOrUpdate(alert *naisiov1.Alert, receivers []*alertmanager.Receiver) ([]*alertmanager.Receiver, error) {
	name := utils.GetCombinedName(alert)
	receiver, err := createReceiver(name, alert)
	if err != nil {
		return nil, err
	}

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
