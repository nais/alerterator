package rules

import (
	"testing"

	"github.com/nais/alerterator/api/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestConfigMapUpdater(t *testing.T) {
	t.Run("Test that alerts get added", func(t *testing.T) {
		configMap, err := addOrUpdateAlerts(fixtures.AlertResource, fixtures.ConfigMapBeforeAlerts)
		assert.NoError(t, err)
		assert.Equal(t, fixtures.ExpectedConfigMapAfterAlerts.Data["aura.yml"], configMap.Data["aura.yml"])
	})
}
