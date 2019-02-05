# uniris-core

Welcome tp the UNIRIS miner source repository ! This software enables you to build the first transaction chain and next generation of blockchain focused on scalability and human oriented.

UNIRIS features:

- Fast transaction mining
- Lower energy consomption than other blockchain
- Designed with a high level of security
- Account management on the blockchain  itself (no need to keep your keys with you ever)
- Smart contract platform powered by a build-in interpreter
- Strong scalability using Artificial Intelligence 
- Secure key exchange to ensure a authorization between miners and clients at anytime

## Supported Operating Systems

UNIRIS currently supports the following systems:

- Linux systems
- Mac OS X

## Executables

UNIRIS project comes with several executables found in the `cmd` directory

| Command       | Description | 
| ------------- |-------------| 
| uniris       | Entry point to the UNIRIS network by running a blockchain miner 

## Architecture

UNIRIS miner software uses a separation of concerns designed by running multiple processes with a microservice philosophy.

| Service name  | Public | Description |
| ------------  | -------| ----------- |
| API           | Yes | Gateway providing endpoints for clients
| Discovery     | Yes | Peer to peer discovery using a built-in Gossip protocol
| Mining        | Yes | Transaction mining (proof of work) and validation confirmations
| Lock          | Yes | Transaction lock and unlock to avoid double spending
| Chain         | Yes | Transaction storage and queries 
| Internal |  No | Forward transaction to other miners (mining, pool requesting) and retrieve shared data between miners

## Running UNIRIS

### Configuration

UNIRIS executable support few command line flags in its scope also configurable with a YAML file

Supported flags:

| Name | Description | Value |
| ---- | ----------- | ---- |
| conf | Configuration file | File location. default: `cmd/uniris/conf.yaml` |
| private-key| Miner private key | Hexadecimal ECDSA private key
| public-key | Miner public key | Hexadecimal ECDSA public key
| network-type | Type of network | `public` (default) or `private` |
|network-interface | Name of the network interface when network type is `private` | network interface name
| discovery-db-type|Discovery database instance type | `mem` (default), `redis` |
| discovery-db-host| Discovery database instance hostname | default: `localhost` |
| discovery-db-port | Discovery database instance port | default: `6379` (Redis) |
| discovery-db-password | Discovery database instance password | default: `''`
| bus-type| Bus messenging instance type | `mem` (default), `amqp` |
| bus-host | Bus messenging instance host | default: `localhost`|
| discovery-notif-user | Bus messenging instance user | default: `guest` (AMQP) |
| bus-password | Bus messenging instance password | default: `guest` (AMQP) |
| external-grpc-port | External GRPC port | default: `5000`|
| internal-grpc-port | Internal GRPC port  | default: `3009`|
| http-port | HTTP port the API | default: `8080`| 

Supported YAML configuration

```yaml
# NETWORK CONFIGURATION
network-type: private
network-interface: lo0

# MINER IDENTIFICATION
public-key: 3059301306072a8648ce3d020106082a8648ce3d0301070342000459c8b568df66798d7f876d94fb0afc516502893d996610632c40f70b830aebf39e0cbee311af4450ec56859d2b8f59ec09a44c7e303d030899aee551de61af2e
private-key: 307702010104201432f00062c0229d19e24c070a713d900da4883788ce3f8bd3fede4c10a36a79a00a06082a8648ce3d030107a1440342000459c8b568df66798d7f876d94fb0afc516502893d996610632c40f70b830aebf39e0cbee311af4450ec56859d2b8f59ec09a44c7e303d030899aee551de61af2e

#SERVICES CONFIGURATION
http-port: 8080
external-grpc-port: 5000
internal-grpc-port: 3009
discovery-seeds: 127.0.0.1:5000:3059301306072a8648ce3d020106082a8648ce3d0301070342000459c8b568df66798d7f876d94fb0afc516502893d996610632c40f70b830aebf39e0cbee311af4450ec56859d2b8f59ec09a44c7e303d030899aee551de61af2e
discovery-db-type: mem
discovery-db-host: localhost
discovery-db-port: 6379
discovery-db-password: ""
bus-type: mem
bus-host: localhost
bus-port: 5672
bus-user: guest
bus-password: guest
```

## Contribution

Thank you for considering to help out with the source code.
We welcome contributions from anyone and are grateful for even the smallest of improvement.

If you want to contribute to UNIRIS, we are using GitFlow approach. So please to fork, fix, commit and send us a pull-request for the maintainers to review and merge into the main code base. 

Please make sure your contributions adhere to our coding guidelines:
- Code must adhere to the official Go guidelines (https://github.com/golang/go/wiki/CodeReviewComments, https://golang.org/doc/effective_go.html)
- Pull request need to be based on the `develop` branch

## Licence

TBD