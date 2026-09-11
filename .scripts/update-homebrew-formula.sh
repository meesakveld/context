#!/usr/bin/env bash

set -euo pipefail

OWNER="meesakveld"
PROJECT="context"
TAP_REPOSITORY="homebrew-tap"
FORMULA_NAME="context"

VERSION="${1:-}"

if [[ -z "$VERSION" ]]; then
  echo "Usage: $0 <version>"
  exit 1
fi

if [[ "$VERSION" == v* ]]; then
  VERSION="${VERSION#v}"
fi

TAG="v${VERSION}"

BASE_URL="https://github.com/${OWNER}/${PROJECT}/releases/download/${TAG}"

CHECKSUMS_URL="${BASE_URL}/checksums.txt"

echo "Updating Homebrew formula for ${PROJECT} ${VERSION}"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

curl -fsSL "$CHECKSUMS_URL" -o "$TMP_DIR/checksums.txt"

get_checksum() {
  local filename="$1"

  awk -v file="$filename" '$2 == file { print $1 }' "$TMP_DIR/checksums.txt"
}

DARWIN_ARM64_SHA="$(get_checksum "${PROJECT}_darwin_arm64.tar.gz")"
DARWIN_AMD64_SHA="$(get_checksum "${PROJECT}_darwin_amd64.tar.gz")"
LINUX_ARM64_SHA="$(get_checksum "${PROJECT}_linux_arm64.tar.gz")"
LINUX_AMD64_SHA="$(get_checksum "${PROJECT}_linux_amd64.tar.gz")"

for checksum in \
  "$DARWIN_ARM64_SHA" \
  "$DARWIN_AMD64_SHA" \
  "$LINUX_ARM64_SHA" \
  "$LINUX_AMD64_SHA"
do
  if [[ -z "$checksum" ]]; then
    echo "Error: Could not find all required checksums."
    exit 1
  fi
done

FORMULA_FILE="$TMP_DIR/${FORMULA_NAME}.rb"

cat > "$FORMULA_FILE" <<EOF
class Context < Formula
  desc "CLI for turning codebases into structured, AI-ready context"
  homepage "https://github.com/${OWNER}/${PROJECT}"
  version "${VERSION}"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "${BASE_URL}/${PROJECT}_darwin_arm64.tar.gz"
      sha256 "${DARWIN_ARM64_SHA}"
    else
      url "${BASE_URL}/${PROJECT}_darwin_amd64.tar.gz"
      sha256 "${DARWIN_AMD64_SHA}"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "${BASE_URL}/${PROJECT}_linux_arm64.tar.gz"
      sha256 "${LINUX_ARM64_SHA}"
    else
      url "${BASE_URL}/${PROJECT}_linux_amd64.tar.gz"
      sha256 "${LINUX_AMD64_SHA}"
    end
  end

  def install
    bin.install "context"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/context --version")
  end
end
EOF

echo "Generated formula:"
cat "$FORMULA_FILE"

if [[ -z "${TAP_GITHUB_TOKEN:-}" ]]; then
  echo "Error: TAP_GITHUB_TOKEN is not set."
  exit 1
fi

TAP_DIR="$TMP_DIR/${TAP_REPOSITORY}"

git clone \
  "https://x-access-token:${TAP_GITHUB_TOKEN}@github.com/${OWNER}/${TAP_REPOSITORY}.git" \
  "$TAP_DIR"

mkdir -p "$TAP_DIR/Formula"

cp "$FORMULA_FILE" "$TAP_DIR/Formula/${FORMULA_NAME}.rb"

cd "$TAP_DIR"

git config user.name "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"

git add "Formula/${FORMULA_NAME}.rb"

if git diff --cached --quiet; then
  echo "Formula is already up to date."
  exit 0
fi

git commit -m "Update ${PROJECT} to ${TAG}"
git push origin main

echo "Homebrew formula updated successfully."