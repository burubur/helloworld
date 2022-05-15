# Helloworld

![Last Binary Build](https://github.com/burubur/helloworld/workflows/Last%20Binary%20Build/badge.svg)

A containerised **HTTP** based **Helloworld** microservice.

## Prerequisites

- [Docker](https://www.docker.com/)
- [Golang](https://golang.org/)
- **Unused port on 8080**

## Installation

```shell
make build
```

## How to Run The Service

```shell
make run
```

## How to Test The Running Service

```shell
make ping
```

## How to Login to The Running Container

> The command below require the `make run` to be performed successfully

```shell
make shell
```
