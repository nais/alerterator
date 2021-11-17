package inhibitions

import (
	"testing"

	"github.com/nais/alerterator/controllers/fixtures"
	alertmanager "github.com/prometheus/alertmanager/config"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestRoutes(t *testing.T) {
	t.Run("Labels should always have team", func(t *testing.T) {
		alert := fixtures.AlertResource()
		config := alertmanager.Config{}
		err := yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYaml), &config)
		assert.NoError(t, err)

		inhibitions, err := AddOrUpdateInhibition(alert, config.InhibitRules)
		assert.NoError(t, err)
		rule := inhibitions[len(inhibitions)-1]
		assert.Contains(t, rule.Equal, model.LabelName("team"))
	})

	t.Run("Simple validation that new config don't override old config", func(t *testing.T) {
		alert := fixtures.AlertResource()
		config := alertmanager.Config{}
		err := yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYaml), &config)
		assert.NoError(t, err)

		inhibitionConfig, err := AddOrUpdateInhibition(alert, config.InhibitRules)
		assert.NoError(t, err)
		assert.Len(t, inhibitionConfig, 3)

	})

	t.Run("Simple deletion validation", func(t *testing.T) {
		alert := fixtures.AlertResource()
		config := alertmanager.Config{}
		err := yaml.Unmarshal([]byte(fixtures.AlertmanagerConfigYaml), &config)
		assert.NoError(t, err)

		inhibitionConfig, err := AddOrUpdateInhibition(alert, config.InhibitRules)
		assert.NoError(t, err)
		assert.Len(t, inhibitionConfig, 3)

		inhibitRules := DeleteInhibition(alert, config.InhibitRules)
		assert.Len(t, inhibitRules, 2)
	})
}
