package mock

type Monitor struct{}

func (w Monitor) CPULoad() (string, error) {
	return "0.62 0.77 0.71 4/972 26361", nil
}

func (w Monitor) FreeDiskSpace() (float64, error) {
	return 212383852, nil
}

func (w Monitor) P2PFactor() (int, error) {
	return 1, nil
}

func (w Monitor) CheckAutodiscoveryProcess(discoveryPort int) error {
	return nil
}
func (w Monitor) CheckMiningProcess() error {
	return nil
}
func (w Monitor) CheckDataProcess() error {
	return nil
}
func (w Monitor) CheckAIProcess() error {
	return nil
}
func (w Monitor) CheckRedisProcess() error {
	return nil
}
func (w Monitor) CheckScyllaDbProcess() error {
	return nil
}
func (w Monitor) CheckRabbitmqProcess() error {
	return nil
}
