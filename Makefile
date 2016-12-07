# Build tool for Azure CS:GO Server Manager
#

NODE_ENV:=production

build:
	# Build Linux release
	mkdir build
	GOOS=linux GOARCH=amd64 go build -o azure-csgo-server-linux/azure-csgo-server-windows src/*
#	ui/node_modules/webpack/bin/webpack.js ui/webpack.config.js app/bundle.js --progress --profile --colors 
	cp -r app/ azure-csgo-server-linux/
	cp conf.json.example azure-csgo-server-linux/conf.json
	zip -r build/azure-csgo-server-linux-x64.zip azure-csgo-server-linux
	rm -rf azure-csgo-server-linux
	# Build Windows release
	GOOS=windows GOARCH=386 go build -o azure-csgo-server-windows/azure-csgo-server-windows.exe src/*
	cp -r app/ azure-csgo-server-windows/
	cp conf.json.example azure-csgo-server-windows/conf.json
    cp -r configs/ azure-csgo-server-windows
    cp -r templates/ azure-csgo-server-windows
	zip -r build/factorio-server-manager-windows.zip azure-csgo-server-windows
	rm -rf azure-csgo-server-windows

