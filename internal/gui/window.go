package gui

import (
	"LibreGate/internal/vpn"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

type Window struct {
	fetcher   *vpn.Fetcher
	connector *vpn.Connector
	app       fyne.App
	window    fyne.Window
}

func NewWindow(fetcher *vpn.Fetcher, connector *vpn.Connector) *Window {
	return &Window{
		fetcher:   fetcher,
		connector: connector,
		app:       app.New(),
	}
}

func (w *Window) Run() error {
	w.window = w.app.NewWindow("LibreGate : Free & Open Source VPN Client for VPN Gate")
	w.window.Resize(fyne.NewSize(900, 1000))

	connectionWidget := NewConnectionWidget(w.connector)
	vpnList := NewVPNList(w.fetcher, w.connector, connectionWidget, w.showError)

	content := container.NewBorder(
		connectionWidget.Container(),
		nil,
		nil,
		nil,
		vpnList,
	)

	w.window.SetContent(content)
	w.window.ShowAndRun()

	return nil
}

func (w *Window) showError(err error) {
	dialog.ShowError(err, w.window)
}
