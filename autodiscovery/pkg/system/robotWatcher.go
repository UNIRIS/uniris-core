package system

import (
	"fmt"
	"net"

	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
)

type robotWatcher struct{}

func NewRobotWatcher() monitoring.RobotWatcher {
	return robotWatcher{}
}

//CheckAutodiscovery check that the autodiscovery daemon is running
func (r robotWatcher) CheckAutodiscoveryProcess(port int) error {
	_, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return err
	}
	return nil
}

//CheckDataProcess check that the autodiscovery daemon is running
func (r robotWatcher) CheckDataProcess() error {
	// TBD
	return nil
}

//CheckMiningProcess check that the autodiscovery daemon is running
func (r robotWatcher) CheckMiningProcess() error {
	//TBD
	return nil
}

//CheckAIProcess check that the autodiscovery daemon is running
func (r robotWatcher) CheckAIProcess() error {
	//TBD
	return nil
}

//CheckRedisProcess check that the ScyallaDB daemon is running
func (r robotWatcher) CheckScyllaDbProcess() error {
	//TBD
	return nil
}

//CheckRedisProcess check that the redis daemon is running
func (r robotWatcher) CheckRedisProcess() error {
	//TBD
	return nil
}

//CheckRabbitmqProcess check that the rabbitmq daemon is running
func (r robotWatcher) CheckRabbitmqProcess() error {
	//TBD
	return nil
}
