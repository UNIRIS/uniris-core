# robot-miner
This repository concerns the UNIRIS software for miners

## Composition

It's composed for 5 micro services:
- API: Handle REST endpoints for clients and SDKs
- Autodiscovery: Gossip discovery
- Replication: Ledger storage and gossip ledgers
- Mining: Transaction mining and approvals
- AI: Artifical Inteligence

The services communicates with Advanced Message Queue Protocol as RabbitMQ

## Contribution

- We use the GitFlow methodology.
- When your work is finished please create a pull-request, otherwise your work will be not merged
