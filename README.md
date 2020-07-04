# Go F.S. C.R.U.D.

_Create_, _Read_, _Update_ and _Delete_ on _File System_ with _Go_.

## Documentation

To start the file server on [localhost](http://localhost:1234), run the main file of this repo.

```sh
go run main.go
```

### Options

```sh
go run main.go -port=:1234 -store=./store
```

#### `-port`

Port to serve the API on. _(Default: `:1234`)_

#### `-store`

Path to the file storage directory. _(Default: `./store`)_

## Endpoints

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/8e05ec219633e401ff14)

### Create

```sh
curl --request POST 'localhost:1234/test.txt' --data-raw 'Booyaka'
```

### Read

```sh
curl --request GET 'localhost:1234/test.txt'
```

### Update

```sh
curl --request PUT 'localhost:1234/test.txt' --data-raw 'Boomshakalakasha'
```

### Delete

```sh
curl --request DELETE 'localhost:1234/test.txt'
```
