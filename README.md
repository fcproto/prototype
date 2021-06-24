# Fog Computing Prototyping Assignment

- [Introduction](#introduction)
- [Components](#components)
    * [Client](#client)
    * [Server](#server)
- [Implementation](#implementation)
    * [Build, run and deploy server](#build--run-and-deploy-server)
        + [Locally](#locally)
        + [Remote](#remote)
    * [Build & Run edge-service](#build---run-edge-service)

## Introduction

## Components

### Client

### Server

## Implementation

### Build, run and deploy server

#### Locally

* Install [air](https://github.com/cosmtrek/air)
* Put GCP credential file `fcproto-credentials.json` in the root of the project
* Run `air`

#### Remote

- Setup gcloud cli

```bash
gcloud builds submit --tag gcr.io/fcproto/server
gcloud run deploy server --image gcr.io/fcproto/server
```

### Build & Run edge-service

```bash
go build -v -o ./bin/edge-service ./cmd/edge-service && ./bin/edge-service
```
