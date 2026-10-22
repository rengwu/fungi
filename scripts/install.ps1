param(
    [string]$Version = 'latest',
    [string]$Directory = $(if ($env:FUNGI_INSTALL_DIR) { $env:FUNGI_INSTALL_DIR } else { Join-Path $HOME '.fungi' }),
    [switch]$NoPath
)
$ErrorActionPreference = 'Stop'
# Windows PowerShell's per-chunk progress rendering can dominate large downloads.
# Emit a few stage messages instead, also useful when an agent captures output.
$ProgressPreference = 'SilentlyContinue'
if ($Version -notmatch '^[A-Za-z0-9.-]+$') { throw 'Invalid release version' }
$base = if ($env:FUNGI_RELEASE_BASE_URL) { $env:FUNGI_RELEASE_BASE_URL } else { 'https://github.com/rengwu/fungi/releases' }
$nativeArch = if ($env:PROCESSOR_ARCHITEW6432) { $env:PROCESSOR_ARCHITEW6432 } else { $env:PROCESSOR_ARCHITECTURE }
$arch = switch ($nativeArch) { 'AMD64' { 'amd64' }; 'ARM64' { 'arm64' }; default { throw "Unsupported CPU: $nativeArch" } }
$name = "fungi-windows-$arch.zip"
$location = if ($Version -eq 'latest') { "$base/latest/download" } else { "$base/download/$Version" }
$temp = Join-Path ([IO.Path]::GetTempPath()) ('fungi-install-' + [guid]::NewGuid().ToString('N'))
[IO.Directory]::CreateDirectory($temp) | Out-Null
try {
    # Windows PowerShell 5.1: no legacy HTML parsing or browser dependency.
    [Net.ServicePointManager]::SecurityProtocol = [Net.ServicePointManager]::SecurityProtocol -bor [Net.SecurityProtocolType]::Tls12
    $checksum = Join-Path $temp 'checksum'
    [Console]::Error.WriteLine("Downloading $name and its checksum...")
    Invoke-WebRequest -UseBasicParsing -Uri "$location/$name.sha256" -OutFile $checksum -TimeoutSec 900
    $fields = (Get-Content -Raw $checksum).Trim() -split '\s+'
    if ($fields.Count -ne 2 -or $fields[0] -notmatch '^[a-fA-F0-9]{64}$' -or $fields[1] -cne $name) { throw 'Invalid release checksum record' }
    $archive = Join-Path $temp $name
    Invoke-WebRequest -UseBasicParsing -Uri "$location/$name" -OutFile $archive -TimeoutSec 900
    [Console]::Error.WriteLine('Verifying archive checksum...')
    if ((Get-FileHash -Algorithm SHA256 $archive).Hash -ine $fields[0]) { throw 'Checksum mismatch; installation unchanged' }
    $package = Join-Path $temp 'package'
    [Console]::Error.WriteLine('Extracting bootstrap package...')
    Expand-Archive -LiteralPath $archive -DestinationPath $package
    [Console]::Error.WriteLine('Creating managed installation...')
    & (Join-Path $package 'fungi.exe') install --archive $archive --sha256 $fields[0] --dir $Directory
    if ($LASTEXITCODE -ne 0) { throw "Fungi installation failed with exit code $LASTEXITCODE" }
    $bin = Join-Path ([IO.Path]::GetFullPath($Directory)) 'bin'
    $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
    if (-not $NoPath -and ($userPath -split ';') -notcontains $bin) {
        [Environment]::SetEnvironmentVariable('Path', (($userPath, $bin | Where-Object { $_ }) -join ';'), 'User')
    }
    if (-not $NoPath -and ($env:Path -split ';') -notcontains $bin) { $env:Path += ";$bin" }
    if ($NoPath) { Write-Host "Fungi is installed at $bin. PATH was left unchanged." }
    else { Write-Host 'Fungi is installed. New terminals inherit its PATH entry. Update with: fungi update' }
} finally {
    Remove-Item -LiteralPath $temp -Recurse -Force
}
