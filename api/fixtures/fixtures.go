package fixtures

import (
	"github.com/nais/alerterator/pkg/apis/alerterator/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var AlertResource = &v1.Alert{
	ObjectMeta: metav1.ObjectMeta{
		Name: "aura",
		Namespace: "aura",
		Labels: map[string]string{
			"alert": "aura",
		},
	},
	Spec: v1.AlertSpec{
		Receivers: v1.Receivers{
			Slack: v1.Slack{
				Channel:     "#nais-alerts-dev",
				PrependText: "<!here>",
			},
			Email: v1.Email{
				To: "test@example.com",
			},
		},
		Alerts: []v1.Rule{
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
	},
}

var MinimalAlertResource = &v1.Alert{
	ObjectMeta: metav1.ObjectMeta{
		Name: "aura",
		Namespace: "aura",
		Labels: map[string]string{
			"alert": "aura",
		},
	},
	Spec: v1.AlertSpec{
		Receivers: v1.Receivers{
			Slack: v1.Slack{
				Channel: "#nais-alerts-dev",
			},
		},
		Alerts: []v1.Rule{
			{
				Alert:  "app is down",
				For:    "2m",
				Expr:   "kube_deployment_status_replicas_unavailable{deployment=\"my-app\"} > 0",
				Action: "kubectl describe pod -l app=my-app",
			},
		},
	},
}

var AlertmanagerConfigYaml = `
global:
  slack_api_url: web-site.com
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
  - name: aura
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
    - receiver: aura
      continue: true
      match:
        alert: aura
    - receiver: testmann
      continue: true
      match:
        alert: testmann`

var AlertmanagerConfigYamlDifferentRoutes = `
route:
  group_by: ['alertname','team', 'kubernetes_namespace']
  group_wait: 100s
  group_interval: 50m
  repeat_interval: 10h
  receiver: default-receiver
  routes: []`

var ConfigMapBeforeAlerts = &corev1.ConfigMap{
	Data: map[string]string{},
}

var ExpectedConfigMapAfterAlerts = &corev1.ConfigMap{
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
      alert: aura
`,
	},
}
