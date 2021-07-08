# Phi

[![Go Reference](https://pkg.go.dev/badge/github.com/clstb/phi.svg)](https://pkg.go.dev/github.com/clstb/phi)
[![Go Report Card](https://goreportcard.com/badge/github.com/clstb/phi)](https://goreportcard.com/report/github.com/clstb/phi)
[![Build Status](https://cloud.drone.io/api/badges/clstb/phi/status.svg)](https://cloud.drone.io/clstb/phi)

---
Phi is an open source system for managing your finances.
Phi provides tools to efficiently ingest finance data and generate financial reports such as trial balance, income statement, balance sheet and journal.  

It is written in go, consisting of multiple containerized microservices. The services communicate with GRPC and use PostgresQL compatible databases for storage.
Currently Phi has following services:
| Service | Spec                        | Description                                                    |
|:------- |:--------------------------- |:-------------------------------------------------------------- |
| Core    | [Link](/proto/core.proto)   | Core functionality such as accounts and transactions           |
| Auth    | [Link](/proto/auth.proto)   | User management and authentication via JWT's                   |
| TinkGW  | [Link](/proto/tinkgw.proto) | Links Phi to bank accounts using Tinks PSD2 API                |

## Features
Besides functionalities of a double entry accounting system phi has following features:
* Bayesian classification of transactions
* Import transactions from your bank account programmatically using [Tink](https://tink.com/)
* Import transactions using csv
* Support for multiple users
* GRPC interface for all services

## Contribute

### Skaffold
Skaffold is a system to rapidly develop distributed systems on a local kubernetes cluster.
With it you can get Phi running on your local machine in a few commands.

1. Install [skaffold](https://skaffold.dev/docs/install/).
2. Install a local kubernetes distribution such as [minikube](https://minikube.sigs.k8s.io/docs/start/) or [kind](https://github.com/kubernetes-sigs/kind).
3. Deploy a PostgresQL compatible database to your local cluster. In this example we will use cockroachdb. Wait for the pods to be up and running.
```sh
kubectl create -f https://raw.githubusercontent.com/cockroachdb/cockroach/master/cloud/kubernetes/cockroachdb-statefulset.yaml
kubectl get pods
```
4. Initialize the cockroachdb cluster.
```sh
kubectl create -f https://raw.githubusercontent.com/cockroachdb/cockroach/master/cloud/kubernetes/cluster-init.yaml
```
5. Connect to the database.
```sh
kubectl run cockroachdb -it --image=cockroachdb/cockroach:v20.2.4 --rm --restart=Never -- sql --insecure --host=cockroachdb-public
```
6. Create needed databases and users.
```sql
CREATE DATABASE phi_core;
CREATE DATABASE phi_auth;
CREATE DATABASE phi_tinkgw;
CREATE USER phi_core;
CREATE USER phi_auth;
CREATE USER phi_tinkgw;
GRANT ALL ON DATABASE phi_core TO phi_core;
GRANT ALL ON DATABASE phi_auth TO phi_auth;
GRANT ALL ON DATABASE phi_tinkgw TO phi_tinkgw;
```
7. Deploy Phi.
```sh
skaffold dev --port-forward
```
