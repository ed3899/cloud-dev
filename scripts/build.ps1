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
  $errorOutput = (packer init -upgrade 2>&1)

  if ($? -ne 0) {
    Write-Host "There was an error while initiating packer."
    Write-Host "Command output:"
    Write-Host $errorOutput
  }

  # Copy .env file to Packer folder
  Copy-Item .env packer/.auto.pkrvars.hcl

  # Set run location  
  $initialLocation = Get-Location
  $targetDirectory = "packer"
  Set-Location -Path $targetDirectory

  # Extract environment variables
  # This step is only needed initially because
  # of AWS Credentials. Later on we can extract
  # environment variables from an env file.
  $variables = Get-Content -Raw -Path ".auto.pkrvars.hcl"
  $environmentVariables = @{}
  $variables -split "`r`n" | ForEach-Object {
    $parts = $_ -split "="
    if ($parts.Length -eq 2) {
        $key = $parts[0].Trim()
        $value = $parts[1].Trim()
        $environmentVariables[$key] = $value
    }
  }

  $sshKeyName = $environmentVariables["AWS_EC2_INSTANCE_SSH_KEY_NAME"] ?? 'AWS_EC2_INSTANCE_SSH_KEY_NAME'
  $sshUser = $environmentVariables["AWS_EC2_INSTANCE_USERNAME"] ?? 'AWS_EC2_INSTANCE_USERNAME'

  # Build  
  Write-Host "Building..."
  packer build . 
  if ($?) {
    Write-Host "Built!"
    Write-Host "Once launched, ssh into your instance with the following command:"
    Write-Host "ssh -i $($sshKeyName.Trim('"')) $($sshUser.Trim('"'))@YOUR_EC2_PUBLIC_IP"
  }

  # Set location back
  Set-Location -Path $initialLocation
}


