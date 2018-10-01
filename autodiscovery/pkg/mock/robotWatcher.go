package mock

type RobotWatcher struct{}

func (r RobotWatcher) CheckAutodiscoveryProcess(port int) error {
	return nil
}

func (r RobotWatcher) CheckDataProcess() error {
	return nil
}

func (r RobotWatcher) CheckMiningProcess() error {
	return nil
}

func (r RobotWatcher) CheckAIProcess() error {
	return nil
}

func (r RobotWatcher) CheckScyllaDbProcess() error {
	return nil
}

func (r RobotWatcher) CheckRedisProcess() error {
	return nil
}

func (r RobotWatcher) CheckRabbitmqProcess() error {
	return nil
}
