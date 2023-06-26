# Change the location where Packer installs plugins on Windows
Please refer to the [docs](https://developer.hashicorp.com/packer/docs/plugins/install-plugins)
Once you now which variables to change, use `setx` or `Set-Item` to change them.
Make sure you have admin privileges and reload VSCode afterwards.

# SSH into your EC2 instance
1. Install [OpenSSH](https://learn.microsoft.com/en-us/windows-server/administration/openssh/openssh_install_firstuse?tabs=gui) client
2. `ssh -i <path_to_your_private_key> <AWS_EC2_INSTANCE_USERNAME>@<your_ec2_instance_public_ip>`