package routes

import (
	"alerterator/controllers/fixtures"
	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestRoutes(t *testing.T) {
	t.Run("Validate that merge of config uses latest values", func(t *testing.T) {
		config := make(map[interface{}]interface{})
		err := yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYaml), config)
		assert.NoError(t, err)

		latestConfig := make(map[interface{}]interface{})
		err = yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYamlDifferentRoutes), latestConfig)
		assert.NoError(t, err)

		routesConfig, err := AddOrUpdateRoute(fixtures.AlertResource, config, latestConfig)
		assert.NoError(t, err)

		assert.Equal(t, []string{"alertname", "team", "kubernetes_namespace"}, routesConfig.GroupBy)
		assert.Equal(t, "50m", routesConfig.GroupInterval)
		assert.Equal(t, "100s", routesConfig.GroupWait)
		assert.Equal(t, "default-receiver", routesConfig.Receiver)
		assert.Equal(t, "10h", routesConfig.RepeatInterval)
	})

	t.Run("Valider at man kan sette route config per route", func(t *testing.T) {
		config := make(map[interface{}]interface{})
		err := yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYaml), config)
		assert.NoError(t, err)

		routesConfig, err := AddOrUpdateRoute(fixtures.AlertResource, config, config)
		assert.NoError(t, err)

		teamRoute := routesConfig.Routes[1]
		assert.Equal(t, "5m", teamRoute.GroupInterval)
		assert.Equal(t, "30s", teamRoute.GroupWait)
		assert.Equal(t, "4h", teamRoute.RepeatInterval)
	})

	t.Run("Valider at group-by kommer igjennom parsing", func(t *testing.T) {
		config := make(map[interface{}]interface{})
		err := yaml.Unmarshal([]byte(fixtures.AlertWithGroupBy), config)
		assert.NoError(t, err)

		routeConfig, err := AddOrUpdateRoute(fixtures.AlertResource, config, config)
		assert.NoError(t, err)

		groupBys := routeConfig.GroupBy
		assert.Equal(t, "slack_channel", groupBys[0])
	})

	t.Run("sikre at liberator-typen støtter group_by", func(t *testing.T) {
		alert := fixtures.MinimalAlertResource
		alert.Spec.Route = nais_io_v1.Route{GroupBy: []string{"label!"}}

		yml, err := yaml.Marshal(alert)
		assert.NoError(t, err)

		parsedAlert := &nais_io_v1.Alert{}
		err = yaml.Unmarshal(yml, parsedAlert)
		assert.NoError(t, err)

		assert.Equal(t, alert.Spec.Route.GroupBy, parsedAlert.Spec.Route.GroupBy)
	})
}
