package mock

import (
	"github.com/uniris/uniris-core/api/pkg/listing"
)

type client struct {
	listing.RobotClient
}

func NewClient() client {
	return client{}
}
