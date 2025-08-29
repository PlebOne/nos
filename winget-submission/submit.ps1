# Winget Package Submission Helper Script
# This script helps submit nos to the Microsoft winget repository

Write-Host "🚀 nos Winget Package Submission Helper" -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Green
Write-Host

# Check if winget-pkgs directory exists
if (-not (Test-Path "../../winget-pkgs")) {
    Write-Host "📥 Microsoft winget-pkgs repository not found!" -ForegroundColor Yellow
    Write-Host "Please run these commands first:" -ForegroundColor Yellow
    Write-Host
    Write-Host "1. Fork https://github.com/microsoft/winget-pkgs" -ForegroundColor Cyan
    Write-Host "2. Clone your fork:" -ForegroundColor Cyan
    Write-Host "   cd .." -ForegroundColor White
    Write-Host "   git clone https://github.com/YOUR_USERNAME/winget-pkgs.git" -ForegroundColor White
    Write-Host "   cd nos" -ForegroundColor White
    Write-Host
    Write-Host "Then run this script again." -ForegroundColor Yellow
    exit 1
}

Write-Host "📁 Copying manifest files to winget-pkgs repository..." -ForegroundColor Blue

# Create directory structure
New-Item -Path "../../winget-pkgs/manifests/p/PlebOne/nos/1.1.1" -ItemType Directory -Force | Out-Null

# Copy manifest files
Copy-Item "winget-submission/manifests/p/PlebOne/nos/1.1.1/*" "../../winget-pkgs/manifests/p/PlebOne/nos/1.1.1/"

Write-Host "✅ Files copied successfully!" -ForegroundColor Green
Write-Host

Set-Location "../../winget-pkgs"

Write-Host "🌿 Creating branch for submission..." -ForegroundColor Blue
git checkout -b add-plebOne-nos-1.1.1

Write-Host "📝 Adding manifest files..." -ForegroundColor Blue
git add manifests/p/PlebOne/nos/1.1.1/

Write-Host "💾 Committing changes..." -ForegroundColor Blue
$commitMessage = @"
Add PlebOne.nos version 1.1.1

- Package: nos - A beautiful command-line client for posting to Nostr
- Publisher: PlebOne
- Version: 1.1.1
- License: MIT
- Architectures: x64, ARM64
- Installer Type: Portable (zip)

Release URL: https://github.com/PlebOne/nos/releases/tag/v1.1.1
"@

git commit -m $commitMessage

Write-Host "🚀 Pushing to your fork..." -ForegroundColor Blue
git push origin add-plebOne-nos-1.1.1

Write-Host
Write-Host "🎉 Success! Next steps:" -ForegroundColor Green
Write-Host "========================" -ForegroundColor Green
Write-Host "1. Go to your winget-pkgs fork on GitHub" -ForegroundColor Cyan
Write-Host "2. Click 'Compare & pull request'" -ForegroundColor Cyan
Write-Host "3. Use title: Add PlebOne.nos version 1.1.1" -ForegroundColor Cyan
Write-Host "4. Copy the PR description from winget-submission/README.md" -ForegroundColor Cyan
Write-Host "5. Submit the pull request" -ForegroundColor Cyan
Write-Host
Write-Host "Expected approval time: 1-7 days" -ForegroundColor Yellow
Write-Host "Once approved, users can install with: winget install PlebOne.nos" -ForegroundColor Yellow
Write-Host

Set-Location "../../nos"
