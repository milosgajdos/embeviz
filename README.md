# embeviz

[![Build Status](https://github.com/milosgajdos/embeviz/workflows/CI/badge.svg)](https://github.com/milosgajdos/embeviz/actions?query=workflow%3ACI)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/milosgajdos/embeviz)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache--2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A simple app that helps you visualize vector embeddings.

**THIS PROJECT IS WILDLY EXPERIMENTAL! USE AT YOUR OWN RISK OF SANITY!**

The app consists of two components:
* an API (written in Go, using [gofiber framework](https://docs.gofiber.io/))
* an SPA (written in JavaScript using [Reactjs](https://react.dev/) and [React Router](https://reactrouter.com/en/main))

The SPA is served as a static asset on `/ui` URL path when you start the app.

This project leverages the [go-embeddings](https://github.com/milosgajdos/go-embeddings) Go module for fetching embeddings from various API providers like OpenAI, etc.

As a result of this you must supply specific environment variables that are used to initialized the API clients for fetching the embeddings.

**NOTE:** the API stores the embeddings in an in-memory "DB" (it's a major Go maps hack)

# Build

Build the Go binary:
```shell
go get ./... && go build
```

SPA:
```shell
cd ui && npm install && npm run build
```

# Run

Once you've built the Go binary and bundled the webapp you can simply run the following command:
```shell
go run ./...
```

You should now be able to access the app on [http://localhost:5050/ui](http://localhost:5050/ui).

# TODO

* [ ] Clean up the code: both the Go and React
* [ ] Add support for a vector DB to store the embeddings in
* [ ] Embed the SPA and all its assets into the Go binary
