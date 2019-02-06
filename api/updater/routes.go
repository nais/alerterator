package updater

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/mitchellh/mapstructure"
	"github.com/nais/alerterator/pkg/apis/alerterator/v1alpha1"
)

type matchConfig struct {
	Team string `mapstructure:"team"`
}

type routeConfig struct {
	Receiver string      `mapstructure:"receiver"`
	Continue bool        `mapstructure:"continue"`
	Match    matchConfig `mapstructure:"match"`
}

type routesConfig struct {
	GroupBy        []string      `mapstructure:"group_by"`
	GroupWait      string        `mapstructure:"group_wait"`
	GroupInterval  string        `mapstructure:"group_interval"`
	RepeatInterval string        `mapstructure:"repeat_interval"`
	Receiver       string        `mapstructure:"receiver"`
	Routes         []routeConfig `mapstructure:"routes"`
}

func missingAlertRoute(alert string, routes []routeConfig) bool {
	for i := 0; i < len(routes); i++ {
		route := routes[i]
		if route.Receiver == alert {
			return false
		}
	}
	return true
}

func AddOrUpdateRoutes(alert *v1alpha1.Alert, alertManager map[interface{}]interface{}) error {
	var route routesConfig
	err := mapstructure.Decode(alertManager["route"], &route)
	if err != nil {
		return fmt.Errorf("failed while decoding map structure: %s", err)
	}

	if missingAlertRoute(alert.Name, route.Routes) {
		glog.Infof("Route missing for %s", alert.Name)
		routes := routeConfig{
			Receiver: alert.Name,
			Continue: true,
			Match: matchConfig{
				Team: alert.GetTeamName(),
			},
		}
		route.Routes = append(route.Routes, routes)
		alertManager["route"] = route
	}

	return nil
}
