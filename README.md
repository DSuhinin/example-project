# example-test-task
Example Project

### Prerequisite
Tools that should be installed before run/lint/test:
- Makefile
- Docker
- Golang

### Repository Structure
- `./app` - contains application itself.
- `./build` - contains everything related to service build in Docker.
- `./core` - contains like a `core/common` functionality.
- `./doc` - API docs in `swagger` format. To run/see docs just open `./docs/index.html`.
- `./test` - contains everything related to `integration` tests like: fixtures, test helpers, test data, etc.

### Useful Commands
- `make lint` - to run linter over the code.
- `make service_start_demo` - to run service that will be pointed to the `live` `fixer` service.
- `make service_stop` - stop service.
- `make service_test` - to run `integration` tests over the service.

### How to play with the service

- after `make service_start_demo` service should be available under `http://127.0.0.1:8080` host.
- all endpoints could be found in `swagger` doc.
- firstly `key` should be created through the `POST /keys` endpoint. All `keys` endpoints protected by `Authorization: AdminKey %placehoder for admin key%` header. 
`AdminKey` value is by default equal to `supersecurekey` which could be configured in `./build/currency-exchange-service-demo.env` in `ADMIN_KEY` env variable.
- to call `GET /currencies/exchange`, `key` from the previous step should be provided in a header: `Authorization: Key %placehoder for key%`. 



