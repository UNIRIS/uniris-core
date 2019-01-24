package system

import (
	"net"
	"time"

	"github.com/beevik/ntp"
	"github.com/uniris/uniris-core/pkg/gossip"
)

const (
	cdns          = "uniris.io"
	ntpRetry      = 3
	upmaxOffset   = 300
	downmaxOffset = -300
)

var cntp = [...]string{"1.pool.ntp.org", "2.pool.ntp.org", "3.pool.ntp.org", "4.pool.ntp.org"}

type pNet struct{}

func NewPeerNetworker() gossip.PeerNetworker {
	return pNet{}
}

func (pn pNet) CheckNtpState() error {
	for _, ntps := range cntp {
		r, err := ntp.QueryWithOptions(ntps, ntp.QueryOptions{Version: 4})
		if err == nil {
			if (int64(r.ClockOffset/time.Second) < downmaxOffset) || (int64(r.ClockOffset/time.Second) > upmaxOffset) {
				for i := 0; i < ntpRetry; i++ {
					r, err := ntp.QueryWithOptions(ntps, ntp.QueryOptions{Version: 4})
					if err == nil {
						if (int64(r.ClockOffset/time.Second) > downmaxOffset) || (int64(r.ClockOffset/time.Second) < upmaxOffset) {
							return nil
						}
					}
				}
				return gossip.ErrNTPShift
			}
			return nil
		}
	}
	return gossip.ErrNTPFailure
}

func (pn pNet) CheckInternetState() error {
	_, err := net.LookupIP(cdns)
	if err != nil {
		return err
	}
	return nil
}
