# osbpsql

[![Go Report Card](https://goreportcard.com/badge/github.com/knqyf263/osbpsql)](https://goreportcard.com/report/github.com/knqyf263/osbpsql)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://github.com/knqyf263/osbpsql/blob/master/LICENSE)


An implementation of the Open Service Broker API for PostgreSQL.

Using this tool with CI, it makes it easier to create database for testing (create and drop in a short period).

CLOUD FOUNDRY and OPEN SERVICE BROKER are trademarks of the CloudFoundry.org Foundation in the United States and other countries.

## Supported Services
- Create database (database name is generated automatically)
- Create user (username and password is generated automatically)

## Installation and Usage

### Prerequisites

- [Kubernetes](https://kubernetes.io/) 1.7+ with RBAC enabled
- A working Helm installation
- [Service Catalog](https://github.com/kubernetes-incubator/service-catalog)
- [Helm](https://github.com/kubernetes/helm)
- [Optional] [`svcat`: Service Catalog CLI](https://github.com/kubernetes-incubator/service-catalog/tree/master/cmd/svcat)


### Install
Use Helm to install Open Service Broker for PostgreSQL onto your Kubernetes cluster.

Installation of this chart is simple. First, add the osbpsql charts repository to your local list:

```
$ helm repo add osbpsql https://knqyf263.github.io/osbpsql/
```

Next, install from the osbpsql repo:

```
$ export DB_HOST=[YOUR POSTGRESQL HOST]
$ export DB_PORT=:5432
$ export DB_USER=postgres
$ export DB_PASSWORD=postgres
$ export DB_DATABASE=postgres

$ helm install osbpsql/osbpsql --name osbpsql \
  --set db.host=$DB_HOST \
  --set db.port=$DB_PORT \
  --set db.user=$DB_USER \
  --set db.password=$DB_PASSWORD \
  --set db.database=$DB_DATABASE
```

### Uninstall
```
$ helm delete osbpsql --purge
```
The command removes all the Kubernetes components associated with the chart and deletes the release.


## Configuration

The following tables lists the configurable parameters of the Service
Broker for PostgreSQL chart and their default values.

| Parameter                   | Description | Default |
| --------------------------- | ----------- | ------- |
| `image.repository`          | Docker image location, _without_ the tag. | `"knqyf263/osbpsql"` |
| `image.tag`                 | Tag / version of the Docker image. | `"latest"` |
| `image.pullPolicy`          | `"IfNotPresent"`, `"Always"`, or `"Never"`; When launching a pod, this option indicates when to pull the OSBS Docker image. | `"IfNotPresent"` |
| `basicAuth.username`        | Specifies the basic auth username that clients (e.g. the Kubernetes Service Catalog) must use when connecting to OSBA. | `"username"`; __Do not use this default value in production!__ |
| `basicAuth.password`        | Specifies the basic auth password that clients (e.g. the Kubernetes Service Catalog) must use when connecting to OSBA. | `"password"`; __Do not use this default value in production!__ |
| `db.host`        | Database hostname | `"localhost"` |
| `db.port`        | Database port | `":5432"` |
| `db.username`        | Database username | `"postgres"` |
| `db.password`        | Database password | `"postgres"` |
| `db.database`        | Database name | `"postgres"` |


Specify a value for each option using the `--set <key>=<value>` switch on the
`helm install` command. That switch can be invoked multiple times to set
multiple options.

Alternatively, copy the charts default values to a file, edit the file to your
liking, and reference that file in your `helm install` command:

```console
$ helm inspect values knqyf263/osbpsql > my-values.yaml
$ vim my-values.yaml
$ helm install knqyf263/osbpsql --name osbs --values my-values.yaml
```

### Provisioning

With the Kubernetes Service Catalog and Open Service Broker for PostgreSQL both installed on your Kubernetes cluster,
try creating a ServiceInstance resource to see service provisioning in action.

For example, the following will create database on PostgreSQL:

```bash
$ vi database-service.yaml
```

```console
apiVersion: servicecatalog.k8s.io/v1beta1
kind: ServiceInstance
metadata:
  name: test-psql-database
  namespace: default
spec:
  clusterServiceClassExternalName: psql-database
  clusterServicePlanExternalName: standard
```

```bash
# Create database for test using service catalog (CREATE DATABASE dbname)
$ kubectl create -f database-service.yaml
```

After the ServiceInstance resource is submitted, you can view its status:

```bash
# using kubectl
$ kubectl get serviceinstance test-psql-database -o yaml 

# using svcat(Service Catalog CLI)
$ svcat describe instance test-psql-database
```

You'll see output that includes a status indicating that asynchronous provisioning is ongoing. Eventually,
that status will change to indicate that asynchronous provisioning is complete.

### Binding

Upon provision success, bind to the instance:

```bash
$ vi database-binding.yaml
```

```console
apiVersion: servicecatalog.k8s.io/v1beta1
kind: ServiceBinding
metadata:
  name: test-psql-binding
  namespace: default
spec:
  instanceRef:
    name: test-psql-database
  secretName: test-psql-secret
```

```
$ kubectl create -d database-binding.yaml
```

To check the status of the binding:

```bash
# using kubectl
$ kubectl get servicebinding test-psql-binding -o yaml

# using svcat(Service Catalog CLI)
$ svcat describe binding test-psql-binding
```

You'll see some output indicating that the binding was successful.
Once it is, a secret named test-psql-secret will be written that contains the database name in it.

You can observe that this secret exists and has been populated:

```bash
kubectl get secret test-psql-secret -o yaml
```

This secret can be used just as any other.

### Unbinding

To unbind:

```bash
$ kubectl delete -f database-binding.yaml
```

### Deprovisioning

To deprovision:

```bash
$ kubectl delete -f database-service.yaml
```

## Author

  * Teppei Fukuda