param(
    [Parameter(Mandatory=$true)][string]$Version,
    [Parameter(Mandatory=$true)][string]$InstallerUrlX64,
    [Parameter(Mandatory=$true)][string]$Sha256X64,
    [string]$InstallerUrlArm64 = "",
    [string]$Sha256Arm64 = ""
)

$templatePath = Join-Path $PSScriptRoot 'nos-windows-manifest-template.yaml'
$outDir = Join-Path $PSScriptRoot ("PlebOne\nos\$Version")
if (!(Test-Path $outDir)) { New-Item -ItemType Directory -Path $outDir -Force | Out-Null }

$template = Get-Content $templatePath -Raw

$replaced = $template -replace '{{VERSION}}', $Version -replace '{{INSTALLER_URL}}', $InstallerUrlX64 -replace '{{INSTALLER_SHA256}}', $Sha256X64 -replace '{{INSTALLER_URL_ARM64}}', $InstallerUrlArm64 -replace '{{INSTALLER_SHA256_ARM64}}', $Sha256Arm64

$outFile = Join-Path $outDir 'manifest.yaml'
$replaced | Out-File -FilePath $outFile -Encoding utf8

Write-Host "Generated winget manifest at: $outFile"
