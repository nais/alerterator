package overrides

import (
	alertmanager "github.com/prometheus/alertmanager/config"
	commoncfg "github.com/prometheus/common/config"
)

type GlobalConfig struct {
	HTTPConfig *commoncfg.HTTPClientConfig `yaml:"http_config,omitempty" json:"http_config,omitempty"`

	SMTPFrom         string                `yaml:"smtp_from,omitempty" json:"smtp_from,omitempty"`
	SMTPSmarthost    alertmanager.HostPort `yaml:"smtp_smarthost,omitempty" json:"smtp_smarthost,omitempty"`
	SMTPAuthUsername string                `yaml:"smtp_auth_username,omitempty" json:"smtp_auth_username,omitempty"`
	SMTPAuthPassword string                `yaml:"smtp_auth_password,omitempty" json:"smtp_auth_password,omitempty"`
	SMTPRequireTLS   bool                  `yaml:"smtp_require_tls" json:"smtp_require_tls,omitempty"`
	SlackAPIURL      string                `yaml:"slack_api_url,omitempty" json:"slack_api_url,omitempty"`
}

type Config struct {
	Global            *GlobalConfig                   `yaml:"global,omitempty" json:"global,omitempty"`
	Route             *alertmanager.Route             `yaml:"route,omitempty" json:"route,omitempty"`
	InhibitRules      []*alertmanager.InhibitRule     `yaml:"inhibit_rules,omitempty" json:"inhibit_rules,omitempty"`
	Receivers         []*alertmanager.Receiver        `yaml:"receivers,omitempty" json:"receivers,omitempty"`
	Templates         []string                        `yaml:"templates" json:"templates"`
	MuteTimeIntervals []alertmanager.MuteTimeInterval `yaml:"mute_time_intervals,omitempty" json:"mute_time_intervals,omitempty"`
}
