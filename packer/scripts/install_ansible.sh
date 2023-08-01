#!/bin/bash
set -uo pipefail  # Exit immediately if unset variables are referenced, or pipelines return errors
sudo apt-get install -y software-properties-common
sudo apt-add-repository --yes --update ppa:ansible/ansible
sudo apt-get install -y ansible 
