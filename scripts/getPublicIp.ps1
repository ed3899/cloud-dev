# Stop in case of any errors down the line
$ErrorActionPreference = "Stop"

# Retrieve public IP
$publicIP = [System.DBNull]::Value
try {
    $publicIP = (Invoke-RestMethod -Uri "https://ipinfo.io/ip").Trim()
    Write-Host "Your public IP is: $publicIP"
} catch {
    Write-Host "Error occurred while retrieving the public IP: $_"
}

# Create public IP file
# Set run location  
$initialLocation = Get-Location
$targetDirectory = "pulumi"
Set-Location -Path $targetDirectory

$jsonObject = [PsCustomObject]@{
    publicIP = $publicIP
}
$filePath = "publicIP.json"
try {
    $jsonObject | ConvertTo-Json | Out-File -FilePath $filePath -ErrorAction Stop
    Write-Host "Public IP copied to: $($targetDirectory + '/' + $filePath)"
} catch {
    Write-Host "Error creating the file: $_"
}
# Set location back
Set-Location -Path $initialLocation   

return "$targetDirectory/$filePath"