# [WIP] Go F.S. C.R.U.D.

_Create_, _Read_, _Update_ and _Delete_ on _File System_ with _Go_.

## Documentation

To start the file server on [port 1234](http://localhost:1234), run the main file of this repo.

```sh
go run main.go
```

### Options

You can define the path to the store directory with the `-store` flag.

```sh
go run main.go -store=./store
```

## Endpoints

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/8e05ec219633e401ff14)

### Create

```sh
curl --location --request POST 'localhost:1234/test.txt' --data-raw 'Booyaka'
```

### Read

```sh
curl --location --request GET 'localhost:1234/test.txt'
```

### Update

```sh
curl --location --request PUT 'localhost:1234/test.txt' --data-raw 'Boomshakalakasha'
```

### Delete

```sh
curl --location --request DELETE 'localhost:1234/test.txt'
```
