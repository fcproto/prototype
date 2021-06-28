# Fog Computing Prototyping Assignment
[![CI](https://github.com/fcproto/prototype/workflows/CI/badge.svg?branch=master)](https://github.com/fcproto/prototype/actions?query=workflow%3ACI+branch%3Amaster)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/fcproto/prototype)](https://pkg.go.dev/github.com/fcproto/prototype)

- [Introduction](#introduction)
  * [Demo Video](#demo-video)
- [Components](#components)
  * [Client](#client)
  * [Server](#server)
- [Build, Run And Deploy Server And Client](#build-run-and-deploy-server-and-client)
  * [Server](#server-1)
    + [Locally](#locally)
    + [Remote Build And Deployment With Cloud Run](#remote-build-and-deployment-with-cloud-run)
  * [Client](#client-1)

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
Every ten seconds the gateway sends all recent data to the server.
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

[Data format](https://pkg.go.dev/github.com/fcproto/prototype@v0.0.0-20210627163231-16f2b268c81c/pkg/api#SensorData):
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


## Build, Run And Deploy Server And Client

The backend uses Google Cloud Firestore, which must be set up:

- Create [new project](https://cloud.google.com/resource-manager/docs/creating-managing-projects)
- On the project dashboard open the navigation menu on the left and navigate to Firestore
- Select Cloud Firestore in native mode
- Select your preferred region and create the database

### Server
#### Locally

When running the backend locally it needs a credential file to access Firestore:

- On the project dashboard open the navigation menu on the left and go to _IAM & Admin > Service Accounts_
- Select the default account
- Go to the Keys tab and click _Create new key > JSON_ to download the credential file
- Rename the downloaded credential file to `fcproto-credentials.json` and move it to root of the project

Run the backend locally:

- Install [air](https://github.com/cosmtrek/air) on your system (provides auto reload on code changes and some other nice features)
- Run `air`

#### Remote Build And Deployment With Cloud Run

- Setup the [gcloud cli](https://cloud.google.com/sdk/docs/quickstart) for this project: `gcloud init`
- Run cloud build: `gcloud builds submit --tag gcr.io/<project-id>/<service-name>`
   - When it asks to activate cloud build confirm with `y`
- Deploy: `gcloud run deploy server --image gcr.io/<project-id>/<service-name>`
   - When it asks to allow unauthenticated invocations select yes
- The returned service url can be now used with the client

### Client

- Build: `go build -v -o ./bin/edge-service ./cmd/edge-service`
- Run: `CLOUD_ENDPOINT=<service-url> ./bin/edge-service`

