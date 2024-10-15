package vpn

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

	// Determine the command based on the operating system
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		// macOS command to open a new Terminal window and run OpenVPN
		cmd = exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Terminal" to do script "sudo %s --config %s"`, c.vpnCommand, filePath))
	case "linux":
		// Linux command to open a new terminal window and run OpenVPN
		cmd = exec.Command("gnome-terminal", "--", "bash", "-c", fmt.Sprintf("sudo %s --config %s; exec bash", c.vpnCommand, filePath))
	case "windows":
		// Windows command to open a new Command Prompt window and run OpenVPN
		cmd = exec.Command("cmd", "/C", fmt.Sprintf("start cmd /K \"sudo %s --config %s\"", c.vpnCommand, filePath))
	case "freebsd", "openbsd":
		// FreeBSD and OpenBSD command to open a new terminal window and run OpenVPN
		cmd = exec.Command("xterm", "-e", fmt.Sprintf("sudo %s --config %s", c.vpnCommand, filePath))
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println(cmd.String())
	return cmd.Start()
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
