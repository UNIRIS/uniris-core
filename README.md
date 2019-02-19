# uniris-core

Welcome to the UNIRIS Node source repository ! This software enables you to build the first transaction chain and next generation of blockchain focused on scalability and human oriented.

UNIRIS features:

- Fast transaction mining
- Lower energy consomption than other blockchain
- Designed with a high level of security
- Account management on the blockchain  itself (no need to keep your keys with you ever)
- Smart contract platform powered by a build-in interpreter
- Strong scalability using Artificial Intelligence 
- Secure key exchange to ensure a authorization between nodes and clients at anytime

## Supported Operating Systems

UNIRIS currently supports the following systems:

- Linux systems
- Mac OS X

## Executables

UNIRIS project comes with several executables found in the `cmd` directory

| Command            | Description | 
| -------------------|-------------| 
| uniris-node        | Entry point to the UNIRIS network by running a node 
| uniris-interpreter | UNIRIS smart contract interpreter

## Architecture

UNIRIS node software uses a separation of concerns designed by running multiple processes and services
| Service name | Protocol | Public | Description |
| ------------ | ---------| -------| ----------- |
| API          | HTTP     | Yes    | Gateway providing endpoints for clients
| Discovery    | GRPC     | Yes    | Peer to peer discovery using a built-in Gossip implementation
| Transaction  | GRPC     | Yes    | Transaction storage/queries, mining (proof of work, validation confirmations) and lock 

## Running UNIRIS

### API Access

The OpenAPI is available on the port `80`

### Configuration

UNIRIS executable support few command line flags in its scope and also configurable with a YAML file

Supported flags:

| Name | Description | Value |
| ---- | ----------- | ---- |
| conf | Configuration file | File location. default: `cmd/uniris-node/conf.yaml` |
| private-key|  private key | Hexadecimal ECDSA private key
| public-key |  public key | Hexadecimal ECDSA public key
| network-type | Type of network | `public` (default) or `private` |
| network-interface | Name of the network interface when network type is `private` | network interface name
| discovery-db-type|Discovery database instance type | `mem` (default), `redis` |
| discovery-db-host| Discovery database instance hostname | default: `localhost` |
| discovery-db-port | Discovery database instance port | default: `6379` (Redis) |
| discovery-db-password | Discovery database instance password | default: `''`
| bus-type| Bus messenging instance type | `mem` (default), `amqp` |
| bus-host | Bus messenging instance host | default: `localhost`|
| discovery-notif-user | Bus messenging instance user | default: `uniris` (AMQP) |
| bus-password | Bus messenging instance password | default: `uniris` (AMQP) |
| grpc-port | GRPC port | default: `5000`|

## Contribution

Thank you for considering to help out with the source code.
We welcome contributions from anyone and are grateful for even the smallest of improvement.

If you want to contribute to UNIRIS, we are using GitFlow approach. So please to fork, fix, commit and send us a pull-request for the maintainers to review and merge into the main code base. 

Please make sure your contributions adhere to our coding guidelines:
- Code must adhere to the official Go guidelines (https://github.com/golang/go/wiki/CodeReviewComments, https://golang.org/doc/effective_go.html)
- Pull request need to be based on the `develop` branch

## Licence

TBD