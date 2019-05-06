package updater

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/mitchellh/mapstructure"
	"github.com/nais/alerterator/pkg/apis/alerterator/v1alpha1"
)

type routeConfig struct {
	Receiver string            `mapstructure:"receiver" yaml:"receiver"`
	Continue bool              `mapstructure:"continue" yaml:"continue"`
	Match    map[string]string `mapstructure:"match" yaml:"match"`
}

type routesConfig struct {
	GroupBy        []string      `mapstructure:"group_by" yaml:"group_by"`
	GroupWait      string        `mapstructure:"group_wait" yaml:"group_wait"`
	GroupInterval  string        `mapstructure:"group_interval" yaml:"group_interval"`
	RepeatInterval string        `mapstructure:"repeat_interval" yaml:"repeat_interval"`
	Receiver       string        `mapstructure:"receiver" yaml:"receiver"`
	Routes         []routeConfig `mapstructure:"routes" yaml:"routes"`
}

func missingAlertRoute(alertName string, routes []routeConfig) bool {
	for i := 0; i < len(routes); i++ {
		route := routes[i]
		if route.Receiver == alertName {
			return false
		}
	}
	return true
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

func AddOrUpdateRoutes(alert *v1alpha1.Alert, alertManager map[interface{}]interface{}) error {
	var route routesConfig
	err := mapstructure.Decode(alertManager["route"], &route)
	if err != nil {
		return fmt.Errorf("failed while decoding map structure: %s", err)
	}

	if missingAlertRoute(alert.Name, route.Routes) {
		log.Infof("Route missing for %s", alert.Name)
		routes := routeConfig{
			Receiver: alert.Name,
			Continue: true,
			Match: map[string]string{
				"alert": alert.Name,
			},
		}
		route.Routes = append(route.Routes, routes)
		alertManager["route"] = route
	}

	return nil
}

func DeleteRoute(alert *v1alpha1.Alert, alertManager map[interface{}]interface{}) error {
	var route routesConfig
	err := mapstructure.Decode(alertManager["route"], &route)
	if err != nil {
		return fmt.Errorf("failed while decoding map structure: %s", err)
	}

	index := getAlertRouteIndex(alert.Name, route.Routes)
	if index == -1 {
		log.Infof("No route with the name %s", alert.Name)
		return nil
	}
	log.Info(route.GroupWait)
	route.Routes = append(route.Routes[:index], route.Routes[index+1:]...)
	alertManager["route"] = route

	return nil
}
