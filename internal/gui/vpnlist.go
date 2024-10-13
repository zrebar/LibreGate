package gui

import (
	"fmt"
	"sort"

	"LibreGate/internal/vpn"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type VPNList struct {
	fetcher          *vpn.Fetcher
	connector        *vpn.Connector
	list             *widget.List
	servers          []vpn.VPNServer
	showError        func(error)
	refreshBtn       *widget.Button
	statusLabel      *widget.Label
	connectionWidget *ConnectionWidget
}

func NewVPNList(fetcher *vpn.Fetcher, connector *vpn.Connector, connectionWidget *ConnectionWidget, showError func(error)) *fyne.Container {
	vpnList := &VPNList{
		fetcher:          fetcher,
		connector:        connector,
		showError:        showError,
		statusLabel:      widget.NewLabel(""),
		connectionWidget: connectionWidget,
	}

	vpnList.list = widget.NewList(
		func() int { return len(vpnList.servers) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Country"),
				widget.NewLabel("Hostname"),
				widget.NewLabel("Speed"),
				widget.NewLabel("Ping"),
				widget.NewButton("Connect", func() {}),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			server := vpnList.servers[id]
			items := item.(*fyne.Container).Objects
			items[0].(*widget.Label).SetText(server.Country)
			items[1].(*widget.Label).SetText(server.Hostname)
			items[2].(*widget.Label).SetText(fmt.Sprintf("%d Mbps", server.Speed))
			items[3].(*widget.Label).SetText(fmt.Sprintf("%d ms", server.Ping))
			items[4].(*widget.Button).OnTapped = func() {
				go vpnList.connectToVPN(server)
			}
		},
	)

	vpnList.refreshBtn = widget.NewButton("Refresh", func() {
		go vpnList.refreshServers()
	})

	// Load VPN list from file if it exists
	vpnList.loadServersFromFile()

	return container.NewBorder(
		container.NewVBox(vpnList.refreshBtn, vpnList.statusLabel),
		nil,
		nil,
		nil,
		vpnList.list,
	)
}

func (v *VPNList) refreshServers() {
	v.refreshBtn.Disable()
	v.statusLabel.SetText("Fetching servers...")

	servers, err := v.fetcher.FetchVPNList()
	if err != nil {
		v.showError(fmt.Errorf("failed to fetch VPN servers: %v", err))
		v.statusLabel.SetText("Failed to fetch servers")
		v.refreshBtn.Enable()
		return
	}

	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Score > servers[j].Score
	})

	v.servers = servers
	v.list.Refresh()
	v.statusLabel.SetText(fmt.Sprintf("Found %d servers", len(servers)))
	v.refreshBtn.Enable()
}

func (v *VPNList) loadServersFromFile() {
	servers, err := v.fetcher.LoadVPNList()
	if err != nil {
		v.showError(fmt.Errorf("failed to load VPN servers from file: %v", err))
		v.statusLabel.SetText("Failed to load servers")
		return
	}

	if servers != nil {
		sort.Slice(servers, func(i, j int) bool {
			return servers[i].Score > servers[j].Score
		})

		v.servers = servers
		v.list.Refresh()
		v.statusLabel.SetText(fmt.Sprintf("Loaded %d servers from file", len(servers)))
	}
}

func (v *VPNList) connectToVPN(server vpn.VPNServer) {
	v.statusLabel.SetText(fmt.Sprintf("Connecting to %s...", server.Hostname))
	err := v.connector.Connect(server)
	if err != nil {
		v.showError(fmt.Errorf("failed to connect to VPN: %v", err))
		v.statusLabel.SetText("Connection failed")
		return
	}
	v.connectionWidget.SetConnected(server)
	v.statusLabel.SetText(fmt.Sprintf("Connected to %s", server.Hostname))
}
