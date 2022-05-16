### TinkGW microservice
- Compatability layer for Tink API

### Core microservice
- Registration
- Login
- REST API for UI

### UI microservice
- Nuff said

### User authentication is handled by [ORY](https://console.ory.sh/)
- Proxy to actual ORY cloud project

### Ledger microservice
- Provisions FS for user
- Fills FS with users bean account data

### FAVA microservice
- UI to visualize bean account data
- hosted in separate repo [FAVA](https://github.com/Goofy-Goof/fava)

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

