# Phi :warning: Alpha software. Do not use for production! :warning:

[![Go Reference](https://pkg.go.dev/badge/github.com/clstb/phi.svg)](https://pkg.go.dev/github.com/clstb/phi)
[![Go Report Card](https://goreportcard.com/badge/github.com/clstb/phi)](https://goreportcard.com/report/github.com/clstb/phi)

---
Phi is an open source system for managing your finances.
Phi provides tools to efficiently ingest finance data and generate financial reports such as trial balance, income statement, balance sheet and journal.  

It is written in go, consisting of multiple containerized microservices. The services communicate with GRPC and use PostgresQL compatible databases for storage.
Currently Phi has following services:
| Service | Spec                      | Description                                                    |
|:------- |:------------------------- |:-------------------------------------------------------------- |
| Core    | [Link](/proto/core.proto) | Core functionality such as Accounts, Transactions and Postings |
| Auth    | [Link](/proto/auth.proto) | User management and authentication via JWT's                   |

## Contribute

### Skaffold
Skaffold is a system to rapidly develop distributed systems on a local kubernetes cluster.
With it you can get Phi running on your local machine in a few commands.

1. Install [skaffold](https://skaffold.dev/docs/install/)
2. Install a local kubernetes distribution such as [minikube](https://minikube.sigs.k8s.io/docs/start/) or [kind](https://github.com/kubernetes-sigs/kind)
3. Deploy a PostgresQL compatible database to your local cluster
```sh
kubectl create -f https://raw.githubusercontent.com/cockroachdb/cockroach/master/cloud/kubernetes/cockroachdb-statefulset.yaml
kubectl create -f https://raw.githubusercontent.com/cockroachdb/cockroach/master/cloud/kubernetes/cluster-init.yaml
```
4. Connect to the database
```sh
kubectl run cockroachdb -it --image=cockroachdb/cockroach:v20.2.4 --rm --restart=Never -- sql --insecure --host=cockroachdb-public
```
5. Create needed databases and users
```sql
CREATE DATABASE phi_core;
CREATE DATABASE phi_auth;
CREATE USER phi_core;
CREATE USER phi_auth;
GRANT ALL ON DATABASE phi_core TO phi_core;
GRANT ALL ON DATABASE phi_auth TO phi_auth;
```
6. Deploy Phi
```sh
skaffold dev --port-forward
```
