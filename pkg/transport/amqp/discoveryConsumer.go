package amqp

import (
	"fmt"

	"github.com/streadway/amqp"
)

//ConsumeDiscoveryNotifications intercepts the RabbitMQ discovery notifications and store them into a specific database
//These data will be used by the consensus layer to identify peers and pools
func ConsumeDiscoveryNotifications(host string, user string, pwd string, port int) error {
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

	for err := range consumeQueues(ch) {
		fmt.Println(err)
	}

	return nil
}

type consumeFuncHandler func(data []byte) error

func consumeQueues(ch *amqp.Channel) (errChan chan error) {
	go func() {
		if err := consumeQueue(ch, queueNameDiscoveries, consumeDiscoveryHandler); err != nil {
			errChan <- err
		}
	}()

	go func() {
		if err := consumeQueue(ch, queueNameReachable, consumeReachableHandler); err != nil {
			errChan <- err
		}
	}()

	go func() {
		if err := consumeQueue(ch, queueNameUnreachable, consumeUnreachableHandler); err != nil {
			errChan <- err
		}
	}()
	return
}

func consumeQueue(ch *amqp.Channel, queueName string, h consumeFuncHandler) error {
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
			h(msg.Body)
		}
	}()

	<-forever

	return nil
}

func consumeDiscoveryHandler(data []byte) error {
	p, err := unmarshalPeer(data)
	if err != nil {
		return err
	}
	fmt.Printf("Discovered peer: %s\n", p)
	return nil
}

func consumeReachableHandler(data []byte) error {
	p, err := unmarshalPeerIdentity(data)
	if err != nil {
		return err
	}
	fmt.Printf("Reachable peer: %s\n", p.Endpoint())
	return nil
}

func consumeUnreachableHandler(data []byte) error {
	p, err := unmarshalPeerIdentity(data)
	if err != nil {
		return err
	}
	fmt.Printf("Unreachable peer: %s\n", p.Endpoint())
	return nil
}
