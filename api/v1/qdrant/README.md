 # qdrant vector store

This is an experimental [qdrant](https://qdrant.tech/) support.

If you wish to use it, you have to specify the DSN in the following format:

Insecure connections:
```shell
qdrant://[qdrant_cloud_api_key]@host:port
```

Secure encrypted connections:
```shell
qdrants://[qdrant_cloud_api_key]@host:port
```

If you want to run qdrant store by yourself, you wo't need the qdrant cloud API key so you can specify the following DSN URL i.e. just omit the API key section:
```shell
qdrants://@host:port
```
