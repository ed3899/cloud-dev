#!/bin/bash

# Copy .env file to Pulumi folder
cp .env pulumi/.env

# Copy .env file to Packer folder
cp .env packer/.auto.pkrvars.hcl
