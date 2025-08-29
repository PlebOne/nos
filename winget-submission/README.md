# Winget Package Submission to Microsoft

This directory contains the manifest files for submitting **nos** to the official Microsoft winget repository.

## 📋 Submission Steps

### 1. Fork the Microsoft winget repository
```bash
# Go to: https://github.com/microsoft/winget-pkgs
# Click "Fork" to create your fork
```

### 2. Clone your fork locally
```bash
git clone https://github.com/YOUR_USERNAME/winget-pkgs.git
cd winget-pkgs
```

### 3. Copy the manifest files
Copy the files from `winget-submission/manifests/p/PlebOne/nos/1.1.1/` to:
```
winget-pkgs/manifests/p/PlebOne/nos/1.1.1/
```

### 4. Create a pull request
```bash
# Create a new branch
git checkout -b add-plebOne-nos-1.1.1

# Add the files
git add manifests/p/PlebOne/nos/1.1.1/

# Commit with proper message
git commit -m "Add PlebOne.nos version 1.1.1"

# Push to your fork
git push origin add-plebOne-nos-1.1.1
```

### 5. Submit the PR
- Go to your fork on GitHub
- Click "Compare & pull request"
- Use title: **Add PlebOne.nos version 1.1.1**
- Include description below

## 📝 PR Description Template

```markdown
# Add PlebOne.nos version 1.1.1

## Package Information
- **Package**: nos
- **Publisher**: PlebOne  
- **Version**: 1.1.1
- **License**: MIT
- **Description**: A beautiful command-line client for posting to Nostr

## Validation
✅ Package identifier follows naming convention (PlebOne.nos)
✅ All required manifest files included (installer, locale, version)
✅ URLs and checksums verified from GitHub release
✅ Package builds and runs successfully on Windows

## Release Details
- **Release URL**: https://github.com/PlebOne/nos/releases/tag/v1.1.1
- **Architectures**: x64, ARM64
- **Installer Type**: Portable (zip)
- **File Format**: tar.gz archives

## Testing
Tested installation with:
```bash
winget install --manifest manifests/p/PlebOne/nos/1.1.1/
```

Package installs successfully and `nos` command is available.
```

## 📦 Files Included

1. **PlebOne.nos.yaml** - Version manifest
2. **PlebOne.nos.installer.yaml** - Installer details with checksums
3. **PlebOne.nos.locale.en-US.yaml** - Package metadata and description

## ⏱️ Expected Timeline

- **Automated validation**: Immediate (when PR is created)
- **Human review**: 1-3 days
- **Final approval**: 1-7 days total
- **Package availability**: Within 24 hours of approval

## 🔍 Validation Checklist

Before submitting, ensure:
- [ ] All URLs are accessible
- [ ] SHA256 checksums are correct
- [ ] Package identifier is unique
- [ ] Description is clear and helpful
- [ ] License information is accurate
- [ ] Release notes URL is valid

## 📞 Support

If you encounter issues during submission:
- Review the [winget submission guidelines](https://docs.microsoft.com/en-us/windows/package-manager/package/)
- Check existing issues in the winget-pkgs repository
- Join the Windows Package Manager community discussions

---

**Next Step**: Follow the submission steps above to get **nos** available via `winget install PlebOne.nos` 🚀
