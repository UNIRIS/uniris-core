package redis

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//FormatPeerToHash converts a peer into a Redis hashset
func FormatPeerToHash(p discovery.Peer) map[string]interface{} {
	return map[string]interface{}{
		"publicKey":             p.Identity().PublicKey(),
		"port":                  strconv.Itoa(p.Identity().Port()),
		"ip":                    p.Identity().IP().String(),
		"generationTime":        strconv.Itoa(int(p.HeartbeatState().GenerationTime().Unix())),
		"elapsedHeartbeats":     strconv.Itoa(int(p.HeartbeatState().ElapsedHeartbeats())),
		"status":                string(p.AppState().Status()),
		"cpuLoad":               p.AppState().CPULoad(),
		"freeDiskSpace":         fmt.Sprintf("%f", p.AppState().FreeDiskSpace()),
		"version":               p.AppState().Version(),
		"geoPosition":           fmt.Sprintf("%f;%f", p.AppState().GeoPosition().Lat, p.AppState().GeoPosition().Lon),
		"p2pFactor":             string(p.AppState().P2PFactor()),
		"discoveredPeersNumber": fmt.Sprintf("%d", p.AppState().DiscoveredPeersNumber()),
	}
}

//FormatHashToPeer converts a Redis hashset into a peer
func FormatHashToPeer(hash map[string]string) discovery.Peer {

	pbKey := hash["publicKey"]
	port, _ := strconv.Atoi(hash["port"])
	ip := net.ParseIP(hash["ip"])

	gen, _ := strconv.Atoi(hash["generationTime"])
	generationTime := time.Unix(int64(gen), 0)

	elapsedHeartbeats, _ := strconv.Atoi(hash["elapsedHeartbeats"])
	elpased := int64(elapsedHeartbeats)

	s, _ := strconv.Atoi(hash["status"])
	status := discovery.PeerStatus(s)

	cpuLoad := hash["cpuLoad"]
	freeDiskSpace, _ := strconv.ParseFloat(hash["freeDiskSpace"], 64)
	version := hash["version"]
	p2pFactor, _ := strconv.Atoi(hash["p2pFactor"])
	posArr := strings.Split(hash["geoPosition"], ";")

	lat, _ := strconv.ParseFloat(posArr[0], 64)
	lon, _ := strconv.ParseFloat(posArr[1], 64)

	pos := discovery.PeerPosition{Lat: lat, Lon: lon}

	dpn, _ := strconv.Atoi(hash["discoveredPeersNumber"])

	p := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(ip, port, pbKey),
		discovery.NewPeerHeartbeatState(generationTime, elpased),
		discovery.NewPeerAppState(version, status, pos, cpuLoad, freeDiskSpace, p2pFactor, dpn),
	)
	return p
}

//FormatSeedToHash converts a seed into a Redis hashset
func FormatSeedToHash(seed discovery.Seed) map[string]interface{} {
	return map[string]interface{}{
		"ip":        seed.IP.String(),
		"port":      strconv.Itoa(seed.Port),
		"publicKey": seed.PublicKey,
	}
}

//FormatHashToSeed converts a Redis hashset into a seed
func FormatHashToSeed(hash map[string]string) discovery.Seed {
	port, _ := strconv.Atoi(hash["port"])
	return discovery.Seed{
		IP:        net.ParseIP(hash["ip"]),
		Port:      port,
		PublicKey: hash["publicKey"],
	}
}
