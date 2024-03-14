---
author: Mike
date: 2024-02-06
description: Release notes for Rot v2024.02.
tags:
  - release
title: "What's New in Rot: v2024.02"
type: blog
---

{{< rot-release version="2024.02" >}}

## Features

### Base64

Rot now includes base64 CLI helpers to convert between standard, url encoded, raw and padded base64 strings.  Visit [`rot base64`]({{< ref "/docs/references/cli#base64" >}}) for more information.

### JWTs

Rot can now create JWTs.  Visit [Generate JWTs]({{< ref "/docs/guides/generate-jwts" >}}) for more information.

### SSH

Rot can now create SSH keys and certificates.  Visit [Generate SSH]({{< ref "/docs/guides/generate-ssh" >}}) for more information.

### Signatures

Rot can now create cryptographic signatures.  Visit [Generate Signatures]({{< ref "/docs/guides/generate-signatures" >}}) for more information.

## Enhancements

- Added [`show-certificate`]({{< ref "/docs/references/cli#show-certificate" >}}) to display X.509 certificate details.

## Removals

- Removed [`generate-value`] to generate random values and add them to the configuration, its functionality has been moved into [`add-value`].
