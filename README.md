# Prototype

## Build & Run sensor-test

```bash
go build ./cmd/sensor-test/ && ./sensor-test
```

## Build, run and deploy server

### Locally
```bash
go build -v -o ./bin/server ./cmd/server && ./bin/server
```

### Remote

- Setup gcloud cli

```bash
gcloud builds submit --tag gcr.io/fcproto/server
gcloud run deploy server --image gcr.io/fcproto/server
```
