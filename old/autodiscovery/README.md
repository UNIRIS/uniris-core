# Autodiscovery
Micro-service which manage autodiscovery system using a Gossip implementation.

## Gossip implementation

To gossip with another peer and to perform an autodiscovery exchange, the service use multiple messages to send and receive peer's informations:

- Peer 1 get predefined seed list for the first startup and load in its peer's database
- Peer 1 randomize a peer through its peer's database
- Peer 1 starts to gossip by sending all peers informations that he know (IP, Heartbeat) -> this step is called an SYN request
- Peer X receives the SYN request and returns the list of peers and theire details (IP, Heartbeat, AppState) than Peer1 does not know and ask the same from Peer 1 -> this step is called an ACK request.
- Peer 1 receives the ACK request and send to Peer X the requested informations (IP, Heartbeat, AppState) -> this step is called ACK2 request.

Each cycle every peer repeat this process with a defined number of peers.
