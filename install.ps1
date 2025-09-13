# nos Installation Script for Windows
# This script downloads and installs nos while handling Windows security warnings

param(
    [string]$Version = "latest",
    [string]$InstallPath = "$env:USERPROFILE\bin"
)

Write-Host "🚀 Installing nos - Nostr command-line client" -ForegroundColor Green
Write-Host ""

# Function to get latest version from GitHub
function Get-LatestVersion {
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/PlebOne/nos/releases/latest"
        return $response.tag_name
    }
    catch {
        Write-Host "❌ Failed to get latest version. Using fallback." -ForegroundColor Red
        return "v1.1.4"
    }
}

# Determine version to install
if ($Version -eq "latest") {
    $Version = Get-LatestVersion
    Write-Host "📦 Latest version: $Version" -ForegroundColor Blue
}

# Determine architecture
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
$filename = "nos-$($Version.TrimStart('v'))-windows-$arch.exe"
$url = "https://github.com/PlebOne/nos/releases/download/$Version/$filename"

Write-Host "🌐 Downloading from: $url" -ForegroundColor Blue
Write-Host "📁 Installing to: $InstallPath" -ForegroundColor Blue
Write-Host ""

# Create install directory
if (!(Test-Path $InstallPath)) {
    New-Item -ItemType Directory -Path $InstallPath -Force | Out-Null
    Write-Host "📁 Created directory: $InstallPath" -ForegroundColor Green
}

# Download file
$outputPath = Join-Path $InstallPath "nos.exe"
try {
    Write-Host "⬇️  Downloading nos..." -ForegroundColor Yellow
    
    # Add TLS security protocol support
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    
    Invoke-WebRequest -Uri $url -OutFile $outputPath -UseBasicParsing
    Write-Host "✅ Downloaded successfully" -ForegroundColor Green
}
catch {
    Write-Host "❌ Download failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Unblock the file to prevent Windows security warnings
try {
    Unblock-File -Path $outputPath
    Write-Host "🔓 Unblocked file for Windows security" -ForegroundColor Green
}
catch {
    Write-Host "⚠️  Could not unblock file automatically. You may need to manually approve execution." -ForegroundColor Yellow
}

# Add to PATH if not already there
$currentPath = [Environment]::GetEnvironmentVariable("PATH", [EnvironmentVariableTarget]::User)
if ($currentPath -notlike "*$InstallPath*") {
    Write-Host "🔧 Adding to PATH..." -ForegroundColor Yellow
    $newPath = "$currentPath;$InstallPath"
    [Environment]::SetEnvironmentVariable("PATH", $newPath, [EnvironmentVariableTarget]::User)
    Write-Host "✅ Added $InstallPath to PATH" -ForegroundColor Green
    Write-Host "⚠️  Please restart your terminal or run: `$env:PATH += ';$InstallPath'" -ForegroundColor Yellow
}

# Test installation
Write-Host ""
Write-Host "🧪 Testing installation..." -ForegroundColor Blue
try {
    $env:PATH += ";$InstallPath"
    $version = & "$outputPath" --version 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Installation successful!" -ForegroundColor Green
    } else {
        # Try running without version flag (nos might not have --version)
        Write-Host "✅ nos installed at: $outputPath" -ForegroundColor Green
    }
}
catch {
    Write-Host "✅ nos installed at: $outputPath" -ForegroundColor Green
}

Write-Host ""
Write-Host "🎉 Installation complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Usage examples:" -ForegroundColor Cyan
Write-Host "  nos 'Hello Nostr world!'" -ForegroundColor White
Write-Host "  echo 'Check out #bitcoin' | nos" -ForegroundColor White
Write-Host ""
Write-Host "For more information: https://github.com/PlebOne/nos" -ForegroundColor Blue