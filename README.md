# embeviz

A simple app that helps you visualize data embeddings.

> [!WARNING]
> THIS PROJECT IS WILDLY EXPERIMENTAL AT THE MOMENT! USE AT YOUR OWN RISK!
> IF YOU LIKE CLEAN AND DRY CODE THIS ISN'T GONNA BE YOUR JAM!

[![Build Status](https://github.com/milosgajdos/embeviz/workflows/CI/badge.svg)](https://github.com/milosgajdos/embeviz/actions?query=workflow%3ACI)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/milosgajdos/embeviz)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache--2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

The app consists of two components:
* A JSON API (written in Go, using [gofiber framework](https://docs.gofiber.io/))
* An SPA (written in JavaScript using [ReactJS](https://react.dev/) framework with [React Router](https://reactrouter.com/en/main))

The SPA is accessible on `/ui` URL path when you run the app.

> [!IMPORTANT]
> Before you try accessing the SPA you must build it first. See the [README](./ui/README.md) for more details.

<p align="center">
  <img alt="Embeviz Home" src="./ui/public/home.png" width="45%">
&nbsp; &nbsp; &nbsp; &nbsp;
  <img alt="Embeviz Provider" src="./ui/public/provider.png" width="45%">
</p>

The API provides a swagger API endpoint on `/api/v1/docs` which serves the API documentation powering the SPA.

<p align="center">
  <img alt="Swagger endpoints" src="./ui/public/swagger_endpoints.png" width="45%">
&nbsp; &nbsp; &nbsp; &nbsp;
  <img alt="Swagger models" src="./ui/public/swagger_models.png" width="45%">
</p>

The app leverages the [go-embeddings](https://github.com/milosgajdos/go-embeddings) Go module for fetching embeddings from various API providers like [OpenAI](https://openai.com/), etc.

As a result of this you must supply specific environment variables when you run the app. The environment variables are used to initialize the API clients for fetching the embeddings. See the `README` of the `go-embeddings` module for more details.

> [!WARNING]
> By default the API stores the embeddings in an in-memory store (it's a major Go maps hack!)
> The only vector store currently supported is [qdrant](https://qdrant.tech/). See the [docs](./api/v1/qdrant).

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

> [!IMPORTANT]
> Before you run the app you need to make sure you have set some environment variables required by specific AI embeddings API providers. See the list below

The project relies on the [go-embeddings](https://github.com/milosgajdos/go-embeddings) Go module so we only support specific AI embeddings API providers:
* [OpenAI](https://openai.com/): `OPENAI_API_KEY`
* [Cohere](https://cohere.com/): `COHERE_API_KEY`
* [Google VertexAI](https://cloud.google.com/vertex-ai/docs/generative-ai/learn/overview): `VERTEXAI_TOKEN` (get it by running `gcloud auth print-access-token` once you've set up your GCP project and authenticated locally) and `GOOGLE_PROJECT_ID` (the ID of the GCP project)

> [!NOTE]
> If none of the above environment vars has been set, no AI embeddings provider is loaded and you won't be able to interact with the app.
> The project doesn't allow adding new embeddings providers at the moment.

Once you've built the Go binary and bundled the webapp you can simply run the following command:
```shell
OPENAI_API_KEY="sk-XXXX" COHERE_API_KEY="XXX" go run ./...
```

Alternatively you can also run the following command:
```shell
go run ./...
```

You should now be able to access the SPA on [http://localhost:5050/ui](http://localhost:5050/ui).

The API docs should be available on [http://localhost:5050/api/v1/docs](http://localhost:5050/api/v1/docs).

# TODO

* [ ] Clean up the code: both Go and React
* [ ] Embed the SPA into the Go binary

# Contributing

YES PLEASE!
