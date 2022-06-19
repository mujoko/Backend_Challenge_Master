# create a stock
curl -X POST \
  http://localhost:8080/api/v1/stock \
  -H 'content-type: application/json' \
  -d '{
    "name":"Meat Ball",
    "price":100,
    "availability":1000,
    "is_active":true
}';

# list one stocks
curl -X GET 'http://localhost:8080/api/v1/stocks?limit=1'

# or try getting one using an id (note: your id will be something different)
curl -X GET 'http://localhost:8080/api/v1/stock/1599425402-0970640120-4671038529'