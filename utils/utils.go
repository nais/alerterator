package utils

import (
	"fmt"

	"github.com/nais/liberator/pkg/apis/nais.io/v1"
)

func GetCombinedName(alert *nais_io_v1.Alert) string {
	return fmt.Sprintf("%s-%s", alert.Namespace, alert.Name)
}
