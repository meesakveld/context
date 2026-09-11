# Deployment Guide

This guide explains how to release a new version of `context` and publish it through GitHub Releases and Homebrew.

## Overview

The release flow is fully automated:

```text
Local changes
     │
     ▼
Commit + push
     │
     ▼
Create version tag
     │
     ▼
GitHub Actions
     │
     ├── GoReleaser
     │     ├── Build binaries
     │     ├── Generate checksums
     │     └── Create GitHub Release
     │
     └── Homebrew update
           │
           ▼
     meesakveld/homebrew-tap
           │
           ▼
     brew install / upgrade
```

## 1. Make sure the project is clean

Before creating a release, run:

```bash
go test ./...
go vet ./...
go build ./...
```

All three commands should complete without errors.

Check the Git working tree:

```bash
git status
```

Make sure there are no unintended changes.

---

## 2. Commit and push your changes

Stage the changes:

```bash
git add .
```

Create a commit:

```bash
git commit -m "Describe the changes"
```

Push to `main`:

```bash
git push origin main
```

---

## 3. Choose the new version

The project uses semantic versioning:

```text
MAJOR.MINOR.PATCH
```

Examples:

```text
v1.0.0
v1.1.0
v1.1.1
v1.2.0
v2.0.0
```

Use:

* `PATCH` for bug fixes and small corrections
* `MINOR` for backwards-compatible features
* `MAJOR` for breaking changes

### Example

If the current version is:

```text
v1.1.0
```

and you only improved the CLI help output:

```text
v1.1.1
```

---

## 4. Create and push the tag

Create the tag locally:

```bash
git tag v1.1.1
```

Push the tag:

```bash
git push origin v1.1.1
```

This automatically triggers the GitHub Actions release workflow.

You can check the tags with:

```bash
git tag
```

---

## 5. GitHub Actions

The workflow is located at:

```text
.github/workflows/release.yml
```

It runs whenever a tag matching this pattern is pushed:

```yaml
tags:
  - "v*"
```

The workflow performs two main tasks.

### GoReleaser

GoReleaser builds the following platforms:

```text
macOS ARM64
macOS AMD64
Linux ARM64
Linux AMD64
Windows AMD64
```

It then creates:

```text
*.tar.gz
*.zip
checksums.txt
```

and publishes them to the GitHub Release.

### Homebrew

After GoReleaser succeeds, the workflow runs:

```text
.scripts/update-homebrew-formula.sh
```

This script:

1. Downloads `checksums.txt`
2. Extracts the checksums for macOS and Linux
3. Generates `Formula/context.rb`
4. Clones `meesakveld/homebrew-tap`
5. Updates the formula
6. Commits the change
7. Pushes it to the tap

---

## 6. Wait for the release workflow

After pushing the tag, open the GitHub Actions page for the repository.

The workflow should show:

```text
Release
├── Checkout
├── Set up Go
├── Create GitHub App token
├── Release
└── Update Homebrew formula
```

All steps should be green.

If the `Release` step succeeds but `Update Homebrew formula` fails, the GitHub Release may already exist. Do not immediately delete or recreate the tag.

First investigate the workflow error.

---

## 7. Verify the GitHub Release

Open the repository's Releases page and verify that the new version exists.

For example:

```text
v1.1.1
```

Check that the release contains:

```text
context_darwin_arm64.tar.gz
context_darwin_amd64.tar.gz
context_linux_arm64.tar.gz
context_linux_amd64.tar.gz
context_windows_amd64.zip
checksums.txt
```

---

## 8. Verify the Homebrew formula

The formula should automatically be updated in:

```text
meesakveld/homebrew-tap
```

Specifically:

```text
Formula/context.rb
```

The formula should contain the new version:

```ruby
version "1.1.1"
```

and the correct SHA-256 checksums.

---

## 9. Test the Homebrew installation

Update Homebrew:

```bash
brew update
```

If `context` is already installed:

```bash
brew upgrade context
```

If it is not installed:

```bash
brew install meesakveld/tap/context
```

Then verify:

```bash
context --version
```

Expected:

```text
context v1.1.1
```

Also test:

```bash
context --help
```

and ideally:

```bash
context --clipboard
```

---

# Release Checklist

Use this checklist for every release:

```text
[ ] Code is finished
[ ] go test ./...
[ ] go vet ./...
[ ] go build ./...
[ ] git status checked
[ ] Changes committed
[ ] Changes pushed to main
[ ] Correct semantic version selected
[ ] Version tag created
[ ] Version tag pushed
[ ] GitHub Actions release succeeded
[ ] GitHub Release exists
[ ] All binaries are present
[ ] checksums.txt is present
[ ] Homebrew formula updated
[ ] brew update
[ ] brew upgrade/install context
[ ] context --version
[ ] context --help
```

---

# Important: Do Not Reuse Release Tags

Once a version has been released, do not reuse the same tag.

For example, if:

```text
v1.1.0
```

already exists, do not try to recreate `v1.1.0`.

Instead create:

```text
v1.1.1
```

GitHub releases and GoReleaser can treat existing releases/tags as immutable.

---

# Current Release Setup

The project uses:

```text
Repository:
github.com/meesakveld/context

Homebrew tap:
github.com/meesakveld/homebrew-tap

Binary:
context

Module:
github.com/meesakveld/context

Release workflow:
.github/workflows/release.yml

GoReleaser:
.goreleaser.yaml

Homebrew update script:
.scripts/update-homebrew-formula.sh
```

Users install the stable release with:

```bash
brew install meesakveld/tap/context
```

and upgrade with:

```bash
brew upgrade context
```

---

# Emergency Troubleshooting

## GoReleaser succeeds but Homebrew fails

The GitHub Release may already have been created.

Do **not** immediately create the same tag again.

Check the error in:

```text
Update Homebrew formula
```

Common causes include:

* missing GitHub App credentials
* incorrect repository permissions
* Homebrew tap unavailable
* checksum file unavailable
* script missing from the tagged commit

---

## The workflow cannot find the release script

Verify that the script exists on the commit referenced by the tag:

```bash
git show v1.1.1:.scripts/update-homebrew-formula.sh
```

If this returns the file contents, the tag contains the script.

---

## Homebrew installs the old version

First run:

```bash
brew update
```

Then:

```bash
brew upgrade context
```

You can check which formula Homebrew is using with:

```bash
brew info context
```

---

# Quick Release

Once everything is configured and tested, the normal release process is simply:

```bash
go test ./...
go vet ./...
go build ./...

git add .
git commit -m "Your changes"
git push origin main

git tag v1.1.1
git push origin v1.1.1
```

Then wait for GitHub Actions and verify:

```bash
brew update
brew upgrade context
context --version
```

That's it.
