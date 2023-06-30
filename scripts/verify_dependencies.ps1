# Stop in case of any errors down the line
$ErrorActionPreference = "Stop"

# Check if pulumi is installed
$pulumiExe = "pulumi"
$pulumiInstalled = $null
$arc = $env:PROCESSOR_ARCHITECTURE
try {
    $pulumiInstalled = Get-Command $pulumiExe -ErrorAction Stop
} catch [System.Management.Automation.CommandNotFoundException] {
    Write-Host "Pulumi is not installed. Please go to https://www.pulumi.com/docs/install/ and follow the Windows $($arc) installation instructions"
} catch {
    Write-Host "An error occurred while checking Pulumi installation."
}

# Check if node is installed
$nodeExe = "node"
$nodeInstalled = $null
$systemInfo = Get-WmiObject -Class Win32_OperatingSystem
$architecture = $systemInfo.OSArchitecture
try {
    $nodeInstalled = Get-Command $nodeExe -ErrorAction Stop
} catch [System.Management.Automation.CommandNotFoundException] {
    Write-Host "Node is not installed. Please go to https://nodejs.org/en/download and follow the Windows $($architecture) installation instructions"
} catch {
    Write-Host "An error occurred while checking Node installation."
}

return $true