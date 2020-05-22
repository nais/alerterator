package inhibitions

import (
	"fmt"
	"reflect"
	"github.com/mitchellh/mapstructure"
	v1 "github.com/nais/alerterator/pkg/apis/alerterator/v1"
)

type inhibitionConfig struct {
	Targets   map[string]string `mapstructure:"target_match" yaml:"target_match,omitempty"`
	TargetsRe map[string]string `mapstructure:"target_match_re" yaml:"target_match_re,omitempty"`
	Sources   map[string]string `mapstructure:"source_match" yaml:"source_match,omitempty"`
	SourcesRe map[string]string `mapstructure:"source_match_re" yaml:"source_match_re,omitempty"`
	Labels   []string          `mapstructure:"equal" yaml:"equal,omitempty"`
}

func createInhibitConfig(alert *v1.Alert, rule v1.InhibitRules) inhibitionConfig {
	return inhibitionConfig{
		Targets: rule.Targets,
		Sources: rule.Sources,
		Labels: append([]string{"team"}, rule.Labels...),
	}
}

func AddOrUpdateInhibition(alert *v1.Alert, alertManager map[interface{}]interface{}) ([]inhibitionConfig, error) {
	var inhibitions []inhibitionConfig
	err := mapstructure.Decode(alertManager["receivers"], &inhibitions)
	if err != nil {
		return nil, fmt.Errorf("failed while decoding map structure: %s", err)
	}

	for _, inhibitRule := range alert.Spec.InhibitRules {
		inhibitConfig := createInhibitConfig(alert, inhibitRule)
		index := getInhibitionIndex(inhibitRule.Targets, inhibitions)
		if index != -1 {
			inhibitions[index] = inhibitConfig
		} else {
			inhibitions = append(inhibitions, inhibitConfig)
		}
	}

	return inhibitions, nil
}

func getInhibitionIndex(target map[string]string, inhibitions []inhibitionConfig) int {
	for i := 0; i < len(inhibitions); i++ {
		inhibition := inhibitions[i]

		if reflect.DeepEqual(inhibition.Targets, target) {
			return i
		}
	}
	return -1
}

func DeleteInhibition(alert *v1.Alert, alertManager map[interface{}]interface{}) error {
	var inhibitions []inhibitionConfig
	err := mapstructure.Decode(alertManager["inhibit_rules"], &inhibitions)
	if err != nil {
		return fmt.Errorf("failed while decoding map structure: %s", err)
	}

	for _, inhibitRule := range alert.Spec.InhibitRules {
		index := getInhibitionIndex(inhibitRule.Targets, inhibitions)
		if index != -1 {
			inhibitions = append(inhibitions[:index], inhibitions[index+1:]...)
		}
	}

	alertManager["inhibit_rules"] = inhibitions

	return nil
}
