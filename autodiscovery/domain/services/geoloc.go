package services

import (
	"net"
)

//GeoLoc defines the geolocalization result
type GeoLoc struct {
	IP  net.IP
	Lat float64
	Lon float64
}

//GeolocService represents the interface for peer geolocalization operations
type GeolocService interface {
	Lookup() (*GeoLoc, error)
}
