### TinkGW microservice
- is a compatability layer to Tink API

### Core microservice
- handles user creation, login and main features presented in UI

### UI microservice
- nuff said

### oathkeeper microservice
- proxy to actual ory cloud project

### Compile stubs
    make proto

### Run oauthkeeper
    cd oathkeeper
    docker build --tag  phi-oathkeeper --no-cache .
    docker run --name phi-oathkeeper -p 4455:4455 -p 4456:4456 --env LOG_LEAK_SENSITIVE_VALUES=true phi-oathkeeper


