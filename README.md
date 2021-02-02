Alerterator
===========

Alerterator is a Kubernetes operator for managing Prometheus Alertmanager alerts. With this resource you can easily get
notified (either via Slack or email) when somethings is happening with your app.

As alerts are namespace agnostic, you don't have to have different files for each namespace you are running (although we
don't recommend running you application in different namespaces). You can even make your own personal alert-resources
that only notifes you!

The documentation for how to start using alerts are over at https://doc.nais.io/observability/alerts.

## Local testing

```
kind create cluster --image kindest/node:v1.17.11
kubeclt apply -f ./config/local-test/
make run

# different terminal
kubectl apply -f ./config/samples/alerts.yaml
kubectl get alerts
```
