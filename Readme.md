### TinkGW microservice
- Compatability layer for Tink API

### Core microservice
- Registration
- Login
- REST API for UI
- Coordination of users actions

### UI microservice
- Login page, token request page


### Ledger microservice
- Fills FS with users bean account data
- Serves bean account file to Fava

### FAVA microservice
- UI to visualize data from bean account file
- Forked from [github.com/beancount/fava](https://github.com/beancount/fava)

### Compile stubs
    make proto

Yes, ideally they must be generated during build, and not added to VCS,
however, I'm too lazy (and project is too messy) to set up this properly.

### Do not work, unless TINK admin account used
    rpc ProvisionTinkUser
    rpc GetProviders

### Deploy to local k8s cluster
    skaffold run --tail
- UI is available on [http://localhost:30002/](http://localhost:30002/)
- TinkGW requires TINK_CLIENT_ID and TINK_CLIENT_SECRET


#### Yes, this all is unnecessarily over-engineered, I don't care.

