# Autodiscovery
Micro-service which manage autodiscovery system using a Gossip implementation.

## Architecture

This service follow the clean architecture pattern and its composed by the layers belows:

```
---------------------------------------------------------------------------------------------------------------------
--domain            # Peer entities and buiness related
--adapters          # Transformation layer for infrastructure and domain layer (i.e repository implementation, DTO adapters)
--usecases          # Handles business interactions (manipulates domain entities, define interfaces for infrastructure or adapter layer implementations)
--infrastructure    # Manages thirdparty integrations (i.e. HTTP server, database, message bus)
main.go             # MAIN entry point
---------------------------------------------------------------------------------------------------------------------
```

## Gossip implementation

To gossip with another peer and to perform an autodiscovery exchange, the service use multiple messages to send and receive peer's informations:

- Peer 1 get predefined seed list for the first startup and load in its peer's database
- Peer 1 randomize a peer through its peer's database
- Peer 1 starts to gossip by sending all peers informations that he know (IP, Heartbeat) -> this step is called an SYN request
- Peer X receives the SYN request and returns the list of peers and theire details (IP, Heartbeat, AppState) than Peer1 does not know and ask the same from Peer 1 -> this step is called an ACK request.
- Peer 1 receives the ACK request and send to Peer X the requested informations (IP, Heartbeat, AppState) -> this step is called ACK2 request.

Each cycle every peer repeat this process with a defined number of peers.

## Inter communication services

Once the peers received information from others peers (so after gossip rounds), it publish on the message bus queue the peers discovered or the updated informations.
Using that, other services which subscribe to the discovered peers queue receive the new discovered data and can process them.

For example:
- AI service which will compute to rank peers