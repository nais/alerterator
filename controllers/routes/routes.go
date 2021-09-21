package routes

import (
	"fmt"

	"github.com/nais/alerterator/utils"

	"github.com/mitchellh/mapstructure"
	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
)

type routeConfig struct {
	Receiver       string            `mapstructure:"receiver" yaml:"receiver"`
	Continue       bool              `mapstructure:"continue" yaml:"continue"`
	Match          map[string]string `mapstructure:"match" yaml:"match"`
	GroupBy        []string          `mapstructure:"group_by" yaml:"group_by,omitempty"`
	GroupWait      string            `mapstructure:"group_wait" yaml:"group_wait,omitempty"`
	GroupInterval  string            `mapstructure:"group_interval" yaml:"group_interval,omitempty"`
	RepeatInterval string            `mapstructure:"repeat_interval" yaml:"repeat_interval,omitempty"`
}

type Config struct {
	GroupBy        []string      `mapstructure:"group_by" yaml:"group_by"`
	GroupWait      string        `mapstructure:"group_wait" yaml:"group_wait"`
	GroupInterval  string        `mapstructure:"group_interval" yaml:"group_interval"`
	RepeatInterval string        `mapstructure:"repeat_interval" yaml:"repeat_interval"`
	Receiver       string        `mapstructure:"receiver" yaml:"receiver"`
	Routes         []routeConfig `mapstructure:"routes" yaml:"routes"`
}

func getRouteIndex(alertName string, routes []routeConfig) int {
	for i := range routes {
		if routes[i].Receiver == alertName {
			return i
		}
	}

	return -1
}

func AddOrUpdateRoute(alert *naisiov1.Alert, currentConfig, latestConfig map[interface{}]interface{}) (Config, error) {
	var routes Config
	err := mapstructure.Decode(currentConfig["route"], &routes)
	if err != nil {
		return Config{}, fmt.Errorf("failed while decoding map structure: %s", err)
	}

	alertName := utils.GetCombinedName(alert)

	routeConfig := routeConfig{
		GroupBy:        alert.Spec.Route.GroupBy,
		GroupInterval:  alert.Spec.Route.GroupInterval,
		GroupWait:      alert.Spec.Route.GroupWait,
		RepeatInterval: alert.Spec.Route.RepeatInterval,
		Receiver:       alertName,
		Continue:       true,
		Match: map[string]string{
			"alert": alertName,
		},
	}

	if i := getRouteIndex(alertName, routes.Routes); i != -1 {
		routes.Routes[i] = routeConfig
	} else {
		routes.Routes = append(routes.Routes, routeConfig)
	}

	var latestRoutes Config
	err = mapstructure.Decode(latestConfig["route"], &latestRoutes)
	if err != nil {
		return Config{}, fmt.Errorf("failed while decoding map structure: %s", err)
	}

	latestRoutes.Routes = routes.Routes

	return latestRoutes, nil
}

func getAlertRouteIndex(alertName string, routes []routeConfig) int {
	for i := 0; i < len(routes); i++ {
		route := routes[i]
		if route.Receiver == alertName {
			return i
		}
	}
	return -1
}

func DeleteRoute(alert *naisiov1.Alert, alertManager map[interface{}]interface{}) error {
	var route Config
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
