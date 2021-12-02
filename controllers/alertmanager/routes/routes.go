package routes

import (
	"github.com/nais/alerterator/utils"

	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"

	alertmanager "github.com/prometheus/alertmanager/config"
	model "github.com/prometheus/common/model"
)

func getRouteIndexByName(receiver string, routes []*alertmanager.Route) int {
	for i := range routes {
		if routes[i].Receiver == receiver {
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

	var groupBy []string
	for _, v := range alert.Spec.Route.GroupBy {
		groupBy = append(groupBy, v)
	}

	return &alertmanager.Route{
		GroupByStr:     groupBy,
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

func AddOrUpdate(alert *naisiov1.Alert, routes []*alertmanager.Route) ([]*alertmanager.Route, error) {
	name := utils.GetCombinedName(alert)
	alertRoute, err := createNewRoute(name, alert)
	if err != nil {
		return nil, err
	}

	if i := getRouteIndexByName(name, routes); i != -1 {
		routes[i] = alertRoute
	} else {
		routes = append(routes, alertRoute)
	}

	routes = deleteDuplicates(name, routes)

	return routes, nil
}

func Delete(alert *naisiov1.Alert, routes []*alertmanager.Route) []*alertmanager.Route {
	name := utils.GetCombinedName(alert)
	if i := getRouteIndexByName(name, routes); i != -1 {
		routes = append(routes[:i], routes[i+1:]...)
	}

	return routes
}
