package ports

import "github.com/uniris/uniris-core/autodiscovery/domain"

//Geolocalizer wraps the geoposition lookup
type Geolocalizer interface {
	Lookup() (domain.GeoPosition, error)
}
