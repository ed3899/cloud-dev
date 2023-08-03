# Overview

A cloud development environment you can customize with a wide range of tools.

## Table of contents

- [Requirements](#requirements)
  - [Optionals](#optionals)
- [How-to](#how-to)
  - [Windows](#windows)
- [Tags](#tags)
  - [Cloud providers](#cloud-providers)
    - [aws](#aws)
  - [Containerization](#containerization)
    - [docker](#docker)
  - [IaC](#iac)
    - [pulumi](#pulumi)
  - [Orchestration](#orchestration)
    - [helm](#helm)
    - [kubectl](#kubectl)
    - [minikube](#minikube)
  - [Programming languages](#programming-languages)
    - [dotnet](#dotnet)
    - [go](#nodejs)
    - [node.js](#nodejs)
    - [python](#python)
    - [ruby](#ruby)
    - [rust](#rust)
  - [Terminal](#terminal)
    - [starship](#starship)
  - [Version control](#version-control)
    - [github](#github)
- [Extras](#extras)
  - [Unix aliases](#unix-aliases)
  - [Git aliases](#git-aliases)
- [Q&A](#qa)
- [Contributions](#contributions)

## Requirements

- Windows 10 or a later version is required.
- [Packer](https://developer.hashicorp.com/packer/downloads) version 1.2.0 to less than 2.0.0 is needed for building AMIs. Their management is done via the AWS console.
- [Powershell](https://learn.microsoft.com/en-us/powershell/scripting/install/installing-powershell-on-windows?view=powershell-7.3) version 7.2.11 or higher is required as it's used for writing script utilities.
- [OpenSSH](https://learn.microsoft.com/en-us/windows-server/administration/openssh/openssh_install_firstuse?tabs=powershell) is necessary for SSH access into the instance, specifically the client part.

## Optionals

- [Pulumi](https://www.pulumi.com/docs/install/): version 3.73.0 or greater (deploy and manage your instance from your CLI).
- [Node](https://nodejs.org/en/download): version 18.16.1 or greater (required for the pulumi backend).

## How-to

Follow these steps to run.

### Windows

1. Create an `.env` file at the root of the project
2. Make sure to fill in with the required values.
3. Add your ansible [tags](#tags):

    ```.env
      AWS_ACCESS_KEY = "<CUSTOM_VALUE>"
      AWS_SECRET_KEY = "<CUSTOM_VALUE>"
      AWS_IAM_PROFILE = "<CUSTOM_VALUE>"
      AWS_USER_IDS = ["<CUSTOM_VALUE>"]
      AWS_AMI_NAME = "cloud_dev"
      AWS_INSTANCE_TYPE = "t2.micro" # Must be compatible with your AMI
      AWS_REGION = "us-west-2" # AMI must be available on this region
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
      # Optionals
      AWS_EC2_INSTANCE_VOLUME_TYPE = "gp2" # default
      AWS_EC2_INSTANCE_VOLUME_SIZE = "8" # default
      PULUMI_PERSONAL_ACCESS_TOKEN = "<CUSTOM_VALUE>"
    ```

4. At the root of your project, run the following command:

    ```powershell
    ./scripts/build.ps1
    ```

*If you are unable to run the script, please refer to the [Q&A](#why-is-powershell-not-allowing-me-to-run-scripts) section for guidance on why PowerShell may not be allowing you to run scripts.*

*If you want more information on how to pick the right EC2 instance, please go to [Q&A](#how-to-pick-the-right-ec2-instance).*

#### Post-Build Steps

Once the script completes, follow these steps:

  1. Limit the permissions on your *SSH key*. Please refer to the [Q&A](#how-do-i-fix-the-broad-permissions-error-when-trying-to-ssh-to-my-instance-from-powershell) section for guidance on how to fix the broad permissions error when trying to *SSH* to your instance from *PowerShell*.
  2. In your *AWS EC2 management console*, launch an instance using the recently created *AMI*. You do not need to add *SSH keys* during the launch process.
  3. SSH into your *EC2* instance using the *SSH key* that *Packer* downloaded to the root of your project. For instructions on how to SSH into an *EC2* instance, please refer to the [Q&A](#how-to-ssh-into-an-ec2-instance) section.

If you want to remove your *AMI*, you can do so from the *AWS EC2 management console*.

#### Manage resources from cli

Make sure to add the following to your `.env` file:

```.env
PULUMI_PERSONAL_ACCESS_TOKEN = "CUSTOM_VALUE"
```

To get a Pulumi personal access token, please refer to the [Pulumi documentation](https://www.pulumi.com/docs/pulumi-cloud/access-management/access-tokens/#personal-access-tokens).

Follow the steps below in the given order:

1. Create stack (this will install npm modules if not present):

    ```powershell
    ./scripts/init_stack.ps1
    ```

2. Deploy

    ```powershell
    ./scripts/up.ps1
    ```

    - Once deployed, an SSH config file will be generated at the root of the project. You can use it to [remote SSH from Visual Studio Code](#how-to-remote-ssh-from-vs-code).

3. Destroy

    ```powershell
    ./scripts/down.ps1
    ```

4. Delete stack

    ```powershell
    ./scripts/remove_stack.ps1
    ```

Before removing a stack, always make sure that you have already destroyed the resources deployed.

Make sure to follow these steps carefully to manage your resources effectively.

## Tags

In your root folder, you need to add tags to your environment file. Each tag represents a tool.

For example, modify your `.env` file as follows:

```.env
ANSIBLE_TAGS = ["aws","node_js","docker"]
```

Please note that there are several tools listed below. Some of these tools may have specific requirements such as CPU, RAM, or disk space. Before building an AMI, it is recommended to consult their respective documentations for any specific requirements.

If you require a specific level of customization, you can find all the playbooks located at `packer\ansible\playbooks.`

Make sure to add the appropriate tags to your environment file, and if needed, refer to the playbooks to customize your setup accordingly.

### Cloud providers

#### AWS

Add `aws` to tags.

Ensure the following are present at the environment file:

```.env
AWS_ACCESS_KEY = "<FILL_WITH_CUSTOM_VALUE>"
AWS_SECRET_KEY = "<FILL_WITH_CUSTOM_VALUE>"
AWS_REGION = "<FILL_WITH_CUSTOM_VALUE>"
AWS_EC2_INSTANCE_USERNAME = "<FILL_WITH_CUSTOM_VALUE>"
AWS_EC2_INSTANCE_USERNAME_HOME = "<FILL_WITH_CUSTOM_VALUE>"
```

### Containerization

#### Docker

Add `docker` to tags.

### IaC

#### Pulumi

Add `pulumi` to tags.

### Orchestration

#### Helm

Add `helm` to tags.

#### Kubectl

Add `kubectl` to tags.

#### Minikube

Add `docker` and `minikube` to tags. The order matters as `docker` is a dependency.

### Programming languages

#### Dotnet

Add `dotnet` to tags.

It install `dotnet-sdk-7.0` under the hood. If you need to change this, go to `packer\ansible\playbooks\programming_languages\dotnet.yml` and change it on:

```yaml
  tasks:
    - name: Install dotnet sdk
      ansible.builtin.apt:
        name: <CHANGE_ME_HERE_TO_THE_RIGHT_VERSION>
        update_cache: yes
```

Please note that dotnet follows a different approach for version management. For more information on how to select the correct version, please read the [official documentation](https://learn.microsoft.com/en-us/dotnet/core/versions/selection).

If you have made the necessary changes to the `dotnet.yml` file and specified the correct version, save the file. The `dotnet-sdk-7.0` will be installed according to the updated configuration.

#### Go

Add `go` to tags.

Refer to the official [Go documentation](https://go.dev/doc/manage-install#installing-multiple) to learn how to manage multiple *Go* versions effectively. This guide provides detailed instructions on installing and managing multiple *Go* versions on your system.

Additionally, you can also check out this [Stack Overflow thread](https://stackoverflow.com/a/68087898/11941146), which explains how to alias multiple *Go* versions. This technique allows you to switch between different *Go* versions easily.

#### Node.js

Add `node_js` to tags.

It installs [nvm](https://github.com/nvm-sh/nvm), which is a popular tool for managing multiple versions of *Node.js*.

Once *nvm* is installed, you can use it to install and manage multiple *Node.js* versions on your system. By using *nvm* commands, you can easily switch between different node versions based on your requirements.

#### Python

Add `python` to tags.

It installs [miniconda](https://docs.conda.io/en/latest/miniconda.html#linux-installers) which is a lightweight version of [anaconda](https://www.anaconda.com/). *Miniconda* allows you to manage multiple *Python* versions, virtual environments, and packages efficiently.

If you prefer a different *miniconda* version, you can refer to the download page for alternative options. Once you have chosen the appropriate version, make the following changes in the `packer\ansible\playbooks\programming_languages\python.yml` file:

```yaml
  vars:
    x86_64:
      anaconda_url: <ADEQUATE_URL>
      anaconda_sha256: <ADEQUATE_SHA256>
    aarch64:
      anaconda_url: <ADEQUATE_URL>
      anaconda_sha256: <ADEQUATE_SHA256>
```

Additionally, modify the following section:

```yaml
- name: Ensure needed python version for Anaconda3 is present
  ansible.builtin.apt:
    name: <ADEQUATE_PYTHON_VERSION>
    update_cache: yes
```

Replace `<ADEQUATE_URL>`, `<ADEQUATE_SHA256>`, and `<ADEQUATE_PYTHON_VERSION>` with the appropriate values based on the *miniconda* version you choose.

The reason we suggest using *miniconda* is to save storage space from the start. Installing *anaconda* requires at least 3GB of storage, which would require deploying an *AMI* with a larger volume. The decision to use additional storage is up to you.

### Ruby

Add `ruby` to your tags.

Manage versions with [rbenv](https://github.com/rbenv/rbenv).

#### Rust

Add `rust` to your tags.

Manage versions with [rustup-init](https://forge.rust-lang.org/infra/other-installation-methods.html).

### Terminal

#### Starship

Add `starship` to your tags.

For more information, please consult its [official documentation](https://starship.rs/).

### Version control

#### GitHub

Add `github` to your tags and automate the process of adding *SSH* keys.

1. Update your environment file `.env` by adding the following line:

    ```.env
    GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC = "<CUSTOM_VALUE>"
    ```

2. Make sure to replace `<CUSTOM_VALUE>` with the appropriate value for your [personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token).
3. Grant the necessary permissions to the token. At a minimum, the following permissions are required: "repo" and "admin:public_key".
4. Test your connection using the following command:

  ```bash
  ssh -T git@github.com
  ```

  When prompted, type `yes` and press enter. You should receive a message similar to the following:

  ```bash
  Hi <YOUR_GITHUB_USER>! You've successfully authenticated, but GitHub does not provide shell access.
  ```

## Extras

### Unix aliases

If you're looking for some handy *Unix* aliases, you're in luck! We have included a collection of useful aliases that you can use. To costumize them, navigate to the `packer/ansible/playbooks/base/linux_aliases.yml file`

### Git aliases

We haven't forgotten about *Git*! In fact, we have also included some helpful *Git* aliases for you to use. To costumize them, go to the `packer/ansible/playbooks/base/git.yml` file. These aliases will make your *Git* commands more efficient and easier to remember.

## Q&A

### Why is powershell not allowing me to run scripts?

If you encounter issues running a script due to a non-admin shell or remote policy restrictions, follow these steps to resolve the problem:

1. Open PowerShell with administrative privileges by doing the following:
    - Press Windows key to open the Start menu.
    - Type "PowerShell".
    - Right-click on "Windows PowerShell" (or "PowerShell") and select "Run as administrator".
2. Run the following command in the elevated PowerShell window:

  ```powershell
  Set-ExecutionPolicy RemoteSigned
  ```

## What are the recommended Ubuntu images?

When choosing an image on AWS (Amazon Web Services), it is important to consider various factors. While technically any available image can be picked, some may require additional considerations. If you have successfully tried a new image, we encourage you to open an issue or initiate a pull request.

For now, here is a list of recommended images:

### Ubuntu Jammy 22.04 AMD64

```.env
AWS_REGION = "us-west-2"
AWS_EC2_AMI_NAME_FILTER = "ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-20230516"
AWS_EC2_AMI_ROOT_DEVICE_TYPE = "ebs"
AWS_EC2_AMI_VIRTUALIZATION_TYPE = "hvm"
AWS_EC2_AMI_OWNERS = ["099720109477"]
AWS_EC2_SSH_USERNAME = "ubuntu"
```

By selecting the appropriate image based on the provided details, you can ensure compatibility and meet your requirements when working on AWS.

## How to change the location where Packer installs plugins on Windows?

Please refer to the [official documentation](https://developer.hashicorp.com/packer/docs/plugins/install-plugins) for detailed instructions on installing plugins.

Once you have identified the variables that need to be changed, you can use `setx` or `Set-Item` to modify them.

**Note:** Make sure you have administrative privileges before attempting to modify system variables. After making changes, be sure to reload *VSCode* for the modifications to take effect.

## How to SSH into an *EC2* instance?

1. Install the OpenSSH client on your local machine.
2. Open the command prompt or terminal and run the following command:

    ```bash
    ssh -i <path_to_your_private_key> <AWS_EC2_INSTANCE_USERNAME>@<your_ec2_instance_public_ip>
    ```

3. When prompted to add your instance URL to the list of known hosts, type `yes` and press enter.

Please refer to the *aws docs* if you don't know where to get your *ec2* instance public ip.

## How do I fix the broad permissions error when trying to ssh to my instance from Powershell?

Run the following command from a PowerShell admin shell. You only need to do this once unless you removed the SSH key file.

Make sure the username is the same that will be sshing into the ec2 instance, otherwise change it to the appropiate one.

```powershell
icacls <AWS_EC2_INSTANCE_SSH_KEY_NAME> /inheritance:r /grant:r "$($env:USERNAME):(R,W,D)"
```

## How secure are my *AWS* credentials?

The approach we took to secure your credentials is based on architectural decisions rather than heavy reliance on encryption and decryption.

### Access Limitations at Runtime

- To limit access to your instances, we control access based on your local IP and security groups.
- By giving the account that holds the credentials the least amount of permissions, we facilitate the development process, as you work with plain values.
- This also simplifies our workflow, as we do not have to depend on third-party vaults.

### Mitigating Impact in Case of Breach

- In the unlikely event that an attacker gains access to your AWS keys, you can restrict their impact by granting only the necessary permissions for a specific task.
- We encourage you to set up alerts or other monitoring mechanisms to provide early detection of any unauthorized activity.
- Remember, it is important not to use AWS root credentials to further enhance security.

### Configuration Details

Please note that the information provided above specifically accounts for placing the values inside a cloud instance using *Ansible*. It does not cover how *Packer* utilizes these credentials for initial communication with *AWS*.

### Why not vaults?

We had the option to use tools such as *Ansible Vault* or *Hashicorp Vault* for securing credentials. However, we decided against them due to the following reasons:

- **Complexity of Local Installation**: Using *Ansible Vault* would require this tool to be locally installed, which could increase complexity, especially considering cross-compatibility issues on *Windows* and *MacOS*.
- **Cost**: *Hashicorp Vault* is not free, and even the open-source version adds complexity that would only be necessary for a small portion of the utility of this tool.

By avoiding these tools, we aim to simplify the setup and maintenance process, ensuring a smoother experience for both contributors and users.

### Why not secrets?

We prioritize a simple setup process for your convenience. Therefore, we suggest the following approach when it comes to secrets setup and permissions:

- **Minimize External Dependencies:** Having to set up secrets on a 3rd party app right from the start can hinder the development experience. Hence, we recommend avoiding this additional complexity unless absolutely necessary.
- **Focus on Cloud Provider Credentials:** By limiting the permissions of your cloud provider credentials, you remain in control of ensuring they have the correct permissions. This puts the responsibility squarely on your own shoulders.

Our aim is to strike a balance between simplicity and security, providing you with a streamlined and efficient development experience while maintaining control over the permissions granted to your cloud provider credentials.

## How to remote ssh from VS Code?

  1. Ensure that you have the [Remote - SSH extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-ssh) installed in your VS Code. You can find and install it from the Extensions view (`Ctrl+Shift+X` or View -> Extensions).
  2. Once the extension is installed, click on the Extensions view button and search for "Remote - SSH". Click on the "Install" button to install the extension.
  3. Next, open the Command Palette in VS Code by pressing `Ctrl+Shift+P` (or View -> Command Palette). Type "Remote-SSH: Open Configuration File" and select the option from the list. This will open the SSH config file in the editor.
  4. In the SSH config file, you need to add the configuration for the remote SSH connection. The format of the entry in the config file should be as follows (these values are generated at runtime):

      ```ssh_config
      Host cloud_dev
        HostName 52.10.25.215
        IdentityFile "D:\Documents\DevEnvironments\Cloud\Jammy64\ssh_key"
        User dev
        Port 22
        StrictHostKeyChecking no
        PasswordAuthentication no
        IdentitiesOnly yes
        LogLevel INFO
      ```

  5. Save the SSH config file after adding the configuration.
  6. Go back to the Command Palette (`Ctrl+Shift+P`) and type "Remote-SSH: Connect to Host". You should see the list of configured SSH hosts from the SSH config file.
  7. Select the desired host from the list. It will attempt to establish an SSH connection to the remote machine using the provided configuration.
  8. If everything is set up correctly, a new window will open in VS Code connected to the remote machine via SSH.

## How to pick the right EC2 instance

To help you make a decision based on factors such as cost and development needs, here are some tools that you can use:

- [AWS Pricing Calculator](https://calculator.aws/#/)
- [EC2 Instance Types](https://aws.amazon.com/ec2/instance-types/)
- [EC2 docs](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html)

## How to send files from my host to my instance?

```bash
scp -i /path/to/keypair.pem /path/to/local/file1 /file2 AWS_EC2_INSTANCE_USERNAME@ec2-instance-ip:/path/on/ec2
```

Works both on powershell and bash.

## Contributions

We would appreciate your assistance in making this project more **usable** and **cross-compatible**. There are two ways you can contribute:

  1. **Open an Issue**: If you come across any problems or have suggestions for improvement, please open up an issue describing your concern or idea.
  2. **Submit a Pull Request via Forking**: If you have implemented changes or enhancements and want to share them with us, fork the project, make your modifications, and then submit a pull request.

### Design Philosophy

In our project, we aim to achieve the following objectives:

- **Streamlined Cloud Development Setup**: We strive to provide a simple setup process for your cloud development environment, minimizing the number of commands required. Please note that in order to use our tools effectively, you should already have your cloud credentials set up.
- **Cross-Compatibility with Ubuntu**: While we welcome other Linux distributions, our primary goal is to ensure compatibility with at least 8 different working Ubuntu images.
- **Windows and MacOS Compatibility**: Our project aims to be compatible with both Windows and MacOS environments, allowing users to work seamlessly across different operating systems.
- **Cloud Provider Compatibility**: We want our tools to be compatible with major cloud providers such as AWS, Google, and Azure, enabling you to work with your preferred provider.
- **Consistent User Experience**: We value a single user experience across different platforms, ensuring a seamless and intuitive workflow.
- **Multiple Version Handling**: Our project supports multiple versions of specific tools, enabling flexibility and accommodating diverse requirements.
- **Security and Shared Responsibility**: We prioritize security and advocate for shared responsibility when it comes to managing cloud keys. It is crucial to maintain a secure and controlled environment.
- **Preference for Essential Tools**: We prefer tools that are essential for productive development across all platforms, extending beyond the scope of simple installation commands. Additionally, we encourage leveraging Docker or related technologies for easy accessibility.
- **Readable documentation**: Our project emphasizes the importance of clear and concise documentation. By providing well-documented resources, we aim to enhance understanding and assist users in utilizing our tools effectively.
