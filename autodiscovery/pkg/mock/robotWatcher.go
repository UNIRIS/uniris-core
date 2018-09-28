package mock

import (
	"errors"
	"net"
)

// RobotWatcher assure watching all needed processes on a Robot
type RobotWatcher struct{}

//CheckAutodiscoveryProcess check that the autodiscovery daemon is running
func (r RobotWatcher) CheckAutodiscoveryProcess(port int) error {
	return nil
}

//CheckDataProcess check that the autodiscovery daemon is running
func (r RobotWatcher) CheckDataProcess() error {
	return nil
}

//CheckMiningProcess check that the autodiscovery daemon is running
func (r RobotWatcher) CheckMiningProcess() error {
	return nil
}

//CheckAIProcess check that the autodiscovery daemon is running
func (r RobotWatcher) CheckAIProcess() error {
	return nil
}

//CheckScyllaDbProcess check that the ScyallaDB daemon is running
func (r RobotWatcher) CheckScyllaDbProcess() error {
	return nil
}

//CheckRedisProcess check that the redis daemon is running
func (r RobotWatcher) CheckRedisProcess() error {
	return nil
}

//CheckRabbitmqProcess check that the rabbitmq daemon is running
func (r RobotWatcher) CheckRabbitmqProcess() error {
	return nil
}

type mockSystemNetworker2 struct{}

func (n mockSystemNetworker2) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n mockSystemNetworker2) CheckInternetState() error {
	return nil
}

func (n mockSystemNetworker2) CheckNtpState() error {
	return errors.New("System Clock have a big Offset check the ntp configuration of the system")
}

type mockSystemNetworker3 struct{}

func (n mockSystemNetworker3) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n mockSystemNetworker3) CheckInternetState() error {
	return errors.New("required processes are not running")
}

func (n mockSystemNetworker3) CheckNtpState() error {
	return nil
}
