# Copy .env file to Pulumi folder
Copy-Item .env pulumi/.env

# Copy .env file to Packer folder
Copy-Item .env packer/.auto.pkrvars.hcl
