package inhibitions

import (
	"alerterator/controllers/fixtures"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestRoutes(t *testing.T) {
	t.Run("Labels should always have team", func(t *testing.T) {
		alert := fixtures.AlertResource()
		config := make(map[interface{}]interface{})
		err := yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYaml), config)
		assert.NoError(t, err)

		inhibitionConfig, err := AddOrUpdateInhibition(alert, config)
		assert.NoError(t, err)
		rule := inhibitionConfig[len(inhibitionConfig)-1]
		assert.Contains(t, rule.Labels, "team")
	})

	t.Run("Simple validation that new config don't override old config", func(t *testing.T) {
		alert := fixtures.AlertResource()
		config := make(map[interface{}]interface{})
		err := yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYaml), config)
		assert.NoError(t, err)

		inhibitionConfig, err := AddOrUpdateInhibition(alert, config)
		assert.NoError(t, err)
		assert.Len(t, inhibitionConfig, 3)

	})

	t.Run("Simple deletion validation", func(t *testing.T) {
		alert := fixtures.AlertResource()
		config := make(map[interface{}]interface{})
		err := yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYaml), config)
		assert.NoError(t, err)

		inhibitionConfig, err := AddOrUpdateInhibition(alert, config)
		assert.NoError(t, err)
		assert.Len(t, inhibitionConfig, 3)

		err = DeleteInhibition(alert, config)
		assert.NoError(t, err)
		assert.Len(t, config["inhibit_rules"], 2)

	})
}
