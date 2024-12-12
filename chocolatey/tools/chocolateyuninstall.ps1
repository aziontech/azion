$ErrorActionPreference = 'Stop'; # stop on all errors

# Define the paths
$binDir = Join-Path $env:ChocolateyInstall 'bin'
$azionExecutable = Join-Path $binDir 'azion.exe'

# Uninstallation process
Write-Host "Attempting to remove Azion executable from $binDir..."

if (Test-Path $azionExecutable) {
    try {
        Remove-Item -Path $azionExecutable -Force
        Write-Host "Successfully removed Azion executable from $binDir."
    } catch {
        Write-Warning "Failed to remove Azion executable from $binDir. Error: $_"
    }
} else {
    Write-Warning "Azion executable not found in $binDir. Nothing to remove."
}

# Uninstall logic for registry-based uninstallation (optional, if needed)
$packageArgs = @{
    packageName   = $env:ChocolateyPackageName
    softwareName  = 'azion*'
    fileType      = 'exe'
    silentArgs    = "/qn /norestart"
    validExitCodes= @(0, 3010, 1605, 1614, 1641)
}

[array]$key = Get-UninstallRegistryKey -SoftwareName $packageArgs['softwareName']

if ($key.Count -eq 1) {
    $key | % {
        $packageArgs['file'] = "$($_.UninstallString)"

        if ($packageArgs['fileType'] -eq 'MSI') {
            $packageArgs['silentArgs'] = "$($_.PSChildName) $($packageArgs['silentArgs'])"
            $packageArgs['file'] = ''
        }

        Uninstall-ChocolateyPackage @packageArgs
    }
} elseif ($key.Count -eq 0) {
    Write-Warning "$($packageArgs['packageName']) has already been uninstalled by other means."
} elseif ($key.Count -gt 1) {
    Write-Warning "$($key.Count) matches found!"
    Write-Warning "To prevent accidental data loss, no programs will be uninstalled."
    Write-Warning "Please alert the package maintainer the following keys were matched:"
    $key | % { Write-Warning "- $($_.DisplayName)" }
}