#!/bin/sh
# Install one prebuilt platform archive. No language runtime or package manager.
set -eu
version=latest
setup_runtime=false
install_dir=${FUNGI_INSTALL_DIR:-"$HOME/.fungi"}
base=${FUNGI_RELEASE_BASE_URL:-https://github.com/rengwu/fungi/releases}
while [ "$#" -gt 0 ]; do
  case "$1" in
    --setup-runtime) setup_runtime=true; shift;;
    --version|--dir)
      [ "$#" -ge 2 ] || { echo "Missing value for $1" >&2; exit 2; }
      case "$1" in --version) version=$2;; --dir) install_dir=$2;; esac
      shift 2;;
    *) echo "Usage: install.sh [--version vX.Y.Z] [--dir DIRECTORY] [--setup-runtime]" >&2; exit 2;;
  esac
done
case "$version" in *[!a-zA-Z0-9.-]*|'') echo "Invalid release version" >&2; exit 2;; esac
case "$(uname -s)" in Darwin) platform=darwin;; Linux) platform=linux;; *) echo "Use install.ps1 on Windows" >&2; exit 1;; esac
case "$(uname -m)" in arm64|aarch64) arch=arm64;; x86_64|amd64) arch=amd64;; *) echo "Unsupported CPU" >&2; exit 1;; esac
case "$platform-$arch" in darwin-arm64|linux-amd64) ;; *) echo "No published package for $platform-$arch; this alpha supports Mac arm64, Linux x64 and Windows x64." >&2; exit 1;; esac
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
runtime_status=0
FUNGI_RENDERER_DIR= "$temp/package/fungi" doctor > "$temp/runtime.json" || runtime_status=$?
if [ "$runtime_status" -ne 0 ]; then
  cat "$temp/runtime.json" >&2
  [ "$runtime_status" -eq 3 ] || exit "$runtime_status"
  [ "$setup_runtime" = true ] || { echo "Runtime missing. Re-run with --setup-runtime to install it explicitly." >&2; exit 3; }
  [ "$platform" = linux ] || { echo "The macOS webview is supplied by macOS; update or repair the OS." >&2; exit 3; }
  # Support the verified Ubuntu releases; do not guess other package managers.
  [ -r /etc/os-release ] || { echo "Install WebKitGTK 4.1 using your distribution's instructions." >&2; exit 3; }
  . /etc/os-release
  case "${ID:-}:${VERSION_ID:-}" in
    ubuntu:22.04|ubuntu:24.04) ;;
    *) echo "Automatic runtime setup supports Ubuntu 22.04/24.04 only. Install WebKitGTK 4.1 using your distribution's instructions." >&2; exit 3;;
  esac
  echo "Installing WebKitGTK 4.1 from configured Ubuntu repositories; administrator permission may be requested." >&2
  if [ "$(id -u)" -eq 0 ]; then
    apt-get update >&2
    apt-get install -y libwebkit2gtk-4.1-0 >&2
  else
    sudo apt-get update >&2
    sudo apt-get install -y libwebkit2gtk-4.1-0 >&2
  fi
  FUNGI_RENDERER_DIR= "$temp/package/fungi" doctor >&2
fi
"$temp/package/fungi" install --archive "$temp/$name" --sha256 "$expected" --dir "$install_dir"
printf '\nFungi is installed. Add this directory to PATH if needed:\n  %s/bin\n' "$install_dir" >&2
case "$version" in *-*) echo 'For alpha updates: fungi update --version vX.Y.Z-alpha.N' >&2;; *) echo 'Update with: fungi update' >&2;; esac
