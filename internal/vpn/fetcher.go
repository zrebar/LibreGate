package vpn

import (
	vpngate "LibreGate/pkg"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Fetcher struct {
	client *vpngate.Client
}

func NewFetcher(client *vpngate.Client) *Fetcher {
	return &Fetcher{client: client}
}

func (f *Fetcher) FetchVPNList() ([]VPNServer, error) {
	servers, err := f.client.GetServers()
	if err != nil {
		return nil, err
	}

	var vpnServers []VPNServer
	for _, server := range servers {
		vpnServers = append(vpnServers, VPNServer{
			IP:       server.IP,
			Country:  server.CountryLong,
			Speed:    server.Speed,
			Ping:     server.Ping,
			Score:    server.Score,
			OpenVPN:  server.OpenVPN_ConfigData_Base64,
			Hostname: server.HostName,
		})
	}

	// Save the VPN list to a file
	err = f.SaveVPNList(vpnServers)
	if err != nil {
		return nil, err
	}

	return vpnServers, nil
}

func (f *Fetcher) SaveVPNList(vpnServers []VPNServer) error {
	data, err := json.Marshal(vpnServers)
	if err != nil {
		return err
	}

	filePath := filepath.Join(".", "vpnlist.json")
	return ioutil.WriteFile(filePath, data, 0644)
}

func (f *Fetcher) LoadVPNList() ([]VPNServer, error) {
	filePath := filepath.Join(".", "vpnlist.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var vpnServers []VPNServer
	err = json.Unmarshal(data, &vpnServers)
	if err != nil {
		return nil, err
	}

	return vpnServers, nil
}
