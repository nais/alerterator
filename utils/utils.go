package utils

import (
	"fmt"
	"github.com/nais/alerterator/pkg/apis/alerterator/v1"
)

func GetCombinedName(alert *v1.Alert) string {
	return fmt.Sprintf("%s-%s", alert.Namespace, alert.Name)
}
