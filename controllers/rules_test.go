package controllers

import (
	"github.com/nais/alerterator/controllers/rules"
	"testing"

	"github.com/nais/alerterator/controllers/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestConfigMapUpdater(t *testing.T) {
	t.Run("Test that alerts get added", func(t *testing.T) {
		alert := fixtures.AlertResource()
		configMap, err := rules.AddOrUpdateAlert(alert, *fixtures.ConfigMapBeforeAlerts())
		assert.NoError(t, err)
		assert.Equal(t, fixtures.ExpectedConfigMapAfterAlerts().Data["aura.yml"], configMap.Data["aura.yml"])
	})
}
