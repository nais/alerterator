Alerterator
===========

Alerterator is a Kubernetes operator for managing Prometheus Alertmanager alerts. With this resource you can easily get notified (either via Slack or email) when somethings is happening with your app. 

As alerts are namespace agnostic, you don't have to have different files for each namespace you are running (although we don't recommend running you application in different namespaces). You can even make your own personal alert-resources that only notifes you!

## Spec

| Parameter | Description | Default | Required |
| --------- | ----------- | ------- | :--------: |
| metadata.name | Name for the group of alerts | | x |
| metadata.labels.team | [mailnick/tag](https://github.com/nais/doc/blob/master/content/getting-started/teamadministration.md) | | x |
| spec.receivers.slack.channel | Slack channel to send notifications to | | |
| spec.receivers.slack.preprend_text | Text to prepend every Slack-message (for ex. @here) | | |
| spec.receivers.email.to | The email address to send notifications to| | |
| spec.receivers.email.send_resolved | Whether or not to notify about resolved alerts | | false |
| spec.alerts[].description | Simple description of the triggered alert | | x |
| spec.alerts[].expr | Prometheus expression that triggers an alert | | x |
| spec.alerts[].for | Duration before the alert should trigger | | x |
| spec.alerts[].action | How to resolve this alert | | x |
| spec.alerts[].documentation | URL for docmentation for this alert| | |
| spec.alerts[].sla | Time before the alert should be resolved| | |
| spec.alerts[].severity | Alert level for Slack-messages| | Error |


See [example directory](/example/) for an example-alert.


### Tips

You can also use `annotations` and `labels` from the Prometheus-`expr` result.

For example:
```
{{ $labels.node }} is marked as unschedulable
```

turns into the following when posted to Slack/email:
```
b27apvl00178.preprod.local is marked as unschedulable
```


## Migrating from Naisd

It's pretty straight forward to move alerts from Naisd to Alerterator, and the most notable difference is the removal of name/alert and that annotation-fields has been move to the top-level.

```
alerts:
- alert: appNotAvailable
  expr: kube_deployment_status_replicas_unavailable{deployment="app-name"} > 0
  for: 5m
  annotations:
    action: Read app logs(kubectl logs appname). Read Application events (kubectl descibe deployment appname)
    severity: Warning
```

should be transformed to

```
alerts:
- expr: kube_deployment_status_replicas_unavailable{deployment="app-name"} > 0
  for: 5m
  description: It looks like the app is down
  action: Read app logs(kubectl logs appname). Read Application events (kubectl descibe deployment appname)
  severity: Warning
```

Check out the complete [spec](/#spec) for more information about the different keys.


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
