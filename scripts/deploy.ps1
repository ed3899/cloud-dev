# Stop in case of any errors down the line
$ErrorActionPreference = "Stop"

# Check if packer is installed
$packerExe = "packer"
$packerInstalled = $null
$arc = $env:PROCESSOR_ARCHITECTURE
try {
    $packerInstalled = Get-Command $packerExe -ErrorAction Stop
} catch [System.Management.Automation.CommandNotFoundException] {
    Write-Host "Packer is not installed. Please go to https://developer.hashicorp.com/packer/downloads and follow the Windows $($arc) installation instructions"
} catch {
    Write-Host "An error occurred while checking Packer installation."
}

if ($packerInstalled) {
  # Extract environment variables
  $variables = Get-Content -Raw -Path ".env"
  $environmentVariables = @{}
  $variables -split "`r`n" | ForEach-Object {
    $parts = $_ -split "="
    if ($parts.Length -eq 2) {
        $key = $parts[0].Trim()
        $value = $parts[1].Trim()
        $environmentVariables[$key] = $value
    }
  }  

  # Access environment variables
  $sshKeyName = $environmentVariables["AWS_EC2_INSTANCE_SSH_KEY_NAME"]

  # Copy .env file to Packer folder
  Copy-Item .env packer/.auto.pkrvars.hcl
  # Copy .env file to Pulumi folder
  Copy-Item .env pulumi/.env

  # Set run location  
  $initialLocation = Get-Location
  $targetDirectory = "packer"
  Set-Location -Path $targetDirectory

  # Build  
  Write-Host "Building..."
  packer build .
  Write-Host "Done!"

  # Set location back
  Set-Location -Path $initialLocation
}


