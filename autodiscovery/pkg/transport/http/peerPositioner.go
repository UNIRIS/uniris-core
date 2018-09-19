package http

import (
	"encoding/json"
	"net/http"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

type PeerPositioner struct{}

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
