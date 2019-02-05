package updater

import (
	"fmt"
	"github.com/nais/alerterator/pkg/apis/alerterator/v1alpha1"
	"gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
)

type AlertGroups struct {
	Groups []AlertGroup
}

type AlertGroup struct {
	Name  string
	Rules []v1alpha1.Rule
}

func addTeamLabel(rules []v1alpha1.Rule, teamName string) {
	for i := range rules {
		if rules[i].Labels == nil {
			rules[i].Labels = make(map[string]string)
		}

		rules[i].Labels["team"] = teamName
	}
}

func AddOrUpdateAlerts(alert *v1alpha1.Alert, configMap *v1.ConfigMap) (*v1.ConfigMap, error) {
	addTeamLabel(alert.Spec.Alerts, alert.GetTeamName())
	alertGroup := AlertGroup{Name: alert.Name, Rules: alert.Spec.Alerts}
	alertGroups := AlertGroups{Groups: []AlertGroup{alertGroup}}

	alertGroupYamlBytes, err := yaml.Marshal(alertGroups)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %v to yaml\n", alertGroup)
	}

	if configMap.Data == nil {
		configMap.Data = make(map[string]string)
	}

	configMap.Data[alert.Name+".yml"] = string(alertGroupYamlBytes)

	return configMap, nil
}
