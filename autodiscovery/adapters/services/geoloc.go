package services

import (
	"net"

	"github.com/uniris/uniris-core/autodiscovery/domain/services"
)

//GeoService implements the IGeoService interface
type GeoService struct{}

//Lookup retrieves the current geolocalization information
func (g GeoService) Lookup() (*services.GeoLoc, error) {

	// //Request the ip-api web api
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

	return &services.GeoLoc{
		IP:  net.ParseIP("127.0.0.1"),
		Lat: 2.0,
		Lon: 50.1,
	}, nil
}
