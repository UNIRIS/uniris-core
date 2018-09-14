package infrastructure

import (
	"net"

	"github.com/uniris/uniris-core/autodiscovery/domain"
)

//Geolocalizer wraps the geoposition lookup
type Geolocalizer struct{}

//Lookup retrieves the current geolocalization information
func (g Geolocalizer) Lookup() (domain.GeoPosition, error) {
	// resp, err := http.Get("http://ip-api.com/json")
	// if err != nil {
	// 	return nil, err
	// }
	// defer resp.Body.Close()

	// //Deserialize the json
	// decoder := json.NewDecoder(resp.Body)
	// var geo *services.GeoLoc
	// if err = decoder.Decode(&geo); err != nil {
	// 	return nil, err
	// }
	// return geo, nil

	return domain.GeoPosition{
		IP:  net.ParseIP("127.0.0.1"),
		Lat: 2.0,
		Lon: 50.1,
	}, nil
}
