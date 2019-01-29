package transaction

import "net"

type Pool []PoolMember

type PoolMember struct {
	ip   net.IP
	port int
}

func NewPoolMember(ip net.IP, port int) PoolMember {
	return PoolMember{
		ip:   ip,
		port: port,
	}
}

func (pm PoolMember) IP() net.IP {
	return pm.ip
}

func (pm PoolMember) Port() int {
	return pm.port
}
