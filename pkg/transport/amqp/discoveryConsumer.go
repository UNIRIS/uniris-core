package amqp

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/uniris/uniris-core/pkg/logging"
	"net"

	"github.com/uniris/uniris-core/pkg/consensus"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"

	"github.com/streadway/amqp"
)

//ConsumeDiscoveryNotifications intercepts the RabbitMQ discovery notifications and store them into a specific database
//These data will be used by the consensus layer to identify peers and pools
func ConsumeDiscoveryNotifications(host string, user string, pwd string, port int, nodeWriter consensus.NodeWriter, sharedKeyWriter shared.KeyWriter, l logging.Logger) error {
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

	for err := range consumeQueues(ch, nodeWriter, sharedKeyWriter, l) {
		l.Error(err.Error())
	}

	return nil
}

type consumeFuncHandler func(w consensus.NodeWriter, data []byte, l logging.Logger) error

func consumeQueues(ch *amqp.Channel, nodeWriter consensus.NodeWriter, sharedKeyWriter shared.KeyWriter, l logging.Logger) (errChan chan error) {
	go func() {

		//TODO: remove sharedKeyWriter when the authorized key handling will be implemented
		handler := func(w consensus.NodeWriter, data []byte, l logging.Logger) error {
			return consumeDiscoveryHandler(nodeWriter, sharedKeyWriter, data, l)
		}
		for err := range consumeQueue(ch, queueNameDiscoveries, handler, nodeWriter, l) {
			errChan <- err
		}
	}()

	go func() {
		for err := range consumeQueue(ch, queueNameReachable, consumeReachableHandler, nodeWriter, l) {
			errChan <- err
		}
	}()

	go func() {
		for err := range consumeQueue(ch, queueNameUnreachable, consumeUnreachableHandler, nodeWriter, l) {
			errChan <- err
		}
	}()
	return
}

func consumeQueue(ch *amqp.Channel, queueName string, h consumeFuncHandler, w consensus.NodeWriter, l logging.Logger) (errChan chan error) {
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		errChan <- err
		return
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
		errChan <- err
		return
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			if err := h(w, msg.Body, l); err != nil {
				errChan <- err
			}
		}
	}()

	<-forever
	return
}

func consumeDiscoveryHandler(nodeWriter consensus.NodeWriter, sharedKeyWriter shared.KeyWriter, data []byte, l logging.Logger) error {
	var n node
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}

	patch := consensus.ComputeGeoPatch(n.AppState.GeoPosition.Latitude, n.AppState.GeoPosition.Longitude)

	pubBytes, err := hex.DecodeString(string(n.Identity.PublicKey))
	if err != nil {
		return err
	}
	publicKey, err := crypto.ParsePublicKey(pubBytes)
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
	l.Info("Discovered node stored: " + string(n.Identity.PublicKey))

	//TODO: remove sharedKeyWriter when the authorized key handling will be implemented
	return sharedKeyWriter.WriteAuthorizedNode(publicKey)
}

func consumeReachableHandler(w consensus.NodeWriter, publicKey []byte, l logging.Logger) error {
	pub, err := crypto.ParsePublicKey(publicKey)
	if err != nil {
		return err
	}
	if err := w.WriteReachableNode(pub); err != nil {
		return err
	}
	l.Info("Node " + string(publicKey) + " stored as reachable")
	return nil
}

func consumeUnreachableHandler(w consensus.NodeWriter, publicKey []byte, l logging.Logger) error {
	pub, err := crypto.ParsePublicKey(publicKey)
	if err != nil {
		return err
	}
	if err := w.WriteUnreachableNode(pub); err != nil {
		return err
	}
	l.Info("Node " + string(publicKey) + " stored as unreachable")
	return nil
}
