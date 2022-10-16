$ErrorActionPreference = "Stop"
Try {
	del Robintris64.exe
} 
Catch {
	Write-Host "An error occurred del 64"
}
Finally {
    Write-Host "Delete of Robintris64.exe finished. Building new one."
}

$Env:GOARCH = 'amd64'
go build -o Robintris64.exe -tags=ebitenginesinglethread main.go
Remove-Item Env:GOARCH

Try {
	del Robintris32.exe
} 
Catch {
	Write-Host "An error occurred del 32"
}
Finally {
    Write-Host "Delete of Robintris32.exe finished. Building new one."
}

$Env:GOARCH = '386'
go build -o Robintris32.exe -tags=ebitenginesinglethread main.go
Remove-Item Env:GOARCH

Try {
	del RobintrisARM64.exe
} 
Catch {
	Write-Host "An error occurred del ARM64"
}
Finally {
    Write-Host "Delete of RobintrisARM64.exe finished. Building new one."
}

$Env:GOARCH = 'arm64'
go build -o RobintrisARM64.exe -tags=ebitenginesinglethread main.go
Remove-Item Env:GOARCH

Try {
	del Robintris.wasm
} 
Catch {
	Write-Host "An error occurred del wasm"
}
Finally {
    Write-Host "Delete of Robintris.wasm finished. Building new one."
}

$Env:GOOS = 'js'
$Env:GOARCH = 'wasm'
go build -o Robintris.wasm -tags=ebitenginesinglethread -ldflags="-s -w" main.go
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
