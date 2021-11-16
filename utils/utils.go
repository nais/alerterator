package utils

import (
	"fmt"

	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
)

func GetCombinedName(alert *naisiov1.Alert) string {
	return fmt.Sprintf("%s-%s", alert.Namespace, alert.Name)
}
