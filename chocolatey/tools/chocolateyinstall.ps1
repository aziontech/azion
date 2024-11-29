$ErrorActionPreference = 'Stop'; # Stop on all errors

# Define paths
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$binDir     = Join-Path $env:ChocolateyInstall 'bin'
$outputFile = Join-Path $toolsDir 'azion.exe'

# Define package details
$url        = 'http://downloads.azion.com/windows/x86_64/azion'
$checksum = '{{CHECKSUM}}'
$silentArgs = ''
$packageArgs = @{
    packageName   = 'azion'
    unzipLocation = $toolsDir
    fileType      = 'exe'
    url           = $url

    softwareName  = 'azion*'

    checksum      = $checksum         
    checksumType  = 'sha256'   

    silentArgs    = $silentArgs
}

# Install the package
Install-ChocolateyPackage @packageArgs

# Download the file
Write-Host "Downloading azion executable from $url to $outputFile"
Invoke-WebRequest -Uri $url -OutFile $outputFile

# Ensure Chocolatey's bin directory exists
if (-Not (Test-Path $binDir)) {
    New-Item -ItemType Directory -Path $binDir | Out-Null
}

# Copy the executable to Chocolatey's bin directory
Write-Host "Copying $outputFile to $binDir"
Copy-Item -Path $outputFile -Destination $binDir -Force

# Ensure the executable is available globally
Write-Host "Installation complete. Azion executable is now available globally."