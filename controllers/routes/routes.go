package routes

import (
	"github.com/nais/alerterator/utils"

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
	var groupWait, groupInterval, repeatInterval *model.Duration

	if len(alert.Spec.Route.GroupWait) > 0 {
		gw, err := model.ParseDuration(alert.Spec.Route.GroupWait)
		if err != nil {
			return nil, err
		}
		groupWait = &gw
	}
	if len(alert.Spec.Route.GroupInterval) > 0 {
		gi, err := model.ParseDuration(alert.Spec.Route.GroupInterval)
		if err != nil {
			return nil, err
		}

		groupInterval = &gi
	}
	if len(alert.Spec.Route.RepeatInterval) > 0 {
		ri, err := model.ParseDuration(alert.Spec.Route.RepeatInterval)
		if err != nil {
			return nil, err
		}
		repeatInterval = &ri
	}

	var groupBy []model.LabelName
	for _, v := range alert.Spec.Route.GroupBy {
		groupBy = append(groupBy, model.LabelName(v))
	}

	return &alertmanager.Route{
		GroupBy:        groupBy,
		GroupInterval:  groupInterval,
		GroupWait:      groupWait,
		RepeatInterval: repeatInterval,
		Receiver:       name,
		Continue:       true,
		Match: map[string]string{
			"alert": name,
		},
	}, nil
}

func deleteDuplicates(name string, routes []*alertmanager.Route) []*alertmanager.Route {
	var indices []int
	for i := range routes {
		if routes[i].Receiver == name {
			indices = append(indices, i)
		}
	}

	if len(indices) > 1 {
		for i := 1; i < len(indices); i++ {
			routes = append(routes[:i], routes[i+1:]...)
		}
	}

	return routes
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

	routes = deleteDuplicates(name, routes)

	return routes, nil
}

func DeleteRoute(alert *naisiov1.Alert, routes []*alertmanager.Route) []*alertmanager.Route {
	name := utils.GetCombinedName(alert)
	if i := getRouteIndex(name, routes); i != -1 {
		routes = append(routes[:i], routes[i+1:]...)
	}

	return routes
}
