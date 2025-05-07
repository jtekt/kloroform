# Kloroform

A CLI tool written in Go that allows to scale down all or up all deployments and statefulsets in a Kubernetes cluster

This code is based on [this sample](https://github.com/kubernetes/client-go/tree/master/examples/out-of-cluster-client-configuration)

## Usage

Kloroform can be used to either scale deployments and statefulsets down or back up

### Scaling down

To scale down all deployments and statefulsets in a cluster:

```
./kloroform
```

For specific namespaces:

```
./kloroform -namespaces=mynamespace,myothernamespace
```

To ignore namespaces:

```
./kloroform -exceptions=my-namespace,my-other-namespace
```

Specifying kubeconfig path:

```
./kloroform -kubeconfig=/home/myuser/.kube/config
```

### Scaling back up

To scale back up all deployments and statefulsets in a cluster:

```
./kloroform -wake
```

Specifying namespaces and kubeconfig follows the same pattern as for scaling down

## Development

To build executable:

```
go build
```
