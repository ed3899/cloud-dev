# TODO improve for readability

# Tags
Add tags to your environment file located at the root folder. Each tag represents a tool.

For example:

```
# .env
ANSIBLE_TAGS = ["aws","node_js","docker"]
```

Tools listed in the section below.

Some tools may have specific requirements (i.e cpus, RAM, disk space). Please consult their respective documentation before building an AMI.
## Cloud providers
### AWS
Add `aws` to tags.

Ensure the following are present at the environment file:
```
AWS_ACCESS_KEY = "<FILL_WITH_CUSTOM_VALUE>"
AWS_SECRET_KEY = "<FILL_WITH_CUSTOM_VALUE>"
AWS_REGION = "<FILL_WITH_CUSTOM_VALUE>"
AWS_EC2_INSTANCE_USERNAME = "<FILL_WITH_CUSTOM_VALUE>"
AWS_EC2_INSTANCE_USERNAME_HOME = "<FILL_WITH_CUSTOM_VALUE>"
```
## Containerization
### Docker
Add `docker` to tags.
## IaC
### Pulumi
Add `pulumi` to tags.
## Orchestration
### Helm
Add `helm` to tags.
### Kubectl
Add `kubectl` to tags.
### Minikube
Add `docker` and `minikube` to tags.
The order matters as `docker` is a dependency.
## Programming languages
### Node.js
Add `node_js` to tags.
It uses [nvm](https://github.com/nvm-sh/nvm) to manage multiple node versions.

# Recommended base images
Although you could technically pick any image available so far on AWS, some of them may need aditional
considerations.

If you've already tried out a new image with success, please feel free open up an issue or initiate a pull request.

For now, here's the list of recommended images:

**Ubuntu Jammy 22.04 AMD64**
```
AWS_REGION = "us-west-2"
AWS_EC2_AMI_NAME_FILTER = "ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-20230516"
AWS_EC2_AMI_ROOT_DEVICE_TYPE = "ebs"
AWS_EC2_AMI_VIRTUALIZATION_TYPE = "hvm"
AWS_EC2_AMI_OWNERS = ["099720109477"]
AWS_EC2_SSH_USERNAME = "ubuntu"
```

# Change the location where Packer installs plugins on Windows
Please refer to the [docs](https://developer.hashicorp.com/packer/docs/plugins/install-plugins)
Once you now which variables to change, use `setx` or `Set-Item` to change them.
Make sure you have admin privileges and reload VSCode afterwards.

# SSH into your EC2 instance
1. Install [OpenSSH](https://learn.microsoft.com/en-us/windows-server/administration/openssh/openssh_install_firstuse?tabs=gui) client
2. `ssh -i <path_to_your_private_key> <AWS_EC2_INSTANCE_USERNAME>@<your_ec2_instance_public_ip>`
3. Type 'yes' when prompted to add your instance url to the list of known hosts

Please refer to the aws docs if you don't know where to get your ec2 instance public ip.

# Security for AWS Credentials
The approach that we took to secure them is based on architectural decisions rather than encryption and 
decryption (both at rest or in transit).

We limit the access to your instance based on you local ip, vpc and security groups at runtime.

So, we suggest you give only the least amount of permissions to the account that holds the credentials.

This facilitates the development process as you would be dealing with plain values. It also makes things easier for us as we don't have to depend on 3rd party vaults.

In case an attacker gets ahold of your aws keys (by somehow passing all the walls we've set), you can limit the impact by only giving those keys the needed permissions for a given task. Feel free to add alerts or anything in that regard. Needless to say, don't use aws root credentials.

This only accounts for placing the values inside a cloud instance using Ansible. It doesn't account for how *Packer* uses them to initially communicate with *AWS*.

## Why not vaults?
We could have opted for options such as *Ansible Vault* or *Hashicorp Vault*.

We didn't go for the primer because it would require that tool to be locally installed, and cross compatibility(*Windows* and *MacOs*) increases complexity.

The latter is not free, and the open source version adds again complexity that would only be necessary for a small portion of the utility of this tool.

## Why not secrets?
Again, we would like you to have a simple setup process. If you needed secrets setup on a 3rd party app from the get go, that would hinder the development experience. It's already enough to have to set up an account with your cloud provider.

Limit the permissions of the credentials from your cloud provider. The responsibility of ensuring they hold the right permissions lies on yourself.