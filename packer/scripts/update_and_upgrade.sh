#!/bin/bash
set -uo pipefail  # Exit immediately if unset variables are referenced, or pipelines return errors
sudo apt-get update -qq
sudo apt-get upgrade -qq
