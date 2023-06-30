# Requirements
- [Packer](https://developer.hashicorp.com/packer/downloads)
- Windows 10 or greater
- Powershell 7.2.11 or greater

# How-to
## Windows
Add an `.env` file at the root of the project, include your ansible [tags](#tags):
```
AWS_ACCESS_KEY = "<CUSTOM_VALUE>"
AWS_SECRET_KEY = "<CUSTOM_VALUE>"
AWS_IAM_PROFILE = "<CUSTOM_VALUE>"
AWS_USER_IDS = ["<CUSTOM_VALUE>"]
AWS_AMI_NAME = "cloud_dev"
AWS_INSTANCE_TYPE = "t2.micro"
AWS_REGION = "us-west-2"
AWS_EC2_AMI_NAME_FILTER = "ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-20230516"
AWS_EC2_AMI_ROOT_DEVICE_TYPE = "ebs"
AWS_EC2_AMI_VIRTUALIZATION_TYPE = "hvm"
AWS_EC2_AMI_OWNERS = ["099720109477"]
AWS_EC2_SSH_USERNAME = "ubuntu"
AWS_EC2_INSTANCE_USERNAME = "dev"
AWS_EC2_INSTANCE_USERNAME_HOME = "home"
AWS_EC2_INSTANCE_USERNAME_PASSWORD = "test12345"
AWS_EC2_INSTANCE_SSH_KEY_NAME = "ssh_key"
GIT_USERNAME = "<CUSTOM_VALUE>"
GIT_EMAIL = "<CUSTOM_VALUE>
ANSIBLE_TAGS = ["<CUSTOM_VALUE>"]
```
Some values were prefilled to give you an example. Feel free to replace those as well.

In case any required value is missing you will be prompted to fill it. The command will also provide a brief description of why this value is needed.

Then, at the root of the project run:
```
./scripts/deploy.ps1
```
Can't run the script? Go to [Q&A](#why-is-powershell-not-allowing-me-to-run-scripts)

Once complete, you should:
- Limit your ssh key permissions. Go to [Q&A](#how-do-i-fix-the-broad-permissions-error-when-trying-to-ssh-to-my-instance-from-powershell)
- Go to your *AWS EC2* management console and launch an instance from the recently created *AMI*, no need to add ssh keys.
- SSH into your *EC2* instance, with the key that packer downloaded to the root of the project, to verify connectivity. Go to [Q&A](#how-to-ssh-into-an-ec2-instance)

If you want to remove your *AMI*, do so from the *AWS EC2* management console.

# Tags
Add tags to your environment file located at the 
root folder. Each tag represents a tool.

For example:
```
# .env
ANSIBLE_TAGS = ["aws","node_js","docker"]
```

Tools listed in the section below.

Some tools may have specific requirements (i.e cpus, RAM, disk space). Please consult their respective documentations before building an AMI.

As always, if you think you need a specific level of custommization, all playbooks are located at `packer\ansible\playbooks`.

We welcome any [contributions](#contributions) if you think others can benefit from your changes!
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
Add `docker` and `minikube` to tags. The order matters as `docker` is a dependency.
## Programming languages
### Dotnet
Add `dotnet` to tags.

It install `dotnet-sdk-7.0` under the hood. If you need to change this, go to `packer\ansible\playbooks\programming_languages\dotnet.yml` and change it on: 
```
  tasks:
    - name: Install dotnet sdk
      ansible.builtin.apt:
        name: <CHANGE_ME_HERE_TO_THE_RIGHT_VERSION>
        update_cache: yes
```
Dotnet follows a different approach for version management. Please read [this](https://learn.microsoft.com/en-us/dotnet/core/versions/selection).
### Go
Add `go` to tags.

To manage multiple Go versions, please refer to this [thread](https://go.dev/doc/manage-install#installing-multiple) and this [one](https://stackoverflow.com/a/68087898/11941146). The latter shows how to alias multiple Go versions.
### Node.js
Add `node_js` to tags.

It uses [nvm](https://github.com/nvm-sh/nvm) to manage multiple node versions.
### Python
Add `python` to tags.

It installs [miniconda](https://docs.conda.io/en/latest/miniconda.html#linux-installers) under the hood which is a lightweight version of [anaconda](https://www.anaconda.com/). You can manage multiple python versions, virtual environment and packages with it.

If want a different *miniconda* version. Please refer to this [download](https://docs.conda.io/en/latest/miniconda.html#linux-installers) page and change the following values accordingly at `packer\ansible\playbooks\programming_languages\python.yml`
```
  vars:
    x86_64:
      anaconda_url: <ADEQUATE_URL>
      anaconda_sha256: <ADEQUATE_SHA256>
    aarch64:
      anaconda_url: <ADEQUATE_URL>
      anaconda_sha256: <ADEQUATE_SHA256>
```
As well as:
```
- name: Ensure needed python version for Anaconda3 is present
  ansible.builtin.apt:
    name: <ADEQUATE_PYTHON_VERSION>
    update_cache: yes
```
The reason we picked *miniconda* is because we make the assumption you would want to avoid paying for extra storage from the very beginning. Installing anaconda requires at least *3gb* of storage which would require us to deploy an *AMI* with a bigger volume. The decision then relies on you if want additional storage.
## Terminal
### Starship
Add `starship` to your tags.

For more information consult its [docs](https://starship.rs/).
## Version control
### GitHub
Add `github` to your tags. Automate the process of adding ssh keys.

Add the following to your environment file.
```
#.env
GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC = "<CUSTOM_VALUE>"
```
Make sure you give the right permissions to the token (i.e "repo", "admin:public_key" are the minimum). For more granular control consult the [docs](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token).

You can test your connection with:
```
ssh -T git@github.com
```
Type `yes`, click enter and you should get a message like the following:
```
Hi <YOUR_GITHUB_USER>! You've successfully authenticated, but GitHub does not provide shell access.
```

# Q&A
## Why is powershell not allowing me to run scripts?
The most common cause is that you are running the script from a non-admin shell or that your remote policy is not allowing you. In any case, open up a Powershell shell as an admin and run the following command.
```
Set-ExecutionPolicy RemoteSigned
```
## What are the recommended Ubuntu images?
Although you could technically pick any image available so far on *AWS*, some of them may need aditional
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
## How to change the location where Packer installs plugins on Windows?
Please refer to the [docs](https://developer.hashicorp.com/packer/docs/plugins/install-plugins)
Once you now which variables to change, use `setx` or `Set-Item` to change them.

Make sure you have admin privileges and reload VSCode afterwards.

## How to SSH into an *EC2* instance?
1. Install [OpenSSH](https://learn.microsoft.com/en-us/windows-server/administration/openssh/openssh_install_firstuse?tabs=gui) client
2. `ssh -i <path_to_your_private_key> <AWS_EC2_INSTANCE_USERNAME>@<your_ec2_instance_public_ip>`
3. Type 'yes' when prompted to add your instance url to the list of known hosts

Please refer to the aws docs if you don't know where to get your ec2 instance public ip.
## How do I fix the broad permissions error when trying to ssh to my instance from Powershell?
Run the following command from a powershell admin shell. You only need to do this once unless to remove the ssh key file.
```
icacls <AWS_EC2_INSTANCE_SSH_KEY_NAME> /inheritance:r /grant:r "$($env:USERNAME):(R,W)"
```
## How secure are my *AWS* credentials?
The approach that we took to secure them is based on architectural decisions rather than encryption and 
decryption (both at rest or in transit).

We limit the access to your instance based on you local ip, vpc and security groups at runtime.

So, we suggest you give only the least amount of permissions to the account that holds the credentials.

This facilitates the development process as you would be dealing with plain values. It also makes things easier for us as we don't have to depend on 3rd party vaults.

In case an attacker gets ahold of your aws keys (by somehow passing all the walls we've set), you can limit the impact by only giving those keys the needed permissions for a given task. Feel free to add alerts or anything in that regard. Needless to say, don't use aws root credentials.

This only accounts for placing the values inside a cloud instance using *Ansible*. It doesn't account for how *Packer* uses them to initially communicate with *AWS*.
### Why not vaults?
We could have opted for options such as *Ansible Vault* or *Hashicorp Vault*.

We didn't go for the primer because it would require that tool to be locally installed, and cross compatibility(*Windows* and *MacOs*) increases complexity.

The latter is not free, and the open source version adds again complexity that would only be necessary for a small portion of the utility of this tool.
### Why not secrets?
Again, we would like you to have a simple setup process. If you needed secrets setup on a 3rd party app from the get go, that would hinder the development experience. It's already enough to have to set up an account with your cloud provider.

Limit the permissions of the credentials from your cloud provider. The responsibility of ensuring they hold the right permissions lies on yourself.

# Contributions
Please, help us make this project more **usable** and **cross-compatible**. Feel free to open up either an **issue** or a **pull request via forking**.
## Design philosophy
We strive towards:
- Setting up an entire cloud development environment with as few commands as possible (assuming you've already setup your cloud credentials)
- Being cross Ubuntu compatible (any other distros are welcomed but for now having at least 8 working *Ubuntu* images is our primary goal)
- Being *Windows* and *MacOs* compatible
- Being *AWS*, *Google* and *Azure* compatible
- Single user experience
- Multiple version handling regarding specific tools
- Security and shared responsibility regarding cloud keys
- Prefer tools that are core essentials to a productive development environment across the board and that go beyond being installed with 1 single command or that can be easily accesible via *Docker* or related technologies.


# TODO improve for readability