### TinkGW microservice
- Compatability layer for Tink API

### Core microservice
- Registration
- Login
- REST API for UI

### UI microservice
- Nuff said

### Oathkeeper microservice
- Proxy to actual ORY cloud project

### Ledger microservice
- Provisions FS for user
- Fills FS with users bean account data

### FAVA microservice
- UI to visualize bean account data
- hosted in separate repo [FAVA](https://github.com/Goofy-Goof/fava)

### Launch FAVA
    bean-example > .data/example.beancount
    fava .data/example.beancount

### Compile stubs
    make clean
    make proto

### Run oauthkeeper
    cd oathkeeper
    docker build --tag  phi-oathkeeper --no-cache .
    docker run --name phi-oathkeeper -p 4455:4455 -p 4456:4456 --env LOG_LEAK_SENSITIVE_VALUES=true phi-oathkeeper

### Run UI
    cd ui
    pnpm install
    export NODE_OPTIONS=--openssl-legacy-provider
    pnpm start

### Do not work, unless TINK admin account used
    rpc ProvisionTinkUser
    rpc GetProviders

