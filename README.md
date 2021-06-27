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

Our application for the prototyping assignment is placed in the domain of connected cars.
Cars collect data with onboard sensors that is relevant for other cars in the vicinity, like current speed, position and direction.
This data is shared among cars via the cloud to help for example autopilot systems.

The application is written in [Go](https://golang.org/) and the cloud part is hosted on [GCP](https://cloud.google.com/).

### Demo Video

[https://www.dropbox.com/s/3afh5fif4zs8aaz/prototype-demo.mov?dl=0](https://www.dropbox.com/s/3afh5fif4zs8aaz/prototype-demo.mov?dl=0)

## Components

The application consists of a local client and a server (cloud) part.

Each client represents a car, consisting of a gateway connected to the internet and various (simulated) sensors.
The gateway continuously collects, aggregates and locally stores the data from the sensors.
Every 10 seconds the gateway sends all recent data to the server.
In turn, the cars query the server for the most recent data of other cars in the vicinity.

![FCProto](https://user-images.githubusercontent.com/15909811/123540954-ecc34b80-d741-11eb-9419-3ae42e13ee89.png)

### Client

The client collects new data every second from the simulated sensors and transmits the measured values in batches to the cloud endpoint every ten seconds. It also supports higher sampling rates combined with different aggregations functions (mean average, minimum, and maximum). For example, the temperature is measured 20 times per second, and the aggregated average temperature is collected every second.
In our demo scenario the clients were running on different Raspberry Pi boards.

### Server

The server is a cloud-hosted HTTP backend hosted on [Cloud Run](https://cloud.google.com/run). The data is stored with [Cloud Firestore](https://firebase.google.com/docs/firestore).  
The public endpoint is available at https://server-ix6omulhiq-lm.a.run.app
It provides routes to store and retrieve sensor data formatted as JSON:

Routes:
- `POST /`  Store sensor data
- `GET /` Retrieve all data
- `GET /near/:client-id` Get most recent data from the nearest cars. The client sends its own id so the backend can query the data based on the latest info of the client
- `GET /status` Get information about the latest activity of the backend

Data format:
```json
{
    "clientId": "e65d62a0ee46894268cd0dd5",
    "timestamp": "2021-06-26T20:08:35.978612Z",
    "sensors": {
        "compass": {
            "rotation": 247.81560906021548
        },
        "gps": {
            "acceleration": -1.4474247587797677,
            "lat": 52.514659,
            "lon": 13.44130264174134,
            "speed": 10.643034842899365
        },
        "temperature/env": {
            "temperature": 25.784328651657326
        },
        "temperature/track": {
            "temperature": 31.17649297748596
        }
    }
}
```

## Implementation

### Build, run and deploy server

#### Locally

* Install [air](https://github.com/cosmtrek/air)
* Put GCP credential file `fcproto-credentials.json` in the root of the project (the credentials are needed for the Cloud Firestore access)
* Run `air`

#### Remote

- Setup [gcloud cli](https://cloud.google.com/sdk/docs/quickstart)
- Run cloud build: `gcloud builds submit --tag gcr.io/fcproto/server`
- Deploy: `gcloud run deploy server --image gcr.io/fcproto/server`

### Build & Run edge-service

```bash
go build -v -o ./bin/edge-service ./cmd/edge-service && ./bin/edge-service
```
