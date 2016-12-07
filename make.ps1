SET NODE_ENV=Production

rmdir -Recurse build
mkdir build
cd src
go build -o ../build/azure-csgo-server-windows.exe
cd ../
copy -R app/ build/
copy conf.json build/
copy -R configs/ build/
copy -R templates/ build/
copy -R certs/ build/