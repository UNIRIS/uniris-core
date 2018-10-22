package amqp

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
	"github.com/uniris/uniris-core/autodiscovery/pkg/system"
)

const (
	queueName = "discoveries"
)

type notifier struct {
	amqpURI string
}

//Notify notifies a new peers has been discovered
func (n notifier) Notify(p discovery.Peer) error {
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
	b, err := n.serialize(p)
	if err != nil {
		return err
	}
	if err := ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         b,
		DeliveryMode: amqp.Persistent,
	}); err != nil {
		return err
	}
	return nil
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
