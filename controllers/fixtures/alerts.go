package fixtures

import (
	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func boolp(i bool) *bool {
	return &i
}

func AlertResource() *naisiov1.Alert {
	return &naisiov1.Alert{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aura",
			Namespace: "aura",
			Labels: map[string]string{
				"team": "aura",
			},
		},
		Spec: naisiov1.AlertSpec{
			Route: naisiov1.Route{
				RepeatInterval: "4h",
				GroupWait:      "30s",
				GroupInterval:  "5m",
				GroupBy:        []naisiov1.LabelName{"my_alert", "my_team", "my_kubernetes_namespace"},
			},
			Receivers: naisiov1.Receivers{
				Slack: naisiov1.Slack{
					Channel:      "#nais-alerts-dev",
					PrependText:  "<!here>",
					IconEmoji:    ":fire:",
					IconUrl:      "https://url.emoji",
					Username:     "Username",
					SendResolved: boolp(false),
				},
				Email: naisiov1.Email{
					To:           "test@example.com",
					SendResolved: true,
				},
				SMS: naisiov1.SMS{
					Recipients:   "12346789",
					SendResolved: boolp(false),
				},
				Webhook: naisiov1.Webhook{
					URL:          "http://historymanager.nais",
					MaxAlerts:    0,
					SendResolved: boolp(true),
					HttpConfig: naisiov1.HttpConfig{
						ProxyUrl: "http://no-proxy.nav",
						TLSConfig: naisiov1.TLSConfig{
							InsecureSkipVerify: false,
						},
					},
				},
			},
			Alerts: []naisiov1.Rule{
				{
					Alert:         "app is down",
					For:           "2m",
					Expr:          "kube_deployment_status_replicas_unavailable{deployment=\"my-app\"} > 0",
					Documentation: "some documentation, or link to documentation",
					Action:        "kubectl describe pod -l app=my-app",
					Description:   "this is a description of the alert",
					SLA:           "we need to fix this ASAP",
					Severity:      "#eeeeee",
				},
			},
			InhibitRules: []naisiov1.InhibitRules{
				{
					Targets: map[string]string{
						"alert": "kube_deployment_status_replicas_unavailable",
					},
					Sources: map[string]string{
						"alert": "naisCluster",
					},
				},
			},
		},
	}
}

func MinimalAlertResource() *naisiov1.Alert {
	return &naisiov1.Alert{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aura",
			Namespace: "aura",
			Labels: map[string]string{
				"alert": "aura",
			},
		},
		Spec: naisiov1.AlertSpec{
			Receivers: naisiov1.Receivers{
				Slack: naisiov1.Slack{
					Channel: "#nais-alerts-dev",
				},
			},
			Alerts: []naisiov1.Rule{
				{
					Alert:  "app is down",
					For:    "2m",
					Expr:   "kube_deployment_status_replicas_unavailable{deployment=\"my-app\"} > 0",
					Action: "kubectl describe pod -l app=my-app",
				},
			},
		},
	}
}
