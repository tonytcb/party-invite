# Party Invite App

The goal of this project is design a Golang application serving an HTTP API to filter customers close to a specific location, when the input is a file containing customers and its geolocation.

More details in the [PROBLEM](PROBLEM.md) file.

## Design Solution

The application architecture follows the principles of the **Clean Architecture**, originally described by Robert C. Martin. The foundation of this kind of architecture is the dependency injection, producing systems that are independent of external agents, highly testable and easier to maintain.

You can read more about Clean Architecture [here](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).

## Tools

- [Golang 1.21](https://go.dev/)
- [Docker](https://www.docker.com/)
- [Docker-compose](https://docs.docker.com/compose/)

## API

### Filter Customers endpoint

- Method: `POST`
- Path: `/filter-customers`
- Params:
- - `file`: file containing a list of customers formatted as a JSON, each one in its own line. See an example [here](./Data/customers.txt).
- Response: A JSON containing the customers near to the specified location.

### Commands

- `make help` to see all commands;
- `make up` starts the app serving http api;
- `make test` to run all tests.

## Configurations

Instead of hardcode configurations, like `distance from base location` and `http port`, we are using a [.dot](./app.env) to define and easily change such parameters.

## TODO

- [ ] Implement a simple middleware
- [ ] Improve logger package using a third-party package, like logrus
- [ ] Add OpenTelemetry traces
- [ ] Implement integration tests
- [ ] Hot reload for docker development environment
- [ ] Docker file for production with multi stages
- [ ] Cache requests
- [ ] Idempotent API
- [ ] Decouple input and output from filter handler
- [ ] Add a concurrency mechanism on usecase layer to calculate distances 
- [ ] 