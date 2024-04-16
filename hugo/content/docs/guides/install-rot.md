---
categories:
- guide
description: How to install Rot
title: Install Rot
weight: 10
---

Installing Rot depends on how you want to run it.  Rot is available as a [binary](#binary) or a [container](#container).

## Binary

Rot binaries are available on [GitHub](https://github.com/candiddev/rot/releases).

{{< tabpane text=true >}}
{{< tab header="Linux amd64" >}}
{{< highlight bash >}}
curl -L https://github.com/candiddev/rot/releases/latest/download/rot_linux_amd64.gz -O | gzip -d > rot
chmod +x rot
sha256sum rot | grep $(curl -L https://github.com/candiddev/rot/releases/latest/download/rot_linux_amd64.sha256)
{{< /highlight >}}
{{< /tab >}}

{{< tab header="Linux arm" >}}
{{< highlight bash >}}
curl -L https://github.com/candiddev/rot/releases/latest/download/rot_linux_arm.gz -O | gzip -d > rot
chmod +x rot
sha256sum rot | grep $(curl -L https://github.com/candiddev/rot/releases/latest/download/rot_linux_arm.sha256)
{{< /highlight >}}
{{< /tab >}}

{{< tab header="Linux arm64" >}}
{{< highlight bash >}}
curl -L https://github.com/candiddev/rot/releases/latest/download/rot_linux_arm64.gz -O | gzip -d > rot
chmod +x rot
sha256sum rot | grep $(curl -L https://github.com/candiddev/rot/releases/latest/download/rot_linux_arm64.sha256)
{{< /tab >}}
{{< /highlight >}}

{{< tab header="macOS amd64" >}}
{{< highlight bash >}}
curl -L https://github.com/candiddev/rot/releases/latest/download/rot_darwin_amd64.gz -O | gzip -d > rot
chmod +x rot
sha256sum rot | grep $(curl -L https://github.com/candiddev/rot/releases/latest/download/rot_darwin_amd64.sha256)
{{< /highlight >}}
{{< /tab >}}

{{< tab header="macOS arm64" >}}
{{< highlight bash >}}
curl -L https://github.com/candiddev/rot/releases/latest/download/rot_darwin_arm64.gz -O | gzip -d > rot
chmod +x rot
sha256sum rot | grep $(curl -L https://github.com/candiddev/rot/releases/latest/download/rot_darwin_arm64.sha256)
{{< /highlight >}}
{{< /tab >}}
{{< /tabpane >}}

{{% alert title="Updating Rot" color="info" %}}
Rot can be updated by replacing the binary with the latest version.
{{% /alert %}}

## Container

Rot containers are available on [GitHub](https://github.com/candiddev/rot/pkgs/container/rot).

You can create an alias to run rot as a container:

{{< highlight bash >}}
alias rot='docker run -u $(id -u):$(id -g) -it --rm -v $(pwd):/work -w /work ghcr.io/candiddev/rot:latest'
{{< /highlight >}}

## SBOM

Rot ships with a Software Bill of Materials (SBOM) manifest generated using Cyclonedx.  The manifest is available with the other release assets at https://github.com/candiddev/rot/releases.
