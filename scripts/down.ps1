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

  Write-Host "Destroying..."
  $env:PULUMI_ACCESS_TOKEN = $environmentVariables["PULUMI_PERSONAL_ACCESS_TOKEN"]
  $envVarsAbsent = ($env:PULUMI_ACCESS_TOKEN -eq $null) -or ($environmentVariables["AWS_ACCESS_KEY"] -eq $null) -or ($environmentVariables["AWS_SECRET_KEY"] -eq $null) -or ($environmentVariables["AWS_REGION"] -eq $null)
  if ($envVarsAbsent) {
    throw "Error: environment variables PULUMI_PERSONAL_ACCESS_TOKEN, AWS_ACCESS_KEY, AWS_SECRET_KEY and AWS_REGION must be present"
  }
  
  pulumi config set --secret aws:accessKey $environmentVariables["AWS_ACCESS_KEY"]
  pulumi config set --secret aws:secretKey $environmentVariables["AWS_SECRET_KEY"]
  pulumi config set --secret aws:region $environmentVariables["AWS_REGION"]
  pulumi destroy
  if ($?) {
    Write-Host "Destroyed!"
  }

  # Set location back
  Set-Location -Path $initialLocation   
}