---
categories:
- guide
description: How to initialize Rot
title: Initialize Rot
weight: 20
---

In this guide, we'll go over initializing a new Rot configuration.

## Initialization Process

Out of the box, Rot doesn't know about any keys or values.  Rot can perform ad-hoc encryption and decryption using [`rot encrypt`]({{< ref "/docs/references/cli#encrypt" >}}) and [`rot decrypt`]({{< ref "/docs/references/cli#decrypt" >}}), but long term storage and sharing of secrets must be done using a shared configuration.

Rot is initialized using [`rot init`]({{< ref "/docs/references/cli#init" >}}).  This command requires a name for a new key, and can optionally be provided with an existing public key.  This command performs these actions:

- Generate the [Rot Private Key]({{< ref "/docs/references/cryptography#rot-private-key" >}}) and [Rot Public Key]({{< ref "/docs/references/cryptography#rot-public-key" >}}).
- If a public key is not specified, generate a [User Private Key]({{< ref "/docs/references/cryptography#user-private-key" >}}) and [User Public Key]({{< ref "/docs/references/cryptography#user-public-key" >}}).  Rot will also save the User Private Key to the [keyPath]({{< ref "/docs/references/config#keypath" >}}), defaulting to `.rot-keys`
- Encrypt the Rot Private Key with the User Public Key and create a [signature]({{< ref "/docs/references/config#signature" >}}) of the User Public Key.
- Populate the [config values]({{< ref "/docs/references/config" >}}) and save the configuration to [`configPath`]({{< ref "/docs/references/config#configpath" >}}).

If `rot init` is ran again, a warning will appear to prevent accidental overwriting.
