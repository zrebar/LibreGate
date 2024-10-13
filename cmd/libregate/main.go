package main

import (
	"LibreGate/internal/config"
	"LibreGate/internal/gui"
	"LibreGate/internal/vpn"
	vpngate "LibreGate/pkg"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client := vpngate.NewClient()
	fetcher := vpn.NewFetcher(client)
	connector := vpn.NewConnector(cfg.VPNCommand)

	window := gui.NewWindow(fetcher, connector)
	if err := window.Run(); err != nil {
		log.Fatalf("Failed to run GUI: %v", err)
	}
}
