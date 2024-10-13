package vpngate

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const vpnGateCSVURL = "http://www.vpngate.net/api/iphone/"

type Client struct {
	httpClient *http.Client
}

type Server struct {
	HostName                  string
	IP                        string
	Score                     int
	Ping                      int
	Speed                     int
	CountryLong               string
	CountryShort              string
	NumVPNSessions            int
	Uptime                    int
	TotalUsers                int
	TotalTraffic              float64
	LogType                   string
	Operator                  string
	Message                   string
	OpenVPN_ConfigData_Base64 string
}

func NewClient() *Client {
	return &Client{httpClient: &http.Client{}}
}

func (c *Client) GetServers() ([]Server, error) {
	resp, err := c.httpClient.Get(vpnGateCSVURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	var servers []Server
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Print the record array
		//fmt.Println(record)

		if record[0] == "*" || record[0] == "*vpn_servers" {
			continue
		}

		server, err := parseServer(record)
		if err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}

	return servers, nil
}

func parseServer(record []string) (Server, error) {
	if len(record) < 15 {
		return Server{}, fmt.Errorf("invalid record length")
	}

	score, _ := strconv.Atoi(record[2])
	ping, _ := strconv.Atoi(record[3])
	speed, _ := strconv.Atoi(record[4])
	numSessions, _ := strconv.Atoi(record[7])
	uptime, _ := strconv.Atoi(record[8])
	totalUsers, _ := strconv.Atoi(record[9])
	totalTraffic, _ := strconv.ParseFloat(record[10], 64)

	return Server{
		HostName:                  record[0],
		IP:                        record[1],
		Score:                     score,
		Ping:                      ping,
		Speed:                     speed,
		CountryLong:               record[5],
		CountryShort:              record[6],
		NumVPNSessions:            numSessions,
		Uptime:                    uptime,
		TotalUsers:                totalUsers,
		TotalTraffic:              totalTraffic,
		LogType:                   record[11],
		Operator:                  record[12],
		Message:                   record[13],
		OpenVPN_ConfigData_Base64: record[14],
	}, nil
}
