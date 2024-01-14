# Swagger Docs

Swagger docs are generated through [swaggo/swag](https://github.com/swaggo/swag) from Go annotations in the Godoc comments.

## Install `swag`

```shell
go install github.com/swaggo/swag/cmd/swag@latest
```

## Generate docs

```shell
swag init -g api.go -o http/docs/
```
