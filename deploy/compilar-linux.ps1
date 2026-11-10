# Compila Minimarket SaaS (React + Go) para un servidor EC2 Linux (x86_64) y deja todo listo en la carpeta dist\
$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSScriptRoot
Set-Location -LiteralPath $root

$portableGo = Join-Path $root '..\..\work\toolchain\go\bin\go.exe'
if (Get-Command go -ErrorAction SilentlyContinue) { $go = 'go' }
elseif (Test-Path -LiteralPath $portableGo) { $go = (Resolve-Path $portableGo).Path }
else { throw 'Instala Go 1.22 o posterior para compilar el proyecto.' }

if ((Test-Path 'frontend\package.json') -and (Get-Command npm -ErrorAction SilentlyContinue)) {
    if (-not (Test-Path 'frontend\node_modules')) {
        Write-Host 'Instalando dependencias de React...'
        npm --prefix frontend install
    }
    Write-Host 'Compilando frontend React (Vite) en internal/presentation/web...'
    npm --prefix frontend run build
    if ($LASTEXITCODE -ne 0) { throw 'La compilación del frontend React falló.' }
}

$dist = Join-Path $root 'dist'
New-Item -ItemType Directory -Force $dist | Out-Null

Write-Host 'Ejecutando pruebas automatizadas de Go...'
& $go test ./...
if ($LASTEXITCODE -ne 0) { throw 'Las pruebas fallaron; no se genera el paquete.' }

Write-Host 'Compilando binario único para Linux amd64 (AWS EC2)...'
$env:GOOS = 'linux'; $env:GOARCH = 'amd64'; $env:CGO_ENABLED = '0'
try {
    & $go build -buildvcs=false -trimpath -ldflags '-s -w' -o (Join-Path $dist 'minimarket') ./cmd/api
    if ($LASTEXITCODE -ne 0) { throw 'La compilación falló.' }
} finally {
    Remove-Item Env:GOOS, Env:GOARCH, Env:CGO_ENABLED -ErrorAction SilentlyContinue
}

$serverFiles = 'instalar.sh', 'minimarket.service', 'instalar-https.sh', 'caddy.service', 'Caddyfile'
foreach ($f in $serverFiles) { Copy-Item -Force (Join-Path $PSScriptRoot $f) $dist }
# Copia los datos actuales solo si existen (el instalador no sobrescribe datos ya presentes en el servidor)
if (Test-Path 'data.json') { Copy-Item -Force 'data.json' $dist }

# Normaliza los saltos de línea a LF para que bash y systemd los lean bien en Linux
foreach ($f in $serverFiles) {
    $p = Join-Path $dist $f
    $t = [IO.File]::ReadAllText($p) -replace "`r`n", "`n"
    [IO.File]::WriteAllText($p, $t, (New-Object Text.UTF8Encoding $false))
}

Write-Host "`nListo. Archivos para subir a AWS EC2 en: $dist"
Get-ChildItem $dist | Format-Table Name, Length
