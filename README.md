## Run locally

- Start postgres
- Prepare environment, fill DB parameters:

``` bash
$ source env-sample
```

- Build and run:

```bash
$ export GO111MODULE=on
$ export GOFLAGS=-mod=vendor
$ go mod download
$ go build -o go-mux-api.bin
$ ./go-mux-api.bin
```

Server is listening on localhost:8010

## Test

```bash
$ go test -v
=== RUN   TestEmptyTable
--- PASS: TestEmptyTable (0.00s)
=== RUN   TestGetNonExistentProduct
--- PASS: TestGetNonExistentProduct (0.00s)
=== RUN   TestCreateProduct
--- PASS: TestCreateProduct (0.00s)
=== RUN   TestGetProduct
--- PASS: TestGetProduct (0.00s)
=== RUN   TestUpdateProduct
--- PASS: TestUpdateProduct (0.01s)
=== RUN   TestDeleteProduct
--- PASS: TestDeleteProduct (0.01s)
PASS
ok      _/home/tom/r/go-mux-api 0.034s
```

## License

If I got Accepted to Work, then you can keep this
else the licence for this I put as Apache