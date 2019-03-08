package amqp

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"github.com/uniris/uniris-core/pkg/discovery"
)

//NewDiscoveryNotifier creates a discovery notifier using AMQP
func NewDiscoveryNotifier(host, user, pwd string, port int) discovery.Notifier {
	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%d", user, pwd, host, port)
	return notifier{amqpURI}
}

const (
	queueNameDiscoveries = "autodiscovery_discoveries"
	queueNameReachable   = "autodiscovery_reacheable"
	queueNameUnreachable = "autodiscovery_unreacheable"
)

type notifier struct {
	amqpURI string
}

func (n notifier) NotifyDiscovery(p discovery.Peer) error {
	b, err := marshalNode(p)
	if err != nil {
		return err
	}
	fmt.Printf("Discovered peer: %s\n", p)
	return n.notifyQueue(b, "application/json", queueNameDiscoveries)
}

func (n notifier) NotifyReachable(publicKey string) error {
	return n.notifyQueue([]byte(publicKey), "text/plain", queueNameReachable)
}

func (n notifier) NotifyUnreachable(publicKey string) error {
	return n.notifyQueue([]byte(publicKey), "text/plain", queueNameUnreachable)
}

func (n notifier) notifyQueue(data []byte, contentType string, queueName string) error {
	conn, err := amqp.Dial(n.amqpURI)
	defer conn.Close()

	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	defer ch.Close()

	if err != nil {
		return err
	}
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

	if err := ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType:  contentType,
		Body:         data,
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
	}); err != nil {
		return err
	}

	return nil
}
