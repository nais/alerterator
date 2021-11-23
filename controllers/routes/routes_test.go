package routes

import (
	"testing"

	"github.com/nais/alerterator/controllers/fixtures"
	"github.com/nais/alerterator/controllers/overrides"
	"github.com/nais/alerterator/utils"
	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"

	alertmanager "github.com/prometheus/alertmanager/config"
	model "github.com/prometheus/common/model"
)

func TestRoutes(t *testing.T) {
	t.Run("Validate that merge of config uses latest values", func(t *testing.T) {
		config := overrides.Config{}
		err := yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYaml), &config)
		assert.NoError(t, err)

		routes, err := AddOrUpdateRoute(fixtures.AlertResource(), config.Route.Routes)
		assert.NoError(t, err)
		assert.Len(t, routes, 1)

		route := routes[0]
		assert.Equal(t, []model.LabelName{"alertname", "team", "kubernetes_namespace"}, route.GroupBy)
		assert.Equal(t, "5m", route.GroupInterval.String())
		assert.Equal(t, "30s", route.GroupWait.String())
		assert.Equal(t, "aura-aura", route.Receiver)
		assert.Equal(t, "4h", route.RepeatInterval.String())
	})

	t.Run("sikre at liberator-typen støtter group_by", func(t *testing.T) {
		alert := fixtures.MinimalAlertResource()
		alert.Spec.Route = nais_io_v1.Route{GroupBy: []string{"label!"}}

		yml, err := yaml.Marshal(alert)
		assert.NoError(t, err)

		parsedAlert := &nais_io_v1.Alert{}
		err = yaml.Unmarshal(yml, parsedAlert)
		assert.NoError(t, err)

		assert.Equal(t, alert.Spec.Route.GroupBy, parsedAlert.Spec.Route.GroupBy)
	})

	t.Run("Valider at man kan legge til ny route", func(t *testing.T) {
		config := overrides.Config{}
		err := yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYaml), &config)
		assert.NoError(t, err)

		newAlert := fixtures.AlertResource()
		newAlert.Name = "newalert-does-not-exist"
		routes, err := AddOrUpdateRoute(newAlert, config.Route.Routes)
		assert.NoError(t, err)

		found := false
		for _, alert := range routes {
			if alert.Receiver == utils.GetCombinedName(newAlert) {
				found = true
			}
		}
		assert.True(t, found)
	})

	t.Run("Ensure that unset duration are 0", func(t *testing.T) {
		naisAlert := fixtures.MinimalAlertResource()
		naisAlert.Spec.Route.GroupInterval = ""
		naisAlert.Spec.Route.GroupWait = ""
		naisAlert.Spec.Route.RepeatInterval = ""
		name := utils.GetCombinedName(naisAlert)
		route, err := createNewRoute(name, naisAlert)
		assert.NoError(t, err)
		assert.Nil(t, route.GroupInterval)
		assert.Nil(t, route.GroupInterval)
		assert.Nil(t, route.RepeatInterval)
	})

	t.Run("Ensure duplicated routes are deleted", func(t *testing.T) {
		config := overrides.Config{}
		err := yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYaml), &config)
		assert.NoError(t, err)

		name := "aura-aura"
		duplicatedRoute := &alertmanager.Route{
			Receiver: name,
		}
		config.Route.Routes = append(config.Route.Routes, duplicatedRoute)
		assert.Len(t, config.Route.Routes, 2)
		config.Route.Routes = deleteDuplicates(name, config.Route.Routes)
		assert.Len(t, config.Route.Routes, 1)
	})
}
