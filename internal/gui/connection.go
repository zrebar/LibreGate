package gui

import (
	"fmt"
	"time"

	"LibreGate/internal/vpn"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ConnectionStatus struct {
	widget.Label
	connected bool
	server    *vpn.VPNServer
	startTime time.Time
}

func NewConnectionStatus() *ConnectionStatus {
	status := &ConnectionStatus{}
	status.ExtendBaseWidget(status)
	status.updateDisplay()
	return status
}

func (c *ConnectionStatus) SetConnected(server vpn.VPNServer) {
	c.connected = true
	c.server = &server
	c.startTime = time.Now()
	c.updateDisplay()
}

func (c *ConnectionStatus) SetDisconnected() {
	c.connected = false
	c.server = nil
	c.updateDisplay()
}

func (c *ConnectionStatus) updateDisplay() {
	if !c.connected {
		c.SetText("Not connected")
		return
	}

	duration := time.Since(c.startTime).Round(time.Second)
	c.SetText(fmt.Sprintf("Connected to %s (%s) for %s", c.server.Hostname, c.server.Country, duration))
}

func (c *ConnectionStatus) StartTimer() {
	go func() {
		for c.connected {
			time.Sleep(time.Second)
			c.updateDisplay()
		}
	}()
}

type ConnectionWidget struct {
	container     *fyne.Container
	status        *ConnectionStatus
	disconnectBtn *widget.Button
}

func NewConnectionWidget(connector *vpn.Connector) *ConnectionWidget {
	status := NewConnectionStatus()
	disconnectBtn := widget.NewButton("Disconnect", nil)
	disconnectBtn.Disable()

	cw := &ConnectionWidget{
		status:        status,
		disconnectBtn: disconnectBtn,
	}

	disconnectBtn.OnTapped = func() {
		go func() {
			err := connector.Disconnect()
			if err != nil {
				//fmt.Println(err)
				return
			}
			cw.SetDisconnected()
		}()
	}

	cw.container = container.NewVBox(status, disconnectBtn)
	return cw
}

func (cw *ConnectionWidget) SetConnected(server vpn.VPNServer) {
	cw.status.SetConnected(server)
	cw.status.StartTimer()
	cw.disconnectBtn.Enable()
}

func (cw *ConnectionWidget) SetDisconnected() {
	cw.status.SetDisconnected()
	cw.disconnectBtn.Disable()
}

func (cw *ConnectionWidget) Container() *fyne.Container {
	return cw.container
}
