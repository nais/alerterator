package fixtures

const AlertmanagerBaseConfigYaml = `
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
route:
  group_by: ['alertname', 'team', 'kubernetes_namespace']
  group_wait: 10s
  group_interval: 5m
  repeat_interval: 1h
  receiver: default-receiver
  routes:
    - receiver: default-receiver
      continue: true
      match:
        alert: default
`

const AlertmanagerOldConfigYaml = `
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
  - name: aura-aura
    slack_configs:
    - channel: '#nais-alerts-dev'
      username: 'Alertmanager in preprod-fss'
      send_resolved: true
      title: '{{ template "nais-alert.title" . }}'
      text: '{{ template "nais-alert.text" . }}'
route:
  group_by: ['alertname', 'team', 'kubernetes_namespace']
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
    - receiver: aura-aura
      continue: true
      group_by: ["group_by"]
      match:
        alert: aura-aura
`
const AlertmanagerChangedConfigYaml = `global:
  http_config:
    proxy_url: http://webproxy.nais:8088
    follow_redirects: true
  smtp_from: srvKubernetesAlarm@nav.no
  smtp_smarthost: smtp.preprod.local:26
  smtp_auth_username: blarg
  smtp_auth_password: blorg
  smtp_require_tls: false
  slack_api_url: https://web-site.com
route:
  receiver: default-receiver
  group_by:
  - alertname
  - team
  - kubernetes_namespace
  continue: false
  routes:
  - receiver: default-receiver
    match:
      alert: default
    continue: true
  - receiver: aura-aura
    group_by:
    - my_alert
    - my_team
    - my_kubernetes_namespace
    match:
      alert: aura-aura
    continue: true
    group_wait: 30s
    group_interval: 5m
    repeat_interval: 4h
  group_wait: 10s
  group_interval: 5m
  repeat_interval: 1h
inhibit_rules:
- source_matchers:
  - alert="naisCluster"
  target_matchers:
  - alert="kube_deployment_status_replicas_unavailable"
  equal:
  - team
receivers:
- name: default-receiver
  slack_configs:
  - send_resolved: true
    channel: '#nais-alerts-default'
    username: Alertmanager in preprod-fss
    color: '{{ if eq .Status "firing" }}danger{{ else }}good{{ end }}'
    title: '{{ template "nais-alert.title" . }}'
    title_link: '{{ template "slack.default.titlelink" . }}'
    pretext: '{{ template "slack.default.pretext" . }}'
    text: '{{ template "nais-alert.text" . }}'
    short_fields: false
    footer: '{{ template "slack.default.footer" . }}'
    fallback: '{{ template "slack.default.fallback" . }}'
    callback_id: '{{ template "slack.default.callbackid" . }}'
    icon_emoji: '{{ template "slack.default.iconemoji" . }}'
    icon_url: '{{ template "slack.default.iconurl" . }}'
    link_names: false
- name: aura-aura
  email_configs:
  - send_resolved: true
    to: test@example.com
  slack_configs:
  - send_resolved: false
    channel: '#nais-alerts-dev'
    username: Username
    color: '{{ template "nais-alert.color" . }}'
    title: '{{ template "nais-alert.title" . }}'
    text: '{{ template "nais-alert.text" . }}'
    short_fields: false
    icon_emoji: ':fire:'
    icon_url: https://url.emoji
    link_names: false
  webhook_configs:
  - send_resolved: false
    http_config:
      follow_redirects: false
    url: http://smsmanager/sms
    max_alerts: 0
  - send_resolved: true
    http_config:
      proxy_url: http://no-proxy.nav
      follow_redirects: false
    url: http://historymanager.nais
    max_alerts: 0
templates:
- /etc/config/alert.tmpl
`
