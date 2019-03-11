package amqp

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/uniris/uniris-core/pkg/consensus"

	"github.com/streadway/amqp"
)

//ConsumeDiscoveryNotifications intercepts the RabbitMQ discovery notifications and store them into a specific database
//These data will be used by the consensus layer to identify peers and pools
func ConsumeDiscoveryNotifications(host string, user string, pwd string, port int, w consensus.NodeWriter) error {
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

	for err := range consumeQueues(ch, w) {
		fmt.Println(err)
	}

	return nil
}

type consumeFuncHandler func(w consensus.NodeWriter, data []byte) error

func consumeQueues(ch *amqp.Channel, w consensus.NodeWriter) (errChan chan error) {
	go func() {
		if err := consumeQueue(ch, queueNameDiscoveries, consumeDiscoveryHandler, w); err != nil {
			errChan <- err
		}
	}()

	go func() {
		if err := consumeQueue(ch, queueNameReachable, consumeReachableHandler, w); err != nil {
			errChan <- err
		}
	}()

	go func() {
		if err := consumeQueue(ch, queueNameUnreachable, consumeUnreachableHandler, w); err != nil {
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

func consumeDiscoveryHandler(w consensus.NodeWriter, data []byte) error {
	var n node
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}

	patch := consensus.ComputeGeoPatch(n.AppState.GeoPosition.Latitude, n.AppState.GeoPosition.Longitude)
	node := consensus.NewNode(net.ParseIP(n.Identity.IP),
		n.Identity.Port,
		n.Identity.PublicKey,
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
	if err := w.WriteDiscoveredNode(node); err != nil {
		return err
	}
	fmt.Printf("Discovered node stored: %s\n", n.Identity.PublicKey)
	return nil
}

func consumeReachableHandler(w consensus.NodeWriter, publicKey []byte) error {
	if err := w.WriteReachableNode(string(publicKey)); err != nil {
		return err
	}
	fmt.Printf("Node %s stored as reachable\n", publicKey)
	return nil
}

func consumeUnreachableHandler(w consensus.NodeWriter, publicKey []byte) error {
	if err := w.WriteUnreachableNode(string(publicKey)); err != nil {
		return err
	}
	fmt.Printf("Node %s stored as unreachable\n", publicKey)
	return nil
}
