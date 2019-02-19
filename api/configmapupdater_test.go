package api

import (
	"testing"

	"github.com/nais/alerterator/api/updater"
	"github.com/nais/alerterator/pkg/apis/alerterator/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConfigMapUpdater(t *testing.T) {
	alert := &v1alpha1.Alert{
		ObjectMeta: metav1.ObjectMeta{
			Name: "aura",
			Labels: map[string]string{
				"team": "aura",
			},
		},
		Spec: v1alpha1.AlertSpec{
			Alerts: []v1alpha1.Rule{
				{
					Alert:  "app is down",
					For:    "2m",
					Expr:   "kube_deployment_status_replicas_unavailable{deployment=\"my-app\"} > 0",
					Action: "kubectl describe pod -l app=my-app",
				},
			},
			Receivers: v1alpha1.Receivers{
				Slack: v1alpha1.Slack{
					Channel: "#nais-alerts-dev",
				},
			},
		},
	}
	configMap := &v1.ConfigMap{
		Data: map[string]string{
			"test.yml": "hello world",
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

	expectedConfigMap := &v1.ConfigMap{
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

	t.Run("Test that alerts get added", func(t *testing.T) {
		ok := configMap.Data["test.yml"]
		assert.NotEmpty(t, ok)
		configMap := updater.DeleteAlert("test", configMap)
		ok = configMap.Data["test.yml"]
		assert.Empty(t, ok)
	})

	t.Run("Test for error if alertmanager.yml is missing", func(t *testing.T) {
		_, err := deleteReceivers(nil, &v1.ConfigMap{})
		assert.Error(t, err)
	})

	t.Run("Test that receiver and route is added correctly", func(t *testing.T) {
		configMap, err := addOrUpdateReceivers(alert, configMap)
		assert.NoError(t, err)
		assert.Equal(t, expectedConfigMap.Data["alertmanager.yml"], configMap.Data["alertmanager.yml"])
	})
}
