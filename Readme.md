### Configure skaffold context
    export KUBE_CONTEXT=$(kubectl config current-context)
### Deploy auth, thinkgw services
    skaffold run --port-forward --tail
