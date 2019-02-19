package api

import (
	"github.com/nais/alerterator/api/fixtures"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
)

func TestConfigMapDeleter(t *testing.T) {
	t.Run("Test for error if alertmanager.yml is missing", func(t *testing.T) {
		_, err := deleteReceivers(nil, &v1.ConfigMap{})
		assert.Error(t, err)
	})

	t.Run("Test that receiver and route is deleted correctly", func(t *testing.T) {
		configMap, err := deleteReceivers(fixtures.AlertResource, fixtures.ConfigMapBeforeDelete)
		assert.NoError(t, err)
		assert.Equal(t, fixtures.ExpectedConfigMapAfterDelete.Data["alertmanager.yml"], configMap.Data["alertmanager.yml"])
	})
}
