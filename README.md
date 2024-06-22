# Ethereum Blockchain Parser

## Pre-requisites
Code was tested on macos, but should be compatible with linux distributions as well. A dockerfile is provided as a convenience to ensure compatibility across all operating systems.

- go v1.22.2
- docker v26.1.1 (Optional)

## Quickstart
To run the base examples (make a single subscription and get notified whenever new transactions occur)
```sh
go run ./cmd/main.go
```

Alternatively, to run the application with docker
```sh
docker build . -t trustwallet
docker run --rm trustwallet
```

Entrypoint of the code is in `./cmd/main.go` file. User can easily modify the subscription parameters in the main function. The main function should give a good idea of how the notification service can utilize the parser interface. 

To build the binary
```sh
go build -o ./build/main ./cmd/main.go
```

## Project Structure
    .
    ├── cmd                  # Entrypoint to an application of the library
    ├── internal             
    │   ├── mocks            # Mocks used for unit testing
    │   ├── parser           # Pluggable Blockchain Parser
    │   ├── storage          # Pluggable Storage Engine
    │   └── subscriber       # Subscriber package is used to maintain subscriptions from clients and make the corresponding upstream subscriptions
    ├── pkg                  # Reusable packages that can be imported by other projects
    │   ├── ethereum         # Ethereum wrapper clients to interact with the blockchain
    │   ├── types            # Custom types used by blockchain parser
    └── README.md

## Tests
- External libraries were used for creating mocks
- Tests are only implemented for subscriber and storage package since it encapsulates most of the business logic required for notification service obtaining the right transactions, it should eventually be extended to the rest of the packages

To run tests
```sh
go test ./...
```

## Notes
- Subscriber needs to be started to begin indexing transactions for observers
- Clean up has not been implemented, in ideal world, we will have a channel to receive signals for stopping the subscriber service
- JSON RPC REST calls are used to retrieve all new transactions in the universe. In the ideal world, it is better to have a websockets subscription to the Ethereum client to reduce the total number of network calls required
- It was a deliberate choice to store all transaction data per address in memory, if transaction grows, application might run out of memory. For purpose of only notifying the delta transactions, this is not necessary as we can remove alerted transactions from memory. Rationale for doing so is to allow future extensibility when more than a single client uses the same parser, and multiple offsets are stored and replayed at any single point in time
- Parser implementation is currently kept in internal folder. For purposes of making this library extensible, it needs an additional wrapper in the package module to expose it with the right visibility to external packages making use of our library
- Interfaces is used for storage engine. The default implementation is to store all data in memory. This can easily be extended to persistent storage and dependency injecting in the main function
- There is a limit of 500k http requests to cloudflare ethereum rpc endpoints, this application only starts polling from the latest block at 15s interval, it is unlikely that running this application will hit the limit
- The RPC endpoint is polled at 15s interval, however, the same polling loop checks if latest block has gone ahead and directly makes the next poll if so, this allows us to always stay in sync with 12s mainnet block time, even though poll interval is slower than block progress speed