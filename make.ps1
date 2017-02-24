SET NODE_ENV=Production

rmdir -Recurse -Force build

mkdir build
cd ui
npm run build
cd ../src
glide install
go build -o ../build/azure-csgo-server.exe
cd ../

copy -R app/ build/
copy conf.json build/
copy -R configs/ build/
copy -R templates/ build/
copy -R certs/ build/
copy web.config build/