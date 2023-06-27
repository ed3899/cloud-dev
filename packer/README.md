# Change the location where Packer installs plugins on Windows
Please refer to the [docs](https://developer.hashicorp.com/packer/docs/plugins/install-plugins)
Once you now which variables to change, use `setx` or `Set-Item` to change them.
Make sure you have admin privileges and reload VSCode afterwards.

# SSH into your EC2 instance
1. Install [OpenSSH](https://learn.microsoft.com/en-us/windows-server/administration/openssh/openssh_install_firstuse?tabs=gui) client
2. `ssh -i <path_to_your_private_key> <AWS_EC2_INSTANCE_USERNAME>@<your_ec2_instance_public_ip>`

Please refer to the aws docs if you don't know where to get your ec2 instance public ip.

# Security for AWS Credentials
The approach that we took to secure them is based on architectural decisions rather than encryption and 
decryption (both at rest or in transit).

We limit the access to your instance based on you local ip, vpc and security groups at runtime.

So, we suggest you give only the least amount of permissions to the account that holds the credentials.

This facilitates the development process as you would be dealing with plain values. It also makes things easier for us as we don't have to depend on 3rd party vaults.

In case an attacker get ahold of your aws keys (by somehow passing all the walls we've set), you can limit the impact by only giving those keys the needed permissions for a given task. Feel free to add alerts or anything in that regard. Needless to say, don't use aws root credentials.

This only accounts for placing the values inside a cloud instance using Ansible. It doesn't account for how *Packer* uses them to initially communicate with *AWS*.

## Why not vaults?
We could have opted for options such as *Ansible Vault* or *Hashicorp Vault*.

We didn't go for the primer because it would require that tool to be locally installed, and cross compatibility(*Windows* and *MacOs*) increases complexity.

The latter is not free, and the open source version adds again complexity that would only be necessary for a small portion of the utility of this tool.

## Why not secrets?
Again, we would like you to have a simple setup process. If you needed secrets setup on a 3rd party app from the get go, that would hinder the development experience. It's already enough to have to set up an account with your cloud provider.

Limit the permissions of the credentials from the cloud provider. The responsibility of ensuring they hold the right permissions lies on yourself.