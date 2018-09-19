package http

import (
	"encoding/json"
	"net/http"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//PeerPositioner imlplements the PeerPositioner interface which provides methods to get the geo position of the peer
type PeerPositioner struct{}

//Position returns the peer's geo position
func (loc PeerPositioner) Position() (pos discovery.PeerPosition, err error) {
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
