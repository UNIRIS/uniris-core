package mock

// RobotWatcher mock
type RobotWatcher struct{}

//CheckAutodiscoveryProcess mock
func (r RobotWatcher) CheckAutodiscoveryProcess(port int) error {
	return nil
}

//CheckDataProcess mock
func (r RobotWatcher) CheckDataProcess() error {
	return nil
}

//CheckMiningProcess mock
func (r RobotWatcher) CheckMiningProcess() error {
	return nil
}

//CheckAIProcess check mock
func (r RobotWatcher) CheckAIProcess() error {
	return nil
}

//CheckScyllaDbProcess mock
func (r RobotWatcher) CheckScyllaDbProcess() error {
	return nil
}

//CheckRedisProcess mock
func (r RobotWatcher) CheckRedisProcess() error {
	return nil
}

//CheckRabbitmqProcess mock
func (r RobotWatcher) CheckRabbitmqProcess() error {
	return nil
}
