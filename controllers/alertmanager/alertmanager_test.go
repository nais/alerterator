package alertmanager

import (
	"testing"

	"github.com/nais/alerterator/controllers/fixtures"
	"github.com/nais/alerterator/controllers/overrides"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestAddOrUpdate(t *testing.T) {
	var oldConfig *overrides.Config
	err := yaml.Unmarshal([]byte(fixtures.AlertmanagerOldConfigYaml), &oldConfig)
	assert.NoError(t, err)

	var newConfig *overrides.Config
	err = yaml.Unmarshal([]byte(fixtures.AlertmanagerBaseConfigYaml), &newConfig)
	assert.NoError(t, err)

	changedConfig, err := addOrUpdate(fixtures.AlertResource(), oldConfig, newConfig)
	assert.NoError(t, err)

	data, err := yaml.Marshal(changedConfig)
	assert.NoError(t, err)
	assert.Equal(t, fixtures.AlertmanagerChangedConfigYaml, string(data))
}
