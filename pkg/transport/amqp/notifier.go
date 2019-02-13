package amqp

import (
	"encoding/json"
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
	b, err := json.Marshal(map[string]interface{}{
		"identity": map[string]interface{}{
			"ip":         p.Identity().IP().String(),
			"port":       p.Identity().Port(),
			"public_key": p.Identity().PublicKey(),
		},
		"heartbeat_state": map[string]interface{}{
			"generation_time":    p.HeartbeatState().GenerationTime(),
			"elapsed_heartbeats": p.HeartbeatState().ElapsedHeartbeats(),
		},
		"app_state": map[string]interface{}{
			"version":         p.AppState().Version(),
			"status":          p.AppState().Status(),
			"cpu_load":        p.AppState().CPULoad(),
			"free_disk_space": p.AppState().FreeDiskSpace(),
			"geo_position": map[string]float64{
				"latitude":  p.AppState().GeoPosition().Latitude(),
				"longitude": p.AppState().GeoPosition().Longitude(),
			},
			"p2p_factor":             p.AppState().P2PFactor(),
			"discovered_peer_number": p.AppState().DiscoveredPeersNumber(),
		},
	})
	if err != nil {
		return err
	}
	return n.notifyQueue(b, "application/json", queueNameDiscoveries)
}

func (n notifier) NotifyReachable(p discovery.PeerIdentity) error {
	pBytes, err := json.Marshal(map[string]interface{}{
		"ip":         p.IP().String(),
		"port":       p.Port(),
		"public_key": p.PublicKey(),
	})
	if err != nil {
		return err
	}
	return n.notifyQueue(pBytes, "application/json", queueNameReachable)
}

func (n notifier) NotifyUnreachable(p discovery.PeerIdentity) error {
	pBytes, err := json.Marshal(map[string]interface{}{
		"ip":         p.IP().String(),
		"port":       p.Port(),
		"public_key": p.PublicKey(),
	})
	if err != nil {
		return err
	}
	return n.notifyQueue(pBytes, "application/json", queueNameUnreachable)
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
