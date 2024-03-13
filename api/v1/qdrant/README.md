 # qdrant vector store

Experimental support for [qdrant](https://qdrant.tech/) vector store.

If you wish to use it, you have to specify the DSN in the following format:

Insecure connections, [qdrant cloud](https://cloud.qdrant.io/login):
```shell
qdrant://[qdrant_cloud_api_key]@host:port
```

Secure encrypted connections, [qdrant cloud](https://cloud.qdrant.io/login):
```shell
qdrants://[qdrant_cloud_api_key]@host:port
```

If you want to run qdrant store by yourself, you don't need the [qdrant cloud](https://cloud.qdrant.io/login) API key so you can specify the following DSN URL i.e. just omit the API key section:
```shell
qdrants://@host:port
```
