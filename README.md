# Overview

A cloud development environment you can customize with a wide range of tools.

## Table of contents

- [Overview](#overview)
  - [Table of contents](#table-of-contents)
  - [Requirements](#requirements)
  - [How-to](#how-to)
    - [Connect](#connect)
  - [Tools](#tools)
    - [Cloud providers](#cloud-providers)
      - [AWS](#aws)
    - [Containerization](#containerization)
      - [Docker](#docker)
    - [IaC](#iac)
      - [Pulumi](#pulumi)
    - [Orchestration](#orchestration)
      - [Helm](#helm)
      - [Kubectl](#kubectl)
      - [Minikube](#minikube)
    - [Programming languages](#programming-languages)
      - [Dotnet](#dotnet)
      - [Go](#go)
      - [Node.js](#nodejs)
      - [Python](#python)
    - [Ruby](#ruby)
      - [Rust](#rust)
    - [Terminal](#terminal)
      - [Starship](#starship)
    - [Version control](#version-control)
      - [GitHub](#github)
  - [Extras](#extras)
    - [Unix aliases](#unix-aliases)
    - [Git aliases](#git-aliases)
  - [Q\&A](#qa)
  - [What are the recommended Ubuntu images?](#what-are-the-recommended-ubuntu-images)
    - [Ubuntu Jammy 22.04 AMD64](#ubuntu-jammy-2204-amd64)
  - [How to SSH into an instance?](#how-to-ssh-into-an-instance)
  - [How do I fix the broad permissions error when trying to ssh to my instance on Windows from Powershell?](#how-do-i-fix-the-broad-permissions-error-when-trying-to-ssh-to-my-instance-on-windows-from-powershell)
  - [How secure are my cloud credentials?](#how-secure-are-my-cloud-credentials)
    - [Access Limitations at Runtime](#access-limitations-at-runtime)
    - [Mitigating Impact in Case of Breach](#mitigating-impact-in-case-of-breach)
    - [Configuration Details](#configuration-details)
    - [Why not vaults?](#why-not-vaults)
    - [Why not secrets?](#why-not-secrets)
  - [How to remote ssh from VS Code?](#how-to-remote-ssh-from-vs-code)
  - [How to pick the right EC2 instance (AWS)](#how-to-pick-the-right-ec2-instance-aws)
  - [How to send files from my host to my instance?](#how-to-send-files-from-my-host-to-my-instance)
  - [Contributions](#contributions)
    - [Design Philosophy](#design-philosophy)

## Requirements

- Windows or MacOs
- OpenSSH is necessary for SSH access into the instance, specifically the client part.

## How-to

1. Download the binary according to your operative system and architecture.
2. Add the `kumo.exe` binary to your PATH
3. Create a new dir and a `kumo.config.yaml` file (For now we are only compatible with AWS)

    ```yaml
      Cloud: aws

      AWS:
        AccessKeyId: CUSTOM_VALUE
        SecretAccessKey: CUSTOM_VALUE
        IamProfile: CUSTOM_VALUE
        UserIds:
          - "CUSTOM_VALUE" #Quotes here are important
        Region: CUSTOM_VALUE
        EC2:
          Instance:
            Type: CUSTOM_VALUE
          Volume:
            Type: CUSTOM_VALUE
            Size: CUSTOM_VALUE

      AMI:
        Base:
          Filter: ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-20230516
          RootDeviceType: ebs
          VirtualizationType: hvm
          Owners:
            - "099720109477"
          User: ubuntu
        Name: CUSTOM_VALUE
        User: CUSTOM_VALUE
        Home: CUSTOM_VALUE
        Password: CUSTOM_VALUE
        Tools:
          - CUSTOM_VALUES
          - CUSTOM_VALUE
          - SEE REFERENCE BELOW

      Git:
        Username: CUSTOM_VALUE
        Email: CUSTOM_VALUE

      GitHub:
        PersonalAccessTokenClassic: xxxxxxxxxxxxx #optional

      Up:
        AmiId: CUSTOM_VALUE #Optional, in case you want to deploy an ami other than the last built
    ```

4. At the root of your project, run the following command:

    ```.exe
    kumo build
    kumo up
    kumo destroy
    ```

5. Boom! You've succesfully completed an entire workflow for building, deploying and destroying a development environment

### Connect

Assumming you've already ran `kumo build` and `kumo up`.

Once kumo completes, follow these steps:

  1. If you are on Windows, please refer to the [Q&A](#how-do-i-fix-the-broad-permissions-error-when-trying-to-ssh-to-my-instance-on-windows-from-powershell) section for guidance on how to fix the broad permissions error when trying to *SSH* to your instance from *PowerShell*.
  2. SSH into your instance with `ssh -i kumossh kumo`

If you want to remove your *AMI*, you can do so from your cloud management console. We follow the same philoshophy as *Packer*. You build it, you manage it.

## Tools

Add them to your `kumo.config.yaml` file as follows:

```yaml
    Tools:
      - aws
      - node_js
```

Please note that there are several tools listed below. Some of these tools may have specific requirements such as CPU, RAM, or disk space. Before building an AMI, it is recommended to consult their respective documentations for any specific requirements.

If you require a specific level of customization, you can find all the playbooks located at `~/kumo/packer/ansible/playbooks`

Make sure to add the appropriate tags to your environment file, and if needed, refer to the playbooks to customize your setup accordingly.

### Cloud providers

#### AWS

Add `aws`.

Ensure the following are present at the environment file:

```yaml
    AWS:
      AccessKeyId: CUSTOM_VALUE
      SecretAccessKey: CUSTOM_VALUE
```

### Containerization

#### Docker

Add `docker`.

### IaC

#### Pulumi

Add `pulumi`.

### Orchestration

#### Helm

Add `helm`.

#### Kubectl

Add `kubectl`.

#### Minikube

Add `docker` and `minikube`. The order matters as `docker` is a dependency.

### Programming languages

#### Dotnet

Add `dotnet`.

It installs `dotnet-sdk-7.0` under the hood. If you need to change this, go to `~/kumo/packer/ansible/playbooks/programming_languages/dotnet.yml` and change it on:

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

Add `go`.

Refer to the official [Go documentation](https://go.dev/doc/manage-install#installing-multiple) to learn how to manage multiple *Go* versions effectively. This guide provides detailed instructions on installing and managing multiple *Go* versions on your system.

Additionally, you can also check out this [Stack Overflow thread](https://stackoverflow.com/a/68087898/11941146), which explains how to alias multiple *Go* versions. This technique allows you to switch between different *Go* versions easily.

#### Node.js

Add `node_js`.

It installs [nvm](https://github.com/nvm-sh/nvm), which is a popular tool for managing multiple versions of *Node.js*.

Once *nvm* is installed, you can use it to install and manage multiple *Node.js* versions on your system. By using *nvm* commands, you can easily switch between different node versions based on your requirements.

#### Python

Add `python` to tools.

It installs [miniconda](https://docs.conda.io/en/latest/miniconda.html#linux-installers) which is a lightweight version of [anaconda](https://www.anaconda.com/). *Miniconda* allows you to manage multiple *Python* versions, virtual environments, and packages efficiently.

If you prefer a different *miniconda* version, you can refer to the download page for alternative options. Once you have chosen the appropriate version, make the following changes in the `~/kumo/packer/ansible/playbooks/programming_languages/python.yml` file:

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

Add `ruby`.

Manage versions with [rbenv](https://github.com/rbenv/rbenv).

#### Rust

Add `rust`.

Manage versions with [rustup-init](https://forge.rust-lang.org/infra/other-installation-methods.html).

### Terminal

#### Starship

Add `starship`.

For more information, please consult its [official documentation](https://starship.rs/).

### Version control

#### GitHub

Add `github` to tools and automate the process of adding *SSH* keys.

1. Update your environment file `.env` by adding the following line:

    ```yaml
      GitHub:
        PersonalAccessTokenClassic: CUSTOM_VALUE
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

If you're looking for some handy *Unix* aliases, you're in luck! We have included a collection of useful aliases that you can use. To costumize them, navigate to the `~/kumo/packer/ansible/playbooks/base/linux_aliases.yml file`

### Git aliases

We haven't forgotten about *Git*! In fact, we have also included some helpful *Git* aliases for you to use. To costumize them, go to the `~/kumo/packer/ansible/playbooks/base/git.yml` file. These aliases will make your *Git* commands more efficient and easier to remember.

## Q&A

## What are the recommended Ubuntu images?

When choosing an image on any cloud, it is important to consider various factors. While technically any available image can be picked, some may require additional considerations. If you have successfully tried a new image, we encourage you to open an issue or initiate a pull request.

For now, here is a list of recommended images:

### Ubuntu Jammy 22.04 AMD64

```yaml
    AMI:
      Base:
        Filter: ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-20230516
        RootDeviceType: ebs
        VirtualizationType: hvm
        Owners:
          - "099720109477"
        User: ubuntu
```

By selecting the appropriate image based on the provided details, you can ensure compatibility and meet your requirements when working on AWS.

## How to SSH into an instance?

1. Install the OpenSSH client on your local machine.
2. Open the command prompt or terminal and run the following command:

    ```bash
    ssh -F kumokey kumo
    ```

3. When prompted to add your instance URL to the list of known hosts, type `yes` and press enter.

## How do I fix the broad permissions error when trying to ssh to my instance on Windows from Powershell?

In you `kumossh` generated after running `kumo up`. Look for the identity file path:

```s
  Host kumo
    HostName 35.89.173.162
    IdentityFile "some/path" # <-- This
    User dev
    Port 22
    StrictHostKeyChecking no
    PasswordAuthentication no
    IdentitiesOnly yes
    LogLevel error
```

Run the following command from a PowerShell admin shell. You only need to do this once unless you removed the SSH key file.

Make sure the username is the same that will be sshing into the ec2 instance, otherwise change it to the appropiate one.

```powershell
icacls "some/path" /inheritance:r /grant:r "$($env:USERNAME):(R,W,D)"
```

## How secure are my cloud credentials?

The approach we took to secure your credentials is based on architectural decisions rather than heavy reliance on encryption and decryption.

### Access Limitations at Runtime

- To limit access to your instances, we control access based on your local IP and security groups.
- So it is recommended that you give the least amount of permissions to your cloud credentials.
- This also simplifies our workflow, as we do not have to depend on third-party vaults.

### Mitigating Impact in Case of Breach

- In the unlikely event that an attacker gains access to your cloud keys, you can restrict their impact by granting only the necessary permissions for a specific task.
- We encourage you to set up alerts or other monitoring mechanisms to provide early detection of any unauthorized activity.
- Remember, it is important not to use admin nor root level credentials to further enhance security.

### Configuration Details

Please note that the information provided above specifically accounts for placing the values inside a cloud instance using *Ansible*. It does not cover how *Packer* utilizes these credentials for initial communication with a cloud provider.

### Why not vaults?

We had the option to use tools such as *Ansible Vault* or *Hashicorp Vault* for securing credentials. However, we decided against them due to the following reasons:

- **Complexity of Local Installation**: Using *Ansible Vault* would require this tool to be locally installed, which could increase complexity, especially considering cross-compatibility issues on *Windows* and *MacOS*.
- **Cost**: *Hashicorp Vault* is not free, and even the open-source version adds complexity that would only be necessary for a small portion of the utility of this tool.

By avoiding these tools, we aim to simplify the setup and maintenance process, ensuring a smoother experience for both contributors and users.

### Why not secrets?

We prioritize a simple setup process for your convenience. Therefore, we stuck to the following principles:

- **Minimize External Dependencies:** Having to set up secrets on a 3rd party app right from the start can hinder the development experience.
- **Focus on Cloud Provider Credentials:** By limiting the permissions of your cloud provider credentials, you remain in control of ensuring they have the correct permissions. This puts the responsibility squarely on your own shoulders.

Our aim is to strike a balance between simplicity and security, providing you with a streamlined and efficient development experience while maintaining control over the permissions granted to your cloud provider credentials.

## How to remote ssh from VS Code?

  1. Ensure that you have the [Remote - SSH extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-ssh) installed in your VS Code. You can find and install it from the Extensions view (`Ctrl+Shift+X` or View -> Extensions).
  2. Once the extension is installed, click on the Extensions view button and search for "Remote - SSH". Click on the "Install" button to install the extension.
  3. Next, open the Command Palette in VS Code by pressing `Ctrl+Shift+P` (or View -> Command Palette). Type "Remote-SSH: Open Configuration File" and select the option from the list. This will open the SSH config file in the editor.
  4. Copy and paste the fields from the `kumossh` file generated after running `kumo up`:

      ```s
      Host cloud_dev
        HostName ip-of-your-instance
        IdentityFile "path/to/private-key"
        User picked-user
        Port 22
        StrictHostKeyChecking no
        PasswordAuthentication no
        IdentitiesOnly yes
        LogLevel INFO
      ```

  5. Save the SSH config file after adding the configuration.
  6. Go back to the Command Palette (`Ctrl+Shift+P`) and type "Remote-SSH: Connect to Host". You should see the list of configured SSH hosts from the SSH config file. Refresh if needed.
  7. Select the desired host from the list. It will attempt to establish an SSH connection to the remote machine using the provided configuration.
  8. If everything is set up correctly, a new window will open in VS Code connected to the remote machine via SSH.

## How to pick the right EC2 instance (AWS)

To help you make a decision based on factors such as cost and development needs, here are some tools that you can use:

- [AWS Pricing Calculator](https://calculator.aws/#/)
- [EC2 Instance Types](https://aws.amazon.com/ec2/instance-types/)
- [EC2 docs](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html)

## How to send files from my host to my instance?

```bash
scp -F ./kumossh /path/to/local/file1 /file2
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
- **Go Zen**: Throughtout the initial prototying phase we notice that the more we try to go agaist the goish way of doing things, the more redesigning was needed. Things like immutability, functional oriented, monads, cloning, etc (which are not by default in go). Long story short, accept *Go* by what it is and don't try to make it behave like another language unless you want to shoot yourself in the foot.
