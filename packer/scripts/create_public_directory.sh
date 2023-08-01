#!/bin/bash
set -e # halt on any error

# Verify that variable is set and contains a valid path.
if [ -z "$AWS_EC2_PUBLIC_DIRECTORY_INTERNAL" ] ; then
    echo "Error: empty directory path."
    exit 1
fi

sudo mkdir "$AWS_EC2_PUBLIC_DIRECTORY_INTERNAL" && \
sudo chown "$AWS_EC2_SSH_USERNAME":"$AWS_EC2_SSH_USERNAME" "$AWS_EC2_PUBLIC_DIRECTORY_INTERNAL" && \
sudo chmod -R 744 "$AWS_EC2_PUBLIC_DIRECTORY_INTERNAL" && \

# Verify that directory was created successfully.
if  [ ! -d "$AWS_EC2_PUBLIC_DIRECTORY_INTERNAL" ]; then
    echo "Error: failed to create directory."
    exit 1
fi

echo "Directory $AWS_EC2_PUBLIC_DIRECTORY_INTERNAL has been created."