package routes

import (
	"github.com/nais/alerterator/api/fixtures"
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
}
