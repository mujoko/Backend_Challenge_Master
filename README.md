## Run locally

- Start postgres
- Prepare environment, change the value POSTGRES_URL at .env file base on your DB
- execute schema.sql to create schema

## Points to Highlight
- Usage of [Chi](https://github.com/go-chi/chi) as the Router.
- Usage of [Zerolog](https://github.com/rs/zerolog) as the Logger.

## Design Decisions & Project Folder Structure
- Store config related files inside the `config` folder.
- Store model inside the `model` folder.
- Store API handler  inside the `controller` folder.
- Store route inside the `router` folder.
- Store logger inside the `util` folder. This folder can be extended for other utility file
- Store handler logger inside the `requestlog` folder. 
- Store main application code at project root level

```
.
├── LICENSE
├── README.md
├── config
│   └── config.go
├── controller
│   └── stock.go
├── go.mod
├── go.sum
├── main.go
├── main_test.go
├── models
│   └── models.go
├── requestlog
│   ├── handler.go
│   └── log_entry.go
├── router
│   └── router.go
├── schema.sql
└── util
    └── logger
        └── logger.go
```



- Build and run:

```bash
$ export GO111MODULE=on
$ export GOFLAGS=-mod=vendor
$ go mod download
$ go run .
$ 
```
 
## Test
 

## License
 