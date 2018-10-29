package amqp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
	"github.com/uniris/uniris-core/autodiscovery/pkg/system"
)

const (
	queueNameDiscoveries = "autodiscovery_discoveries"
	queueNameReachable   = "autodiscovery_reacheable"
	queueNameUnreachable = "autodiscovery_unreacheable"
)

type notifier struct {
	amqpURI string
}

//NotifyDiscoveries notifies for a new peers that has been discovered
func (n notifier) NotifyDiscoveries(p discovery.Peer) error {
	b, err := n.serialize(p)
	if err != nil {
		return err
	}
	return n.notifyQueue(b, queueNameDiscoveries)
}

//NotifyReachable notifies for an unreachable peer which is now reachable
func (n notifier) NotifyReachable(pubk string) error {
	b, err := n.serializePublicKey(pubk)
	if err != nil {
		return err
	}
	return n.notifyQueue(b, queueNameReachable)
}

//NotifyUnreachable notifies for an unreachable peer
func (n notifier) NotifyUnreachable(pubk string) error {
	b, err := n.serializePublicKey(pubk)
	if err != nil {
		return err
	}
	return n.notifyQueue(b, queueNameUnreachable)
}

func (n notifier) notifyQueue(data []byte, queueName string) error {
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
		ContentType:  "application/json",
		Body:         data,
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
	}); err != nil {
		return err
	}
	return nil
}

func (n notifier) serializePublicKey(pubk string) ([]byte, error) {
	type peerPubk struct {
		PublicKey string `json:"pubKey"`
	}

	puk := peerPubk{
		PublicKey: pubk,
	}

	b, err := json.Marshal(puk)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (n notifier) serialize(p discovery.Peer) ([]byte, error) {
	type PeerPosition struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}

	type PeerIdentity struct {
		PublicKey string `json:"publicKey"`
		IP        string `json:"ip"`
		Port      int    `json:"port"`
	}

	type PeerHeartbeatState struct {
		GenerationTime    int64 `json:"generationTime"`
		ElapsedHeartbeats int64 `json:"elapsedHeartbeats"`
	}
	type PeerAppState struct {
		Status                string       `json:"status"`
		CPULoad               string       `json:"cpuLoad"`
		FreeDiskSpace         float64      `json:"freeDiskSpace"`
		Version               string       `json:"version"`
		GeoPosition           PeerPosition `json:"geoPosition"`
		P2PFactor             int          `json:"p2pFactor"`
		DiscoveredPeersNumber int          `json:"discoveredPeersNumber"`
	}

	return json.Marshal(&struct {
		Identity       PeerIdentity
		HeartbeatState PeerHeartbeatState
		AppState       PeerAppState
	}{
		Identity: PeerIdentity{
			PublicKey: p.Identity().PublicKey(),
			IP:        p.Identity().IP().String(),
			Port:      p.Identity().Port(),
		},
		HeartbeatState: PeerHeartbeatState{
			GenerationTime:    p.HeartbeatState().GenerationTime().Unix(),
			ElapsedHeartbeats: p.HeartbeatState().ElapsedHeartbeats(),
		},
		AppState: PeerAppState{
			CPULoad: p.AppState().CPULoad(),
			Status:  p.AppState().Status().String(),
			Version: p.AppState().Version(),
			GeoPosition: PeerPosition{
				Lat: p.AppState().GeoPosition().Lat,
				Lon: p.AppState().GeoPosition().Lon,
			},
			FreeDiskSpace: p.AppState().FreeDiskSpace(),
			P2PFactor:     p.AppState().P2PFactor(),
		},
	})
}

//NewNotifier creates an amqp implementation of the gossip Notifier interface
func NewNotifier(conf system.AMQPConfig) gossip.Notifier {
	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%d", conf.Username, conf.Password, conf.Host, conf.Port)
	return notifier{amqpURI}
}
