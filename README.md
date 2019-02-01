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

| Service name  | Description |
| ------------  | ----------- |
| API           | Provide HTTP/REST endpoints for clients and SDKs
| Discovery     | Peer to peer discovery using a built-in Gossip protocol
| Datamining    | Transaction mining and storage including sharding computation and replication

## Running UNIRIS

### Configuration

UNIRIS executable support few command line flags in its scope also configurable with a YAML file

Supported flags:

| Name | Description | Value |
| ---- | ----------- | ---- |
| conf | Load a configuration file | File location. default: `cmd/uniris/conf.yaml` |
| private-key| Define the miner private key | Hexadecimal ECDSA private key
| public-key | Define the miner public key | Hexadecimal ECDSA public key
| network-type | Define which can of network to use | `public` (default) or `private` |
|network-interface | Define the network interface to pick when the network type is `private` | network interface name
| discovery-port | Define the port of the discovery service | default: `4000`
| discovery-seeds | Define the peer list to seed the gossip discovery | example: `IP:PORT:PUBLICKEY;IP:PORT:PUBLICKEY` |
| discovery-db-type| Define the type of discovery database | `mem` (default), `redis` |
| discovery-db-host| Define the hostname of the discovery database | default: `localhost` |
| discovery-db-port | Define the port of the discovery database | default: `6379` (Redis) |
| discovery-db-password | Define the database password | default: `''`
| discovery-notif-type| Define the type of discovery notifier | `mem` (default), `amqp` |
| discovery-notif-host | Define the hostname discovery notifier | default: `localhost`|
| discovery-notif-user | Define the user of the discovery notifier | default: `guest` (RabbitMQ) |
| discovery-notif-password | Define the password of the discovery notifier | default: `guest` (AMQP) |
| datamining-port | Define the port to the datamining GPRC service | default: `5000`|
| datamining-internal-port | Define the internal port to the datamining service to make other service able to call it | default: `3009`|
| api-port | Define the API port | default: `8080`| 

Supported YAML configuration

```yaml
# NETWORK CONFIGURATION
network-type: private
network-interface: lo0

# MINER IDENTIFICATION
public-key: 3059301306072a8648ce3d020106082a8648ce3d0301070342000459c8b568df66798d7f876d94fb0afc516502893d996610632c40f70b830aebf39e0cbee311af4450ec56859d2b8f59ec09a44c7e303d030899aee551de61af2e
private-key: 307702010104201432f00062c0229d19e24c070a713d900da4883788ce3f8bd3fede4c10a36a79a00a06082a8648ce3d030107a1440342000459c8b568df66798d7f876d94fb0afc516502893d996610632c40f70b830aebf39e0cbee311af4450ec56859d2b8f59ec09a44c7e303d030899aee551de61af2e

#SERVICES CONFIGURATION
api-port: 8080
discovery-port: 4000
discovery-seeds: 127.0.0.1:4000:3059301306072a8648ce3d020106082a8648ce3d0301070342000459c8b568df66798d7f876d94fb0afc516502893d996610632c40f70b830aebf39e0cbee311af4450ec56859d2b8f59ec09a44c7e303d030899aee551de61af2e
discovery-db-type: mem
discovery-db-host: localhost
discovery-db-port: 6379
discovery-db-password: ""
discovery-notif-type: mem
discovery-notif-host: localhost
discovery-notif-port: 5672
discovery-notif-user: guest
discovery-notif-password: guest
datamining-port: 5000
datamining-internal-port: 3009
```

## Contribution

Thank you for considering to help out with the source code.
We welcome contributions from anyone and are grateful for even the smallest of improvement.

If you want to contribute to UNIRIS, we are using GitFlow approach. So please to fork, fix, commit and send us a pull-request for the maintainers to review and merge into the main code base. 

Please make sure your contributions adhere to our coding guidelines:
- Code must adhere to the official Go guidelines (https://github.com/golang/go/wiki/CodeReviewComments, https://golang.org/doc/effective_go.html)
- Code must use the `domain driven design`  and `hexagonal architecture` approaches (aka Clean architecture)
- Pull request need to be based on the `develop` branch

## Licence

TBD