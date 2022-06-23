## Points to Highlight
- Usage of [Gorm](https://pkg.go.dev/gorm.io/gorm@v1.23.6) as the ORM.

## Design Decisions & Project Folder Structure
- Store model ORM inside the `objects` folder.
- Store orm prosgress service inside the `store` folder.
- Store API handler  inside the `handlers` folder.
- Store logger inside the `util` folder. This folder can be extended for other utility file
- Docker and docker compose are at the project root level
- Store main application code at project root level
```
.
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
**Object: Stock**
```go
// Stock object for the API
type Stock struct {
	// Identifier
	ID string `gorm:"primary_key" json:"id,omitempty"`

	// General details
	Name  string  `json:"name,omitempty"`
	Price float64 `json:"price,omitempty"`

	Availability int       `json:"availability,omitempty"`
	IsActive     bool      `json:"is_active,omitempty"`
	CreatedOn    time.Time `json:"created_on,omitempty"`
	UpdatedOn    time.Time `json:"updated_on,omitempty"`
}
```

## Run locally
### without docker compose
- Start postgres
- change file main.go to change the connection to DB or Prepare environment, change the value DB_CONN at .env file base on your DB  
- Build and run

```bash
$ export GO111MODULE=on
$ export GOFLAGS=-mod=vendor
$ go mod download
$ go run .
```
### with docker compose
```bash
$ docker-compose up --build
``` 

### Rest api

**Create A Stock**
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

**Get Stock**
```http request
GET http://localhost:8080/api/v1/stock/1655536052-0638474600-5197384620
Accept: application/json
###
```

**Update Stock's general details**
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

**List at max 10 records**
```http request
GET http://localhost:8080/api/v1/stocks?limit=10
Accept: application/json
###
```
## License
 
