package main

import (
	"os"
	"text/template"
)

//go:generate go run updater.go

type Service struct {
	Name          string // Resource type.
	Interface     string // Kubernetes client interface for using REST.
	Type          string // Real type of the object being modified.
	TransformFunc string // Optional name of function that will be called to copy the old object into the new.
	ClientType    string // Function to call on the clientset to get a REST client
}

var i = []Service{
	{
		Name:          "configMap",
		Interface:     "typed_core_v1.ConfigMapInterface",
		Type:          "*corev1.ConfigMap",
		ClientType:    "CoreV1().ConfigMaps",
	},
}

func main() {
	t, err := template.ParseFiles("updater.go.tpl")
	if err != nil {
		panic(err)
	}
	err = t.Execute(os.Stdout, i)
	if err != nil {
		panic(err)
	}
}