# LibreGate: Free & Open Source VPN Client for VPN Gate

<img src=./icon.png width="50" alt="Free & Open Source VPN Gate (VPNGate) Client for Linux & macOS">

## Overview

LibreGate is a free and open-source VPN client for **Linux** designed to work seamlessly with VPN Gate . It provides a user-friendly interface to connect to VPN servers, ensuring your online privacy and security.
VPN Gate is a free public VPN relay service that allows users to connect to thousands of VPN servers worldwide. LibreGate simplifies the process of connecting to these servers by providing an easy-to-use interface and automating the connection process. 
VPN Gate project is developed by the University of Tsukuba, Japan.
## Features

- **Easy-to-use Interface**: Simple and intuitive GUI built with Fyne.
- **Server Management**: Fetch and load VPN servers with ease.
- **Secure Connections**: Connect to VPN servers using OpenVPN.
- **Cross-Platform**: Works on Linux & macOS (Windows, iOS & Android support planned).

## Screenshots

![Screenshot 1](./assets/LibreGate-VPN-Linux.png "Linux VPN Client for VPN Gate")

[//]: # (![Screenshot 2]&#40;path/to/screenshot2.png&#41;)

## Usage

1. Install OpenVPN:
    ```sh
    sudo apt install openvpn
    ```
2. Download the binary & launch LibreGate:
    ```sh
    ./libregate
    ```
2. Click on the "Fetch Servers" button to retrieve the latest VPN servers.
3. Select a server from the list and click "Connect". This will establish a VPN connection in a new terminal window.
4. Enjoy having easy access to thousands of free VPNs across many countries worldwide and a secure and private browsing experience.

## Installation

### Prerequisites

- Go 1.16 or higher
- OpenVPN
- Fyne

### Steps

1. Clone the repository:
    ```sh
    git clone https://github.com/zrebar/libregate.git
    cd libregate
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

3. Build the project:
    ```sh
    go build -o libregate ./cmd/libregate/main.go
    ```

4. Run the application:
    ```sh
    ./libregate
    ```


## Configuration

LibreGate uses a configuration file located at `~/.config/libregate/config.json`. You can customize the VPN command if needed, modify this only if you want to use a different VPN client (that supports VPNGate).

```json
{
    "vpn_command": "openvpn"
}
```
## Disclaimer

This project is not affiliated with VPN Gate or the University of Tsukuba. It is an independent project developed by the community to provide an easy-to-use VPN client for VPN Gate on Linux & macOS. Use it at your own risk.