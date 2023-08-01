#!/bin/bash
set -e # halt on any error

if [ -z "$AWS_EC2_PUBLIC_DIRECTORY_INTERNAL" ]; then
  echo "Error: AWS_EC2_PUBLIC_DIRECTORY_INTERNAL variable is not set."
  exit 1
fi

sudo rm -rf "$AWS_EC2_PUBLIC_DIRECTORY_INTERNAL"

if [ $? -eq 0 ]; then
  echo "Directory $AWS_EC2_PUBLIC_DIRECTORY_INTERNAL has been deleted."
else
  echo "Error: Failed to delete directory $AWS_EC2_PUBLIC_DIRECTORY_INTERNAL."
  exit 1
fi
