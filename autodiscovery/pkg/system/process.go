package system

import (
	"net"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//CheckAutodiscoveryProcess check that the autodiscovery daemon is running
func CheckAutodiscoveryProcess(rep discovery.Repository) (bool, error) {
	selfpeer, err := rep.GetOwnedPeer()
	if err != nil {
		return false, err
	}
	selfport := selfpeer.Port()
	_, err = net.Dial("tcp", ":"+string(selfport))
	if err != nil {
		return false, err
	}
	return true, nil
}

//CheckDataProcess check that the autodiscovery daemon is running
func CheckDataProcess(rep discovery.Repository) (bool, error) {
	// TBD
	return true, nil
}

//CheckMiningProcess check that the autodiscovery daemon is running
func CheckMiningProcess(rep discovery.Repository) (bool, error) {
	//TBD
	return true, nil
}

//CheckAIProcess check that the autodiscovery daemon is running
func CheckAIProcess(rep discovery.Repository) (bool, error) {
	//TBD
	return true, nil
}

//CheckScyllaProcess check that the autodiscovery daemon is running
func CheckScyllaProcess(rep discovery.Repository) (bool, error) {
	//TBD
	return true, nil
}

//CheckRedisProcess check that the autodiscovery daemon is running
func CheckRedisProcess(rep discovery.Repository) (bool, error) {
	//TBD
	return true, nil
}

//CheckRabitmqProcess check that the autodiscovery daemon is running
func CheckRabitmqProcess(rep discovery.Repository) (bool, error) {
	//TBD
	return true, nil
}
