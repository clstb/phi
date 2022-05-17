### TinkGW microservice
- Compatability layer for Tink API

### Core microservice
- Registration
- Login
- REST API for UI
- Coordination of users actions

### UI microservice
- Nuff said

### User authentication is handled by [ORY](https://console.ory.sh/)

### Ledger microservice
- Fills FS with users bean account data
- Serves bean account file to Fava

### FAVA microservice
- UI to visualize data from bean account file
- hosted in separate repo [Fava](https://github.com/Goofy-Goof/fava)

### Compile stubs
    make clean
    make proto

### Run UI
    cd ui
    pnpm install
    export NODE_OPTIONS=--openssl-legacy-provider
    pnpm start

### Do not work, unless TINK admin account used
    rpc ProvisionTinkUser
    rpc GetProviders

### Deploy to local k8s cluster
    make proto
    skaffold run --tail
- UI is available on [http://localhost:30002/](http://localhost:30002/)
- TinkGW requires TINK_CLIENT_ID and TINK_CLIENT_SECRET
#### Define k8s secret similar to
    apiVersion: v1
      kind: Secret
      metadata:
      name: tinkgw-secret
      data:
      client-id: xxx
      client-secret: xxx
- client-id and client-secret: must be base64 encoded
- apply with `kubectl apply -f k8s/secret.yaml`
