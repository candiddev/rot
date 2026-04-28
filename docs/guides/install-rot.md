---
categories:
- guide
description: How to install Rot
title: Install Rot
weight: 10
---

Installing Rot depends on how you want to run it.  Rot is available as a [binary](#binary) or a [container](#container).

## Binary

Rot binaries are available for various architectures and operating systems:

{{% release %}}

{{% alert title="Updating Rot" color="primary" %}}
Rot can be updated by replacing the binary with the latest version.
{{% /alert %}}

## Container

Rot containers are available on [GitHub](https://github.com/candiddev/rot/pkgs/container/rot).

You can create an alias to run rot as a container:

{{< highlight bash >}}
alias rot='docker run -u $(id -u):$(id -g) -it --rm -v $(pwd):/work -w /work ghcr.io/candiddev/rot:latest'
{{< /highlight >}}

## SBOM

Rot ships with a Software Bill of Materials (SBOM) manifest generated using [CycloneDX](https://cyclonedx.org/).  The `.bom.json` manifest is available with the other [Binary Assets](#binary).


