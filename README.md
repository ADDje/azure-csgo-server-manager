#Azure CS:GO Server Manager

###A tool for managing CS:GO servers hosted on Azure
This tool allows for servers to be deployed, monitored and managed on Azure using various templates and server configuration files.

## Features
- [x] Monitor running servers (VMs) in Azure Resource Group
- [x] Manage deployment templates and parameters
- [x] Deploy multiple servers using deployment templates
- [x] Delete servers and their associated resources and hard drives (.vhd's)
- [x] Upload and change game server configs (server.cfg)
- [ ] View number of players connected to game server
- [ ] Execute RCON commands on the game server
- [ ] Allow viewing of the game server logs
- [x] View management server logs of deployments etc.
- [x] Scheduler. Support for automation hits from [Azure Scheduler](https://azure.microsoft.com/en-gb/services/scheduler/) to complete various preconfigured tasks
- [x] Authentication for protecting against unauthorized users
- [x] Export demos (.dem) to permanent storage

## Future
* Permissions Support. Users and admins (+other?). Users can't see things like storage keys. Maybe they can't delete VMs, etc.
* Cache things. Only really designed for one user right now, could have performance issues
* Expire cookies (configure)
* External API? Could be used for things like automatically configuring maps or team names
* Reboot server manually (or just restart hlds)

## Installation
1. Download the latest release
  * [https://github.com/MetalMichael/azure-csgo-server-manager/releases](https://github.com/MetalMichael/azure-csgo-server-manager/releases)
2. Configure the config.json file to your [azure credentials](https://docs.microsoft.com/en-gb/azure/active-directory/active-directory-protocols-oauth-code) 
3. Run the server binary file
4. Visit [localhost:8090](localhost:8090) in your web browser.

## Development
The backend is built as a REST API via the Go web application.  

It also acts as the webserver to serve the front end react application

All api actions are accessible with the /api route.  The frontend is accessible from /.

#### Requirements
+ Go 1.6
+ NodeJS 4.2.6

#### Building the Go backend
Go Application which manages the CSGO servers.

API requests for managing the servers are sent to /api.

The frontend code is served by a HTTP file server running on /.
```
git clone https://github.com/MetalMichael/azure-csgo-server-manager.git
cd azure-csgo-server-manager
#./install.ps1
./build.ps1
```

#### Building the React frontend
Frontend is built using React and the AdminLTE CSS framework. See app/dist/ for AdminLTE included files and license.

The root of the UI application is served at app/index.html.  Run the npm build script and the Go application during development to get live rebuilding of the UI code.

All necessary CSS and JavaScript files are included for running the UI.

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
