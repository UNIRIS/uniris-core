package mock

//Monitor mock
type Monitor struct{}

//CPULoad mock
func (w Monitor) CPULoad() (string, error) {
	return "0.62 0.77 0.71 4/972 26361", nil
}

//FreeDiskSpace mock
func (w Monitor) FreeDiskSpace() (float64, error) {
	return 212383852, nil
}

//P2PFactor mock
func (w Monitor) P2PFactor() (int, error) {
	return 1, nil
}
