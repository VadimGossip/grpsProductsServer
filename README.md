# grpsProductsServer

### Tools:
- go 1.19
- MongoDB
- Protobuf

### How to use this
Run containers with MongoDB:

```cmd
docker run --rm -d --name products-mongo -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=admin -p 27017:27017 mongo:latest
```

Building and running the application:
```cmd
go build -o app cmd/main.go
./app
```

### Server methods

This project contains following grpc methods:

- Fetch(URL) - requests an external CSV file with a list of products at an external address.
  The CSV file should look like PRODUCT NAME;PRICE. The last price of each product is stored in the database with the date of the request. The number of product price changes is also saved.
- List(paging params, sorting params) - gets a page-by-page list of products with their
  prices, the number of price changes and the dates of their last update.
  There is a possibility of sorting.

### Example:
 Example grps client for this server: https://github.com/VadimGossip/grpsProductsClient