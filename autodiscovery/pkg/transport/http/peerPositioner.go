package http

import (
	"encoding/json"
	"net/http"

	"github.com/uniris/uniris-core/autodiscovery/pkg/bootstraping"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

type httpPositioner struct{}

//Position returns the peer's geo position
func (p httpPositioner) Position() (pos discovery.PeerPosition, err error) {
	resp, err := http.Get("http://ip-api.com/json")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	//Deserialize the json
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&pos); err != nil {
		return
	}
	return
}

//NewPeerPositioner creates an http implementation of the PeerPositionner interface
func NewPeerPositioner() bootstraping.PeerPositionner {
	return httpPositioner{}
}
