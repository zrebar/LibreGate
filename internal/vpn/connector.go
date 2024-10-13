package vpn

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type Connector struct {
	vpnCommand string
	cmd        *exec.Cmd
}

func NewConnector(vpnCommand string) *Connector {
	return &Connector{vpnCommand: vpnCommand}
}

func (c *Connector) Connect(server VPNServer) error {
	configData, err := base64.StdEncoding.DecodeString(server.OpenVPN)
	if err != nil {
		return fmt.Errorf("failed to decode OpenVPN config: %v", err)
	}

	filePath := filepath.Join(".", "openvpn-config.ovpn")
	if err := ioutil.WriteFile(filePath, configData, 0644); err != nil {
		return fmt.Errorf("failed to write OpenVPN config: %v", err)
	}

	// Command to execute OpenVPN in a new terminal window and keep it open
	c.cmd = exec.Command("gnome-terminal", "--", "bash", "-c", fmt.Sprintf("sudo %s --config %s; exec bash", c.vpnCommand, filePath))
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
	fmt.Println(c.cmd.String())
	return c.cmd.Start()
}

func (c *Connector) Disconnect() error {
	if c.cmd == nil || c.cmd.Process == nil {
		return nil
	}

	err := c.cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return fmt.Errorf("failed to send interrupt signal: %v", err)
	}

	err = c.cmd.Wait()
	if err != nil {
		return fmt.Errorf("failed to wait for process to exit: %v", err)
	}

	c.cmd = nil
	return nil
}
