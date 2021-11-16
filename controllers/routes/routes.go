package routes

import (
	"fmt"

	"github.com/nais/alerterator/utils"

	"github.com/mitchellh/mapstructure"
	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"

	alertmanager "github.com/prometheus/alertmanager/config"
	model "github.com/prometheus/common/model"
)

func getRouteIndex(alertName string, routes []*alertmanager.Route) int {
	for i := range routes {
		if routes[i].Receiver == alertName {
			return i
		}
	}

	return -1
}

func createNewRoute(name string, alert *naisiov1.Alert) (*alertmanager.Route, error) {
	groupWait, err := model.ParseDuration(alert.Spec.Route.GroupWait)
	if err != nil {
		return nil, err
	}
	groupInterval, err := model.ParseDuration(alert.Spec.Route.GroupInterval)
	if err != nil {
		return nil, err
	}
	repeatInterval, err := model.ParseDuration(alert.Spec.Route.RepeatInterval)
	if err != nil {
		return nil, err
	}

	var groupBy []model.LabelName
	for _, v := range alert.Spec.Route.GroupBy {
		groupBy = append(groupBy, model.LabelName(v))
	}

	return &alertmanager.Route{
		GroupBy:        groupBy,
		GroupInterval:  &groupInterval,
		GroupWait:      &groupWait,
		RepeatInterval: &repeatInterval,
		Receiver:       name,
		Continue:       true,
		Match: map[string]string{
			"alert": name,
		},
	}, nil
}

func AddOrUpdateRoute(alert *naisiov1.Alert, routes []*alertmanager.Route) ([]*alertmanager.Route, error) {
	alertName := utils.GetCombinedName(alert)
	alertRoute, err := createNewRoute(alertName, alert)
	if err != nil {
		return nil, err
	}

	if i := getRouteIndex(alertName, routes); i != -1 {
		routes[i] = alertRoute
	} else {
		routes = append(routes, alertRoute)
	}

	return routes, nil
}

func getAlertRouteIndex(alertName string, routes []*alertmanager.Route) int {
	for i := 0; i < len(routes); i++ {
		route := routes[i]
		if route.Receiver == alertName {
			return i
		}
	}
	return -1
}

func DeleteRoute(alert *naisiov1.Alert, alertManager map[interface{}]interface{}) error {
	var route alertmanager.Route
	err := mapstructure.Decode(alertManager["route"], &route)
	if err != nil {
		return fmt.Errorf("failed while decoding map structure: %s", err)
	}

	index := getAlertRouteIndex(utils.GetCombinedName(alert), route.Routes)
	if index != -1 {
		route.Routes = append(route.Routes[:index], route.Routes[index+1:]...)
		alertManager["route"] = route
	}

	return nil
}
