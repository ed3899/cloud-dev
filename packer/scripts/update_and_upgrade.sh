#!/bin/bash
set -euo pipefail  # Exit immediately if a command fails, unset variables are referenced, or pipelines return errors

sudo apt-get update -qq
sudo apt-get upgrade -qq
