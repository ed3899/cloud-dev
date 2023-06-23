#!/bin/bash
# Install Ansible on Ubuntu
set -euo pipefail  # Exit immediately if a command fails, unset variables are referenced, or pipelines return errors
sudo apt-get install -y software-properties-common
sudo apt-add-repository --yes --update ppa:ansible/ansible
sudo apt-get install -y ansible 
