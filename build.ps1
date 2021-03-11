Set-Location ./gui

Write-Host "#: building js code..."
yarn build
Write-Host "#: generating js bin data..."
go-bindata -o ../host/bindata.go ./dist/...
Write-Host "#: done"

Set-Location ../host
go build -ldflags="-s -w"
#upx ./host.exe

Move-Item ./host.exe D:\Users\Lukiya\Desktop\spiders -Force
Set-Location ../