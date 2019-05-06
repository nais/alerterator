package api

import (
	"testing"

	"github.com/nais/alerterator/api/fixtures"
	"github.com/nais/alerterator/api/updater"
	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
)

func TestConfigMapUpdater(t *testing.T) {
	t.Run("Test that alerts get added", func(t *testing.T) {
		configMap, err := updater.AddOrUpdateAlerts(fixtures.AlertResource, fixtures.ConfigMapBeforeAlerts)
		assert.NoError(t, err)
		assert.Equal(t, fixtures.ExpectedConfigMapAfterAlerts.Data["aura.yml"], configMap.Data["aura.yml"])
	})

	t.Run("Test for error if alertmanager.yml is missing", func(t *testing.T) {
		_, err := deleteReceivers(nil, &v1.ConfigMap{})
		assert.Error(t, err)
	})

	t.Run("Test that receiver and route is added correctly", func(t *testing.T) {
		configMap, err := addOrUpdateReceivers(fixtures.AlertResource, fixtures.ConfigMapBeforeAdd)
		assert.NoError(t, err)
		assert.Equal(t, fixtures.ExpectedConfigMapAfterReceivers.Data["alertmanager.yml"], configMap.Data["alertmanager.yml"])
	})

	t.Run("Test that Naisd-alerts is not affected by Alerterator", func(t *testing.T) {
		configMap, err := addOrUpdateReceivers(fixtures.AlertResource, fixtures.ConfigMapMixed)
		assert.NoError(t, err)
		assert.Equal(t, fixtures.ConfigMapMixed.Data["alertmanager.yml"], configMap.Data["alertmanager.yml"])
	})
}
