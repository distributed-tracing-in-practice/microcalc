# Calculator Microservice

This is a calculator application where individual operations are implemented as microservices in order to demonstrate tracing concepts.

## Prerequisites

* Go (tested on 1.13)
* Node.JS

## Installation and Usage

Right now, each service runs as a separate process on your machine, so you'll need to have several console windows open.

**Client**

1. Switch to the `web` directory.
2. Run `npm install`
3. Run `npm start`

**Server**

1. From the root directory, run `go run cmd/<service>/main.go`, where `<service>` is the name of the service you wish to run (`api`, `add`, etc.)

## Notes

By default, all of the trace data for each service will be output to standard output. In your browser, you'll need to open the JavaScript console to see it; For the Go applications, it'll be in your console.