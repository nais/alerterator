package routes

import (
	"testing"

	"github.com/nais/alerterator/controllers/fixtures"
	"github.com/nais/alerterator/utils"
	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"

	alertmanager "github.com/prometheus/alertmanager/config"
	model "github.com/prometheus/common/model"
)

func TestRoutes(t *testing.T) {
	t.Run("Validate that merge of config uses latest values", func(t *testing.T) {
		config := alertmanager.Config{}
		err := yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYaml), &config)
		assert.NoError(t, err)

		latestConfig := alertmanager.Config{}
		err = yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYamlDifferentRoutes), &latestConfig)
		assert.NoError(t, err)

		routes, err := AddOrUpdateRoute(fixtures.AlertResource(), config, latestConfig)
		assert.NoError(t, err)

		assert.Len(t, routes, 1)
		route := routes[0]
		assert.Equal(t, []model.LabelName{"alertname", "team", "kubernetes_namespace"}, route.GroupBy)
		assert.Equal(t, "50m", route.GroupInterval.String())
		assert.Equal(t, "100s", route.GroupWait.String())
		assert.Equal(t, "default-receiver", route.Receiver)
		assert.Equal(t, "10h", route.RepeatInterval.String())
	})

	t.Run("Valider at man kan sette route config per route", func(t *testing.T) {
		config := alertmanager.Config{}
		err := yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYaml), &config)
		assert.NoError(t, err)

		routes, err := AddOrUpdateRoute(fixtures.AlertResource(), config, config)
		assert.NoError(t, err)

		teamRoute := routes[1]
		assert.Equal(t, "5m", teamRoute.GroupInterval)
		assert.Equal(t, "30s", teamRoute.GroupWait)
		assert.Equal(t, "4h", teamRoute.RepeatInterval)
	})

	t.Run("Valider at group-by kommer igjennom parsing", func(t *testing.T) {
		config := alertmanager.Config{}
		err := yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYaml), &config)
		assert.NoError(t, err)

		alert := fixtures.MinimalAlertResource()
		alert.Spec.Route = nais_io_v1.Route{GroupBy: []string{"slack_channel"}}
		routeConfig, err := AddOrUpdateRoute(alert, config, config)
		assert.NoError(t, err)
		assert.Len(t, routeConfig, 1)

		groupBys := routeConfig[0].GroupBy
		assert.Equal(t, "slack_channel", groupBys[0])
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
		config := alertmanager.Config{}
		err := yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYaml), &config)
		assert.NoError(t, err)

		templateConfig := alertmanager.Config{}
		err = yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYamlDifferentRoutes), &templateConfig)
		assert.NoError(t, err)

		newAlert := fixtures.AlertResource()
		newAlert.Name = "newalert-does-not-exist"
		routes, err := AddOrUpdateRoute(newAlert, config, templateConfig)
		assert.NoError(t, err)

		found := false
		for _, alert := range routes {
			if alert.Receiver == utils.GetCombinedName(newAlert) {
				found = true
			}
		}
		assert.True(t, found)
	})

	t.Run("Valider at man kan endre eksisterende route", func(t *testing.T) {
		config := alertmanager.Config{}
		err := yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYaml), &config)
		assert.NoError(t, err)

		templateConfig := alertmanager.Config{}
		err = yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYamlDifferentRoutes), &templateConfig)
		assert.NoError(t, err)

		updatedAlert := fixtures.AlertResource()
		updatedAlert.Spec.Route.GroupBy = []string{"updated-group-by"}
		routes, err := AddOrUpdateRoute(updatedAlert, config, templateConfig)
		assert.NoError(t, err)

		assert.Equal(t, "updated-group-by", routes[1].GroupBy[0])
		assert.Len(t, routes[1].GroupBy, 1)
		assert.Len(t, routes, 2)
	})
}
