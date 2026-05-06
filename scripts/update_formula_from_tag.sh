#!/usr/bin/env bash
set -euo pipefail

TAG="${1:-}"
if [[ -z "${TAG}" ]]; then
  echo "usage: $0 <tag> (example: v0.1.0)" >&2
  exit 2
fi

REPO="wangweicheng7/devclean-cli"
FORMULA_FILE="homebrew-tap/Formula/devclean-cli.rb"
VERSION="${TAG#v}"
CHECKSUMS_URL="https://github.com/${REPO}/releases/download/${TAG}/checksums.txt"
ARM_URL="https://github.com/${REPO}/releases/download/${TAG}/devclean_${VERSION}_darwin_arm64.tar.gz"
AMD_URL="https://github.com/${REPO}/releases/download/${TAG}/devclean_${VERSION}_darwin_amd64.tar.gz"

TMP_CHECKSUMS="$(mktemp -t devclean-checksums.XXXXXX.txt)"
cleanup() { rm -f "${TMP_CHECKSUMS}"; }
trap cleanup EXIT

echo "downloading checksums: ${CHECKSUMS_URL}" >&2
curl -L -o "${TMP_CHECKSUMS}" "${CHECKSUMS_URL}"

ARM_SHA="$(awk '/devclean_'"${VERSION}"'_darwin_arm64\.tar\.gz$/ {print $1}' "${TMP_CHECKSUMS}")"
AMD_SHA="$(awk '/devclean_'"${VERSION}"'_darwin_amd64\.tar\.gz$/ {print $1}' "${TMP_CHECKSUMS}")"

if [[ -z "${ARM_SHA}" || -z "${AMD_SHA}" ]]; then
  echo "failed to parse checksums.txt for version ${VERSION}" >&2
  exit 1
fi

cat > "${FORMULA_FILE}" <<EOF
class DevcleanCli < Formula
  desc "macOS developer cleanup CLI (safe-first)"
  homepage "https://github.com/${REPO}"
  version "${VERSION}"

  on_macos do
    if Hardware::CPU.arm?
      url "${ARM_URL}"
      sha256 "${ARM_SHA}"
    else
      url "${AMD_URL}"
      sha256 "${AMD_SHA}"
    end
  end

  def install
    bin.install "devclean"
  end

  test do
    system "#{bin}/devclean", "doctor"
  end
end
EOF

echo "updated ${FORMULA_FILE}" >&2
echo "version: ${VERSION}" >&2
echo "arm64 sha256: ${ARM_SHA}" >&2
echo "amd64 sha256: ${AMD_SHA}" >&2

