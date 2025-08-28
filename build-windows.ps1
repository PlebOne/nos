param(
    [string]$Version = "0.9",
    [string]$Arch = "amd64",
    [string]$OutDir = "dist"
)

Write-Host "Building nos for Windows..."

$BuildDir = Join-Path -Path "build" -ChildPath "nos_${Version}_${Arch}_windows"
if (Test-Path $BuildDir) { Remove-Item -Recurse -Force $BuildDir }
New-Item -ItemType Directory -Force -Path $BuildDir | Out-Null

$BinDir = Join-Path $BuildDir "bin"
New-Item -ItemType Directory -Force -Path $BinDir | Out-Null

$exe = "nos.exe"
$outPath = Join-Path $BinDir $exe

Write-Host "Go building -> $outPath"
& go build -ldflags "-s -w" -o $outPath .

if ($LASTEXITCODE -ne 0) {
    Write-Error "go build failed with exit code $LASTEXITCODE"
    exit $LASTEXITCODE
}

# Copy docs and image if present
if (Test-Path "README.md") { Copy-Item "README.md" -Destination $BuildDir -Force }
if (Test-Path "nos.jpeg") { Copy-Item "nos.jpeg" -Destination (Join-Path $BuildDir "nos.jpeg") -Force }

# Create output dir
if (!(Test-Path $OutDir)) { New-Item -ItemType Directory -Path $OutDir | Out-Null }

$zipName = "nos_${Version}_${Arch}_windows.zip"
$zipPath = Join-Path $OutDir $zipName

if (Test-Path $zipPath) { Remove-Item $zipPath -Force }

Write-Host "Creating ZIP package -> $zipPath"
Compress-Archive -Path (Join-Path $BuildDir "*") -DestinationPath $zipPath -Force

Write-Host "Package created: $zipPath"
