package configmap

import (
	"context"
	"fmt"

	"github.com/nais/alerterator/controllers/overrides"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Get(ctx context.Context, namespacedName types.NamespacedName, client client.Client, configName string) (*overrides.Config, error) {
	var configMap v1.ConfigMap
	err := client.Get(ctx, namespacedName, &configMap)
	if err != nil {
		return nil, fmt.Errorf("failing while retrieving %s configMap: %s", namespacedName.Name, err)
	}

	if configMap.Data == nil {
		return nil, fmt.Errorf("alerterator is not properly set up, %s is empty", namespacedName)
	}

	config := overrides.Config{}
	err = yaml.Unmarshal([]byte(configMap.Data[configName]), &config)
	if err != nil {
		return nil, fmt.Errorf("failed while unmarshling %s: %s", configName, err)
	}

	return &config, nil
}

func Update(ctx context.Context, namespacedName types.NamespacedName, config *overrides.Config, client client.Client, configName string) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed while marshaling: %s", err)
	}

	var configMap v1.ConfigMap
	err = client.Get(ctx, namespacedName, &configMap)
	if err != nil {
		return fmt.Errorf("failing while retrieving %s configMap: %s", namespacedName, err)
	}

	configMap.Data[configName] = string(data)
	err = client.Update(ctx, &configMap)
	if err != nil {
		return fmt.Errorf("failed while updating %s: %s", configName, err)
	}

	return nil
}
