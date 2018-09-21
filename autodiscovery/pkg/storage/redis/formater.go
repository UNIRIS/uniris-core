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
		"publicKey":      string(p.PublicKey()),
		"port":           strconv.Itoa(p.Port()),
		"ip":             p.IP().String(),
		"generationTime": strconv.Itoa(int(p.GenerationTime().Unix())),
		"isOwned":        strconv.FormatBool(p.IsOwned()),
		"status":         string(p.Status()),
		"cpuLoad":        p.CPULoad(),
		"ioWaitRate":     fmt.Sprintf("%f", p.IOWaitRate()),
		"freeDiskSpace":  fmt.Sprintf("%f", p.FreeDiskSpace()),
		"version":        p.Version(),
		"geoPosition":    fmt.Sprintf("%f;%f", p.GeoPosition().Lat, p.GeoPosition().Lon),
		"p2pFactor":      string(p.P2PFactor()),
	}
}

//FormatHashToPeer converts a Redis hashset into a peer
func FormatHashToPeer(h map[string]string) discovery.Peer {

	pbKey := []byte(h["publicKey"])
	port, _ := strconv.Atoi(h["port"])
	ip := net.ParseIP(h["ip"])

	gen, _ := strconv.Atoi(h["generationTime"])
	generationTime := time.Unix(int64(gen), 0)
	isOwned, _ := strconv.ParseBool(h["isOwned"])

	s, _ := strconv.Atoi(h["status"])
	status := discovery.PeerStatus(s)

	cpuLoad := h["cpuLoad"]
	ioWaitRate, _ := strconv.ParseFloat(h["ioWaitRate"], 64)
	freeDiskSpace, _ := strconv.ParseFloat(h["freeDiskSpace"], 64)
	version := h["version"]
	p2pFactor, _ := strconv.Atoi(h["p2pFactor"])
	posArr := strings.Split(h["geoPosition"], ";")

	lat, _ := strconv.ParseFloat(posArr[0], 64)
	lon, _ := strconv.ParseFloat(posArr[1], 64)

	pos := discovery.PeerPosition{Lat: lat, Lon: lon}

	state := discovery.NewState(version, status, pos, cpuLoad, freeDiskSpace, ioWaitRate, p2pFactor)
	p := discovery.NewPeerDetailed(pbKey, ip, port, generationTime, isOwned, state)
	return p
}

//FormatSeedToHash converts a seed into a Redis hashset
func FormatSeedToHash(seed discovery.Seed) map[string]interface{} {
	return map[string]interface{}{
		"ip":   seed.IP.String(),
		"port": strconv.Itoa(seed.Port),
	}
}

//FormatHashToSeed converts a Redis hashset into a seed
func FormatHashToSeed(hash map[string]string) discovery.Seed {
	port, _ := strconv.Atoi(hash["port"])
	return discovery.Seed{
		IP:   net.ParseIP(hash["ip"]),
		Port: port,
	}
}
