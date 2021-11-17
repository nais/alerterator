package inhibitions

import (
	"fmt"
	"reflect"
	"regexp"

	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	alertmanager "github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/pkg/labels"
	"github.com/prometheus/common/model"
)

func createMatchers(source map[string]string) (matchers []*labels.Matcher) {
	for k, v := range source {
		matcher := labels.Matcher{
			Name:  k,
			Value: v,
		}
		matchers = append(matchers, &matcher)
	}

	return
}

func createMatchRegexps(source map[string]string) (map[string]alertmanager.Regexp, error) {
	matchRegexps := make(map[string]alertmanager.Regexp)
	for k, v := range source {
		regex, err := regexp.Compile(v)
		if err != nil {
			return nil, err
		}

		regexp := alertmanager.Regexp{
			Regexp: regex,
		}
		matchRegexps[k] = regexp
	}

	return matchRegexps, nil
}

func createInhibitRule(rule naisiov1.InhibitRules) (*alertmanager.InhibitRule, error) {
	var equals []model.LabelName
	for _, label := range rule.Labels {
		equals = append(equals, model.LabelName(label))
	}
	equals = append(equals, model.LabelName("team"))

	targetMatchRE, err := createMatchRegexps(rule.TargetsRegex)
	if err != nil {
		return nil, err
	}
	sourceMatchRE, err := createMatchRegexps(rule.SourcesRegex)
	if err != nil {
		return nil, err
	}

	return &alertmanager.InhibitRule{
		TargetMatchers: createMatchers(rule.Targets),
		TargetMatchRE:  targetMatchRE,
		SourceMatchers: createMatchers(rule.Sources),
		SourceMatchRE:  sourceMatchRE,
		Equal:          equals,
	}, nil
}

func AddOrUpdateInhibition(alert *naisiov1.Alert, inhibitions []*alertmanager.InhibitRule) ([]*alertmanager.InhibitRule, error) {
	for _, ir := range alert.Spec.InhibitRules {
		inhibitRule, err := createInhibitRule(ir)
		if err != nil {
			return nil, err
		}

		if i := getInhibitionIndex(ir.Targets, inhibitions); i != -1 {
			inhibitions[i] = inhibitRule
		} else {
			inhibitions = append(inhibitions, inhibitRule)
		}
	}

	fmt.Println(inhibitions)
	return inhibitions, nil
}

func getInhibitionIndex(target map[string]string, inhibitions []*alertmanager.InhibitRule) int {
	for i := 0; i < len(inhibitions); i++ {
		inhibition := inhibitions[i]

		if reflect.DeepEqual(inhibition.TargetMatchers, target) {
			return i
		}
	}
	return -1
}

func DeleteInhibition(alert *naisiov1.Alert, inhibitions []*alertmanager.InhibitRule) []*alertmanager.InhibitRule {
	for _, ir := range alert.Spec.InhibitRules {
		if i := getInhibitionIndex(ir.Targets, inhibitions); i != -1 {
			inhibitions = append(inhibitions[:i], inhibitions[i+1:]...)
		}
	}

	return inhibitions
}
