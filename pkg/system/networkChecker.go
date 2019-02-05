package system

import (
	"fmt"
	"net"
	"time"

	"github.com/beevik/ntp"
	"github.com/uniris/uniris-core/pkg/discovery"
)

const (
	cdns          = "uniris.io"
	ntpRetry      = 3
	upmaxOffset   = 300
	downmaxOffset = -300
)

var cntp = [...]string{"1.pool.ntp.org", "2.pool.ntp.org", "3.pool.ntp.org", "4.pool.ntp.org"}

type netCheck struct {
	internGrpcPort int
	externGrpcPort int
}

//NewNetworkChecker creates new network checker
func NewNetworkChecker(internGrpcPort, externGrpcPort int) discovery.NetworkChecker {
	return netCheck{
		internGrpcPort: internGrpcPort,
		externGrpcPort: externGrpcPort,
	}
}

func (n netCheck) CheckNtpState() error {
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
				return discovery.ErrNTPShift
			}
			return nil
		}
	}
	return discovery.ErrNTPFailure
}

func (n netCheck) CheckInternetState() error {
	_, err := net.LookupIP(cdns)
	if err != nil {
		return err
	}
	return nil
}

func (n netCheck) CheckGRPCServers() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", n.internGrpcPort))
	if err != nil {
		return discovery.ErrGRPCServer
	}

	var buffer []byte
	if _, err := conn.Read(buffer); err != nil {
		conn.Close()
		return discovery.ErrGRPCServer
	}

	conn2, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", n.externGrpcPort))
	if err != nil {
		return discovery.ErrGRPCServer
	}
	var buffer2 []byte
	if _, err := conn2.Read(buffer2); err != nil {
		conn2.Close()
		return discovery.ErrGRPCServer
	}
	return nil
}
