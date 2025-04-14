$ErrorActionPreference = 'Stop' # Stop on all errors

# Define paths
$toolsDir   = "$(Split-Path -Parent $MyInvocation.MyCommand.Definition)"
$outputFile = Join-Path $toolsDir 'azion.exe'

# Define package details
$url        = 'http://downloads.azion.com/windows/x86_64/azion'
$checksum   = '{{CHECKSUM}}' # Replaced with actual checksum during deploy

# Download the CLI binary using Chocolatey helper
Get-ChocolateyWebFile -PackageName 'azion' -FileFullPath $outputFile -Url $url -Checksum $checksum -ChecksumType 'sha256'

# No need to manually copy or shim - Chocolatey will shim the .exe in the tools directory automatically
Write-Host "Installation complete. The Azion CLI is now available globally via the command line."
