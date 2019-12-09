Alerterator
===========

Alerterator is a Kubernetes operator for managing Prometheus Alertmanager alerts. With this resource you can easily get notified (either via Slack or email) when somethings is happening with your app. 

As alerts are namespace agnostic, you don't have to have different files for each namespace you are running (although we don't recommend running you application in different namespaces). You can even make your own personal alert-resources that only notifes you!

The documentation for how to start using alerts are over at https://doc.nais.io/observability/alerts.


## Deployment

### Environment

* Kubernetes v1.11.0 or later


### Installation

You can deploy the most recent release of Alerterator by applying to your cluster:

```
kubectl apply -f ./deployment-resources/
```


## Development

It's pretty simple getting started developing, download the code, run `make build` and you should be set to Go.


### Prerequisites

* Kubectl
* The Go programming language, version 1.11 or later
* [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports)
* Docker
* Kubernetes running locally (Minikube), or a cluster to deploy too


### Code generation

In order to use the Kubernetes Go library, we need to use classes that work together with the interfaces in that library. Those classes are mostly boilerplate code, and to ensure healthy and happy developers, we use code generators for that.

When the CRD changes, or additional Kubernetes resources need to be generated, you have to run code generation:

```
make codegen-crd
```


### Testing

There are two types of tests, automatic tests in Go, or running it locally (for example in Minikube)


#### Tests

To run all the automatic tests:
```
make test
```


#### Local

```
kubectl apply -f ./hack/resources/00-namespace.yaml
kubeclt apply -f ./testing/default-app-rules.yaml
kubectl apply -f ./testing/default-alertmanager.yaml
kubeclt apply -f ./pkg/apis/alerterator/v1alpha1/alert.yaml

make build
make local

# different terminal
kubectl apply -f ./example/max_alerts.yaml
kubectl get alerts
```
