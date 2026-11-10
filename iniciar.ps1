$ErrorActionPreference = 'Stop'
Set-Location -LiteralPath $PSScriptRoot
$portableGo = Join-Path $PSScriptRoot '..\..\work\toolchain\go\bin\go.exe'
if (Get-Command go -ErrorAction SilentlyContinue) { & go run -buildvcs=false ./cmd/api } elseif (Test-Path -LiteralPath $portableGo) { & $portableGo run -buildvcs=false ./cmd/api } else { throw 'Instala Go 1.22 o posterior para ejecutar el proyecto.' }
