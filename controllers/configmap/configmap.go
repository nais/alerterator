package configmap

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Get(ctx context.Context, client client.Client, namespacedName types.NamespacedName) (v1.ConfigMap, error) {
	var configMap v1.ConfigMap
	err := client.Get(ctx, namespacedName, &configMap)
	if err != nil {
		return v1.ConfigMap{}, fmt.Errorf("failing while retrieving %s configMap: %s", namespacedName.Name, err)
	}

	return configMap, nil
}

func GetAndUnmarshal(ctx context.Context, client client.Client, namespacedName types.NamespacedName, configName string, out interface{}) error {
	configMap, err := Get(ctx, client, namespacedName)
	if err != nil {
		return err
	}

	if configMap.Data == nil {
		return fmt.Errorf("alerterator is not properly set up, %s is empty", namespacedName)
	}

	err = yaml.Unmarshal([]byte(configMap.Data[configName]), out)
	if err != nil {
		return fmt.Errorf("failed while unmarshling %s: %s", configName, err)
	}

	return nil
}

func UpdateData(ctx context.Context, client client.Client, namespacedName types.NamespacedName, configName string, in string) error {
	configMap, err := Get(ctx, client, namespacedName)
	if err != nil {
		return err
	}

	if configMap.Data == nil {
		configMap.Data = make(map[string]string)
	}

	configMap.Data[configName] = in
	err = client.Update(ctx, &configMap)
	if err != nil {
		return fmt.Errorf("failed while updating %s: %s", configName, err)
	}

	return nil
}

func MarshalAndUpdateData(ctx context.Context, client client.Client, namespacedName types.NamespacedName, configName string, in interface{}) error {
	data, err := yaml.Marshal(in)
	if err != nil {
		return fmt.Errorf("failed while marshaling: %s", err)
	}

	return UpdateData(ctx, client, namespacedName, configName, string(data))
}

func DeleteFileFromData(ctx context.Context, client client.Client, namespacedName types.NamespacedName, fileName string) error {
	var configMap v1.ConfigMap
	err := client.Get(ctx, namespacedName, &configMap)
	if err != nil {
		return fmt.Errorf("failing while retrieving %s configMap: %s", namespacedName, err)
	}
	delete(configMap.Data, fileName)

	err = client.Update(ctx, &configMap)
	if err != nil {
		return fmt.Errorf("failed while updating %s configMaps: %s", namespacedName, err)
	}

	return nil
}
