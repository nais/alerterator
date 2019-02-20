package fixtures

import (
	"github.com/nais/alerterator/pkg/apis/alerterator/v1alpha1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var AlertResource = &v1alpha1.Alert{
	ObjectMeta: metav1.ObjectMeta{
		Name: "aura",
		Labels: map[string]string{
			"team": "aura",
		},
	},
	Spec: v1alpha1.AlertSpec{
		Receivers: v1alpha1.Receivers{
			Slack: v1alpha1.Slack{
				Channel:     "#nais-alerts-dev",
				PrependText: "<!here>",
			},
			Email: v1alpha1.Email{
				To: "test@example.com",
			},
		},
		Alerts: []v1alpha1.Rule{
			{
				Alert:         "app is down",
				For:           "2m",
				Expr:          "kube_deployment_status_replicas_unavailable{deployment=\"my-app\"} > 0",
				Documentation: "some documentation, or link to documentation",
				Action:        "kubectl describe pod -l app=my-app",
				Description:   "this is a description of the alert",
				SLA:           "we need to fix this ASAP",
			},
		},
	},
}

var ConfigMapBeforeDelete = &v1.ConfigMap{
	Data: map[string]string{
		"alertmanager.yml": `
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
        team: aura
    - receiver: testmann
      continue: true
      match:
        team: testmann`,
	},
}

var ExpectedConfigMapAfterDelete = &v1.ConfigMap{
	Data: map[string]string{
		"alertmanager.yml": `global:
  http_config:
    proxy_url: http://webproxy.nais:8088
  slack_api_url: web-site.com
  smtp_auth_password: blorg
  smtp_auth_username: blarg
  smtp_from: srvKubernetesAlarm@nav.no
  smtp_require_tls: false
  smtp_smarthost: smtp.preprod.local:26
receivers:
- name: default-receiver
  slack_configs:
  - channel: '#nais-alerts-default'
    send_resolved: true
    title: '{{ template "nais-alert.title" . }}'
    text: '{{ template "nais-alert.text" . }}'
    username: Alertmanager in preprod-fss
route:
  group_by:
  - alertname
  - team
  - kubernetes_namespace
  group_wait: 10s
  group_interval: 5m
  repeat_interval: 1h
  receiver: default-receiver
  routes:
  - receiver: testmann
    continue: true
    match:
      team: testmann
templates:
- /etc/config/alert.tmpl
`,
	},
}

var ConfigMapBeforeAdd = &v1.ConfigMap{
	Data: map[string]string{
		"alertmanager.yml": `
global:
  http_config:
    proxy_url: http://webproxy.nais:8088
  slack_api_url: web-site.com
  smtp_auth_password: blorg
  smtp_auth_username: blarg
  smtp_from: srvKubernetesAlarm@nav.no
  smtp_require_tls: false
  smtp_smarthost: smtp.preprod.local:26
receivers:
- name: default-receiver
  slack_configs:
  - channel: '#nais-alerts-default'
    send_resolved: true
    title: '{{ template "nais-alert.title" . }}'
    text: '{{ template "nais-alert.text" . }}'
    username: Alertmanager in preprod-fss
route:
  group_by:
  - alertname
  - team
  - kubernetes_namespace
  group_wait: 10s
  group_interval: 5m
  repeat_interval: 1h
  receiver: default-receiver
  routes:
  - receiver: testmann
    continue: true
    match:
      team: testmann
templates:
- /etc/config/alert.tmpl`,
	},
}

var ExpectedConfigMapAfterReceivers = &v1.ConfigMap{
	Data: map[string]string{
		"alertmanager.yml": `global:
  http_config:
    proxy_url: http://webproxy.nais:8088
  slack_api_url: web-site.com
  smtp_auth_password: blorg
  smtp_auth_username: blarg
  smtp_from: srvKubernetesAlarm@nav.no
  smtp_require_tls: false
  smtp_smarthost: smtp.preprod.local:26
receivers:
- name: default-receiver
  slack_configs:
  - channel: '#nais-alerts-default'
    send_resolved: true
    title: '{{ template "nais-alert.title" . }}'
    text: '{{ template "nais-alert.text" . }}'
    username: Alertmanager in preprod-fss
- name: aura
  slack_configs:
  - channel: '#nais-alerts-dev'
    send_resolved: true
    title: '{{ template "nais-alert.title" . }}'
    text: '{{ template "nais-alert.text" . }}'
    username: 'Alertmanager in '
  email:
  - to: test@example.com
    send_resolved: false
route:
  group_by:
  - alertname
  - team
  - kubernetes_namespace
  group_wait: 10s
  group_interval: 5m
  repeat_interval: 1h
  receiver: default-receiver
  routes:
  - receiver: testmann
    continue: true
    match:
      team: testmann
  - receiver: aura
    continue: true
    match:
      team: aura
templates:
- /etc/config/alert.tmpl
`,
	},
}

var ConfigMapBeforeAlerts = &v1.ConfigMap{
	Data: map[string]string{},
}

var ExpectedConfigMapAfterAlerts = &v1.ConfigMap{
	Data: map[string]string{
		"aura.yml": `groups:
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
      sla: we need to fix this ASAP
    labels:
      team: aura
`,
	},
}
