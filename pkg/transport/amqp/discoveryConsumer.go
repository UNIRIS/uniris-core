package amqp

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/uniris/uniris-core/pkg/consensus"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"

	"github.com/streadway/amqp"
)

//ConsumeDiscoveryNotifications intercepts the RabbitMQ discovery notifications and store them into a specific database
//These data will be used by the consensus layer to identify peers and pools
func ConsumeDiscoveryNotifications(host string, user string, pwd string, port int, nodeWriter consensus.NodeWriter, sharedKeyWriter shared.KeyWriter) error {
	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%d", user, pwd, host, port)

	conn, err := amqp.Dial(amqpURI)
	defer conn.Close()
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	defer ch.Close()
	if err != nil {
		return err
	}

	for err := range consumeQueues(ch, nodeWriter, sharedKeyWriter) {
		fmt.Println(err)
	}

	return nil
}

type consumeFuncHandler func(w consensus.NodeWriter, data []byte) error

func consumeQueues(ch *amqp.Channel, nodeWriter consensus.NodeWriter, sharedKeyWriter shared.KeyWriter) (errChan chan error) {
	go func() {

		//TODO: remove sharedKeyWriter when the authorized key handling will be implemented
		if err := consumeQueue(ch, queueNameDiscoveries, func(w consensus.NodeWriter, data []byte) error {
			return consumeDiscoveryHandler(nodeWriter, sharedKeyWriter, data)
		}, nodeWriter); err != nil {
			errChan <- err
		}
	}()

	go func() {
		if err := consumeQueue(ch, queueNameReachable, consumeReachableHandler, nodeWriter); err != nil {
			errChan <- err
		}
	}()

	go func() {
		if err := consumeQueue(ch, queueNameUnreachable, consumeUnreachableHandler, nodeWriter); err != nil {
			errChan <- err
		}
	}()
	return
}

func consumeQueue(ch *amqp.Channel, queueName string, h consumeFuncHandler, w consensus.NodeWriter) error {
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			h(w, msg.Body)
		}
	}()

	<-forever

	return nil
}

func consumeDiscoveryHandler(nodeWriter consensus.NodeWriter, sharedKeyWriter shared.KeyWriter, data []byte) error {
	var n node
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}

	patch := consensus.ComputeGeoPatch(n.AppState.GeoPosition.Latitude, n.AppState.GeoPosition.Longitude)

	publicKey, err := crypto.ParsePublicKey([]byte(n.Identity.PublicKey))
	if err != nil {
		return err
	}

	node := consensus.NewNode(net.ParseIP(n.Identity.IP),
		n.Identity.Port,
		publicKey,
		consensus.NodeStatus(n.AppState.Status),
		n.AppState.CPULoad,
		n.AppState.FreeDiskSpace,
		n.AppState.Version,
		n.AppState.P2PFactor,
		n.AppState.ReachablePeersNumber,
		n.AppState.GeoPosition.Latitude,
		n.AppState.GeoPosition.Longitude,
		patch,
		true)
	if err := nodeWriter.WriteDiscoveredNode(node); err != nil {
		return err
	}
	fmt.Printf("Discovered node stored: %s\n", n.Identity.PublicKey)

	//TODO: remove sharedKeyWriter when the authorized key handling will be implemented
	return sharedKeyWriter.WriteAuthorizedNode(publicKey)
}

func consumeReachableHandler(w consensus.NodeWriter, publicKey []byte) error {
	pub, err := crypto.ParsePublicKey(publicKey)
	if err != nil {
		return err
	}
	if err := w.WriteReachableNode(pub); err != nil {
		return err
	}
	fmt.Printf("Node %s stored as reachable\n", publicKey)
	return nil
}

func consumeUnreachableHandler(w consensus.NodeWriter, publicKey []byte) error {
	pub, err := crypto.ParsePublicKey(publicKey)
	if err != nil {
		return err
	}
	if err := w.WriteUnreachableNode(pub); err != nil {
		return err
	}
	fmt.Printf("Node %s stored as unreachable\n", publicKey)
	return nil
}
