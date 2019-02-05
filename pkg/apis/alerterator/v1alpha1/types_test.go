package v1alpha1_test

import (
	"testing"

	"github.com/nais/alerterator/pkg/apis/alerterator/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestApplication_Hash(t *testing.T) {
	a1, err := v1alpha1.Alert{Spec: v1alpha1.AlertSpec{}}.Hash()
	a2, _ := v1alpha1.Alert{Spec: v1alpha1.AlertSpec{}, ObjectMeta: v1.ObjectMeta{Annotations: map[string]string{"a": "b", "team": "banana"}}}.Hash()
	a3, _ := v1alpha1.Alert{Spec: v1alpha1.AlertSpec{}, ObjectMeta: v1.ObjectMeta{Labels: map[string]string{"a": "b", "team": "banana"}}}.Hash()

	assert.NoError(t, err)
	assert.Equal(t, a1, a2, "matches, as annotations is ignored")
	assert.NotEqual(t, a2, a3, "must not match ")
}

func TestNilFix(t *testing.T) {
	alert := v1alpha1.Alert{}
	assert.Nil(t, alert.Spec.Alerts)
	alert.NilFix()
	assert.NotNil(t, alert.Spec.Receivers)
	assert.NotNil(t, alert.Spec.Alerts)
}
