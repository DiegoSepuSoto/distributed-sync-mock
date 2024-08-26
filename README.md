# Distributed Sync Mock

## Overview

This Golang API will send a range of IDs to use in any service that
requires coordination across instances. This service mimics the use
of a tool like [Zookeeper](https://zookeeper.apache.org/)

## Local workspace

First: clone the repository:

```bash
git clone https://github.com/DiegoSepuSoto/distributed-sync-mock && cd distributed-sync-mock
```

Then, download the dependencies

```bash
go mod download
```

Now you can run the application using the Makefile

```bash
make run
```

The available endpoint is the following:

```bash
curl --location --request POST 'localhost:8079/range'
```

which will create and send the range of IDs

Also, you can access:

- Prometheus metrics at: **localhost:8080/metrics**
- Swagger documentation at: **localhost:8080/swagger/index.html**

### Tech Stack

- Golang library - Echo framework for http server
- Golang library - Logrus for application logs
- Docker
