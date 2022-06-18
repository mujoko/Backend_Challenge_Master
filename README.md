## Points to Highlight
- Usage of [Gorm](https://pkg.go.dev/gorm.io/gorm@v1.23.6) as the ORM.

## Design Decisions & Project Folder Structure
- Store model ORM inside the `objects` folder.
- Store orm prosgress service inside the `store` folder.
- Store API handler  inside the `handlers` folder.
- Store logger inside the `util` folder. This folder can be extended for other utility file
- Docker and docker compose are at the project root level
- Store main application code at project root level

`.
├── Dockerfile
├── LICENSE
├── README.md
├── docker-compose.yml
├── errors
│   └── errors.go
├── go.mod
├── go.sum
├── handlers
│   ├── handlers.go
│   └── helpers.go
├── handlers_test.go
├── main.go
├── main_test.go
├── objects
│   ├── requests.go
│   └── stock.go
├── schema.sql
├── server.go
├── store
│   ├── postgres.go
│   └── store.go
└── util
    └── logger
        └── logger.go
```

### Rest api

## Run locally
### without docker compose
- Start postgres
- Prepare environment, change the value POSTGRES_URL at .env file base on your DB
- execute schema.sql to create schema
- Build and run

```bash
$ export GO111MODULE=on
$ export GOFLAGS=-mod=vendor
$ go mod download
$ go run .
```
### without docker compose
```bash
$ docker-compose up --build
``` 

### Rest api

**Create an Stock**
```http request
POST http://localhost:8080/api/v1/stock
Content-Type: application/json

{{
    "name":"Test",
    "price":1,
    "availability":3,
    "is_active":true
}
###
```

**Get event**
```http request
GET http://localhost:8080/api/v1/stock/1655536052-0638474600-5197384620
Accept: application/json
###
```

**Update event's general details**
```http request
http://localhost:8080/api/v1/stock/details{
    "id": "1655536052-0638474600-5197384620",
    "name":"Test",
    "price":1,
    "availability":1,
    "is_active":true
}

###
```

**List at max 42 events after the event: 20200828011748**
```http request
GET http://localhost:8080/api/v1/stock?limit=42&after=20200828011748
Accept: application/json
###
```
## License
 