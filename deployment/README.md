# Process Reporter Deployment

This repository contains an Ansible playbook and configuration to **deploy and run the `process-reporter` binary** across different operating systems — Linux, macOS, and Windows.

---

## What It Does

The playbook:

* Detects the OS and architecture of each host
* Downloads the appropriate binary from GitHub Releases
* Installs it in the appropriate location:

    * `/usr/local/bin/process-reporter` on Linux/macOS
    * `C:\bin\process-reporter.exe` on Windows
* Runs the binary with `--format csv`
* Shows the output

---

## Structure

| File                  | Description                                               |
|-----------------------|-----------------------------------------------------------|
| `ansible.cfg`         | Ansible configuration                                     |
| `deploy.sh`           | Bash helper to invoke Ansible with `sudo` pre-initialized |
| `hosts`               | Inventory file with host groups and connection settings   |
| `playbook-deploy.yml` | Main Ansible playbook performing the deployment           |

---

## Prerequisites

* Ansible installed
* For **Linux/macOS**:
    * SSH access or `local` connection
* For **Windows**:
    * WinRM enabled and configured
* Environment variables set:

```bash
export DEPLOY_PASSWORD='your_windows_password_here'
```


### Setting Up WinRM on Windows (Non-Production)

To enable Ansible management on a Windows host, **run this PowerShell script as Administrator** on the target machine:

```powershell
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
$url = "https://raw.githubusercontent.com/ansible/ansible-documentation/devel/examples/scripts/ConfigureRemotingForAnsible.ps1"
$file = "$env:temp\ConfigureRemotingForAnsible.ps1"

(New-Object -TypeName System.Net.WebClient).DownloadFile($url, $file)

powershell.exe -ExecutionPolicy ByPass -File $file
```

This script:

* Enables WinRM
* Sets up a listener for HTTP
* Adds firewall exceptions
* Grants the necessary user permissions

⚠️ This setup is suitable for **testing/development** — not production use.

---

## Running the Deployment

Use the helper script:

```bash
export ANSIBLE_PASSWORD=my_secret
./deploy.sh
```

Or limit deployment to specific host/group:

```bash
./deploy.sh -l windows
./deploy.sh -l linux
./deploy.sh -l macos
```

---


## Supported Platforms

| OS      | Architecture |
|---------|--------------|
| Linux   | amd64, arm64 |
| macOS   | amd64, arm64 |
| Windows | amd64, arm64 |

If the OS/arch combo isn't supported, the playbook fails with an explicit error.

---

## What the Playbook Does

For each host:

1. Detects the OS and architecture.
2. Validates if the platform is supported.
3. Downloads latest release from https://github.com/ldziedziul/process-reporter-go/releases
4. Installs the binary:
    * Linux/macOS: `/usr/local/bin`
    * Windows: `C:\bin`
5. Runs the binary: `process-reporter --format csv`
6. Displays the output

---
