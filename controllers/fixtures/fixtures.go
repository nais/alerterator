package fixtures

import (
	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
				GroupBy:        []string{"group-by"},
			},
			Receivers: naisiov1.Receivers{
				Slack: naisiov1.Slack{
					Channel:     "#nais-alerts-dev",
					PrependText: "<!here>",
				},
				Email: naisiov1.Email{
					To: "test@example.com",
				},
				SMS: naisiov1.SMS{
					Recipients: "12346789",
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

const AlertmanagerConfigYaml = `
global:
  slack_api_url: https://web-site.com
  http_config:
    proxy_url: http://webproxy.nais:8088
  smtp_from: srvKubernetesAlarm@nav.no
  smtp_smarthost: smtp.preprod.local:26
  smtp_auth_username: blarg
  smtp_auth_password: blorg
  smtp_require_tls: false
templates:
- '/etc/config/alert.tmpl'
receivers:
  - name: default-receiver
    slack_configs:
    - channel: '#nais-alerts-default'
      send_resolved: true
      title: '{{ template "nais-alert.title" . }}'
      text: '{{ template "nais-alert.text" . }}'
      username: 'Alertmanager in preprod-fss'
  - name: aura-aura
    slack_configs:
    - channel: '#nais-alerts-dev'
      username: 'Alertmanager in preprod-fss'
      send_resolved: true
      title: '{{ template "nais-alert.title" . }}'
      text: '{{ template "nais-alert.text" . }}'
route:
  group_by: ['alertname','team', 'kubernetes_namespace']
  group_wait: 10s
  group_interval: 5m
  repeat_interval: 1h
  receiver: default-receiver
  routes:
    - receiver: aura-aura
      continue: true
      group_by: ["group_by"]
      match:
        alert: aura-aura
inhibit_rules:
  - target_match:
       alertname: 'applikasjon nede'
    source_match:
       alertname: 'nais_down'
    equal: ['team']
  - target_match:
       alertname: 'http_500'
    source_match:
       alertname: 'brreg_down'
    equal: ['team']`

const AlertmanagerConfigYamlDifferentRoutes = `
global:
  slack_api_url: https://web-site.com
receivers:
  - name: aura-aura
    slack_configs:
    - channel: '#nais-alerts-default'
      send_resolved: true
      title: '{{ template "nais-alert.title" . }}'
      text: '{{ template "nais-alert.text" . }}'
      username: 'Alertmanager in preprod-fss'
route:
  group_by: ['alertname','team', 'kubernetes_namespace']
  group_wait: 100s
  group_interval: 50m
  repeat_interval: 10h
  receiver: aura-aura
  routes: []
`

func ConfigMapBeforeAlerts() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		Data: map[string]string{},
	}
}

func ExpectedConfigMapAfterAlerts() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		Data: map[string]string{
			"aura-aura.yml": `groups:
- name: aura
  rules:
  - alert: app is down
    for: 2m
    expr: kube_deployment_status_replicas_unavailable{deployment="my-app"} > 0
    annotations:
      action: kubectl describe pod -l app=my-app
      description: this is a description of the alert
      documentation: some documentation, or link to documentation
      prependText: <!here>
      severity: '#eeeeee'
      sla: we need to fix this ASAP
    labels:
      alert: aura-aura
`,
		},
	}
}

const EmptyRouteConfig = `
route: {}
`
