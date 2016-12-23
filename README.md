#Azure CS:GO Server Manager

###A tool for managing CS:GO servers hosted on Azure
This tool allows for servers to be deployed, monitored and managed on Azure using various templates and server configuration files.

## Features
- [x] Monitor running servers (VMs) in Azure Resource Group
- [x] Manage deployment templates and parameters
- [x] Deploy multiple servers using deployment templates
- [x] Delete servers and their associated resources and hard drives (.vhd's)
- [x] Upload and change server configs (server.cfg)
- [ ] View number of players connected to server
- [ ] Execute RCON commands on the server
- [ ] Allow viewing of the server logs
- [x] Authentication for protecting against unauthorized users

## Installation
1. Download the latest release
  * [https://github.com/MetalMichael/azure-csgo-server-manager/releases](https://github.com/MajorMJR/factorio-server-manager/releases)
2. Configure the config.json file to your [azure credentials](https://docs.microsoft.com/en-gb/azure/active-directory/active-directory-protocols-oauth-code) 
3. Run the server binary file
4. Visit [localhost:8080](localhost:8080) in your web browser.

## Development
The backend is built as a REST API via the Go web application.  

It also acts as the webserver to serve the front end react application

All api actions are accessible with the /api route.  The frontend is accessible from /.

#### Requirements
+ Go 1.6
+ NodeJS 4.2.6

#### Building the Go backend
Go Application which manages the Factorio server.

API requests for managing the Factorio server are sent to /api.

The frontend code is served by a HTTP file server running on /.
```
git clone https://github.com/MetalMichael/azure-csgo-server-manager.git
cd azure-csgo-server-manager
go build
```

#### Building the React frontend
Frontend is built using React and the AdminLTE CSS framework. See app/dist/ for AdminLTE included files and license.

The root of the UI application is served at app/index.html.  Run the npm build script and the Go application during development to get live rebuilding of the UI code.

All necessary CSS and Javascript files are included for running the UI.

Transpiled bundle.js application is output to app/bundle.js, 'npm run build' script starts webpack to build the React application for development
```
 install nodejs (use nvm)
 cd ui/
 npm install
 npm run build
 Start azure-csgo-server-manager binary in another terminal
```

Check out [factorio-server-manager](https://github.com/MajorMJR/factorio-server-manager) for more info

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
