#!/bin/sh
# Install one prebuilt platform archive. No language runtime or package manager.
set -eu
version=latest
install_dir=${FUNGI_INSTALL_DIR:-"$HOME/.fungi"}
base=${FUNGI_RELEASE_BASE_URL:-https://github.com/rengwu/fungi/releases}
while [ "$#" -gt 0 ]; do
  case "$1" in
    --version|--dir)
      [ "$#" -ge 2 ] || { echo "Missing value for $1" >&2; exit 2; }
      case "$1" in --version) version=$2;; --dir) install_dir=$2;; esac
      shift 2;;
    *) echo "Usage: install.sh [--version vX.Y.Z] [--dir DIRECTORY]" >&2; exit 2;;
  esac
done
case "$version" in *[!a-zA-Z0-9.-]*|'') echo "Invalid release version" >&2; exit 2;; esac
case "$(uname -s)" in Darwin) platform=darwin;; Linux) platform=linux;; *) echo "Use install.ps1 on Windows" >&2; exit 1;; esac
case "$(uname -m)" in arm64|aarch64) arch=arm64;; x86_64|amd64) arch=amd64;; *) echo "Unsupported CPU" >&2; exit 1;; esac
case "$platform-$arch" in darwin-arm64|linux-amd64) ;; *) echo "No published package for $platform-$arch; published targets are Mac arm64, Linux x64 and Windows x64." >&2; exit 1;; esac
download() {
  if command -v curl >/dev/null 2>&1; then curl --fail --location --silent --show-error --connect-timeout 20 --max-time 900 "$1" --output "$2"
  elif command -v wget >/dev/null 2>&1; then wget --quiet --timeout=20 --tries=1 "$1" -O "$2"
  else echo "curl or wget is required to download Fungi" >&2; return 1; fi
}
if command -v sha256sum >/dev/null 2>&1; then hash_command=sha256sum
elif command -v shasum >/dev/null 2>&1; then hash_command=shasum
else echo "sha256sum or shasum is required to verify the download" >&2; exit 1; fi
temp=$(mktemp -d)
trap 'rm -rf "$temp"' EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM
name="fungi-$platform-$arch.tar.gz"
if [ "$version" = latest ]; then location="$base/latest/download"; else location="$base/download/$version"; fi
download "$location/$name.sha256" "$temp/checksum"
read -r expected checksum_name < "$temp/checksum"
case "$expected" in *[!a-fA-F0-9]*|'') echo "Invalid SHA-256" >&2; exit 1;; esac
[ "${#expected}" -eq 64 ] && [ "$checksum_name" = "$name" ] || { echo "Invalid release checksum record" >&2; exit 1; }
download "$location/$name" "$temp/$name"
if [ "$hash_command" = sha256sum ]; then actual=$(sha256sum "$temp/$name"); else actual=$(shasum -a 256 "$temp/$name"); fi
actual=${actual%% *}
[ "$actual" = "$expected" ] || { echo "Checksum mismatch; installation unchanged" >&2; exit 1; }
mkdir "$temp/package"
tar -xzf "$temp/$name" -C "$temp/package"
"$temp/package/fungi" install --archive "$temp/$name" --sha256 "$expected" --dir "$install_dir"
printf '\nFungi is installed. Add this directory to PATH if needed:\n  %s/bin\n' "$install_dir" >&2
case "$version" in *-*) echo 'For alpha updates: fungi update --version vX.Y.Z-alpha.N' >&2;; *) echo 'Update with: fungi update' >&2;; esac
