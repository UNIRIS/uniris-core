package adapters

import (
	"encoding/json"
	"flag"
	"net"

	"github.com/uniris/uniris-core/autodiscovery/core/domain"
	"github.com/uniris/uniris-core/autodiscovery/infrastructure"
)

//ConfigurationReader retrieves the peer configuration details
type ConfigurationReader struct {
	FileReader infrastructure.FileReader
	HTTPClient infrastructure.HTTPClient

	Port          int
	P2PFactor     int
	PublicKeyFile *string
}

func NewConfigurationReader(r infrastructure.FileReader, c infrastructure.HTTPClient) ConfigurationReader {
	port := flag.Int("port", 3545, "GRPC port")
	p2pFactor := flag.Int("p2p-factor", 1, "P2P replication factor")
	pubKeyFile := flag.String("pub-key-file", "id.pub", "Public key file")

	flag.Parse()

	return ConfigurationReader{
		FileReader:    r,
		HTTPClient:    c,
		Port:          *port,
		P2PFactor:     *p2pFactor,
		PublicKeyFile: pubKeyFile,
	}
}

//GetPublicKey retrieve the peer public key
func (r ConfigurationReader) GetPublicKey() (bytes []byte, err error) {
	if r.PublicKeyFile != nil {
		bytes, err = r.FileReader.ReadFile(*r.PublicKeyFile)
	} else {
		bytes, err = r.FileReader.ReadFile("id.pub")
	}
	return bytes, err
}

//GetPort retrieves the peer port
func (r ConfigurationReader) GetPort() (int, error) {
	return r.Port, nil
}

//GetVersion returns the peer version
func (r ConfigurationReader) GetVersion() (string, error) {
	bytes, err := r.FileReader.ReadFile("version")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

//GetSeeds loads the seed peers from the configuration file
func (r ConfigurationReader) GetSeeds() ([]domain.Peer, error) {
	seedPeerList := make([]domain.Peer, 0)
	bytes, err := r.FileReader.ReadFile("seed.json")
	if err != nil {
		return nil, err
	}

	//Deserialize the seed peers
	json.Unmarshal(bytes, &seedPeerList)
	return seedPeerList, nil
}

//GetGeoPosition retrieves the peer geographic position
func (r ConfigurationReader) GetGeoPosition() (domain.GeoPosition, error) {
	//TODO: call external service
	return domain.GeoPosition{
		Lat: 2.0,
		Lon: 50.1,
	}, nil
}

//GetIP retrieves the peer IP
func (r ConfigurationReader) GetIP() (net.IP, error) {
	//TODO: call external service
	// var ip net.IP
	// r.HTTPClient.GetJSON("", *ip)
	return net.ParseIP("127.0.0.1"), nil
}

//GetP2PFactor retrieves the peer P2P replication factor
func (r ConfigurationReader) GetP2PFactor() (int, error) {
	//TODO
	return 1, nil
}
