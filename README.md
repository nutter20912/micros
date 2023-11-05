# Practice micro service
## Packages
#### config
- [spf13/viper](https://github.com/spf13/viper)
#### service
- [go-micro/go-micro](https://github.com/go-micro/go-micro)
- [go-micro/plugins](https://github.com/go-micro/plugins)
- [gorilla/websocket](github.com/gorilla/websocket)
#### database
- [go-gorm/gorm](https://github.com/go-gorm/gorm)
- [ariga/atlas-provider-gorm](https://github.com/ariga/atlas-provider-gorm)
- [golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- [mongodb/mongo-go-driver](https://github.com/mongodb/mongo-go-driver)
- [redis/go-redis](https://github.com/redis/go-redis)
#### message queue
- [ats-io/nats.go](github.com/nats-io/nats.go)
#### debug tool
- [go-delve/delve](https://github.com/go-delve/delve)
---
## Service
- user
- board
- market
- order
- wallet

---
## Setup
### Config
- generate local env file
    ```=bash
    cp .env.example .env
    ```
- configuration
    ```
    ./config
    ├── board.yaml
    ├── config.go
    ├── market.yaml
    ├── order.yaml
    ├── user.yaml
    └── wallet.yaml
    ```
---

## Scripts
read Makefile

