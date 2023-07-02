# Stop in case of any errors down the line
$ErrorActionPreference = "Stop"

$pass = ./scripts/verify_dependencies.ps1

if ($pass) {
  # Copy .env file to Pulumi folder
  Copy-Item .env pulumi/.env

  # Set run location  
  $initialLocation = Get-Location
  $targetDirectory = "pulumi"
  Set-Location -Path $targetDirectory

  # Extract environment variables
  # This step is only needed initially because
  # of AWS Credentials. Later on we can extract
  # environment variables from an env file.
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

  Write-Host "Creating stack..."
  $env:PULUMI_ACCESS_TOKEN = $environmentVariables["PULUMI_PERSONAL_ACCESS_TOKEN"]
  if ($env:PULUMI_ACCESS_TOKEN -eq $null) {
    throw "Error: PULUMI_PERSONAL_ACCESS_TOKEN environment variable must be present."
  }
  pulumi stack init --stack dev --non-interactive

  # Set location back
  Set-Location -Path $initialLocation   
}