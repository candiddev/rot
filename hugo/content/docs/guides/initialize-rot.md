---
categories:
- guide
description: How to initialize Rot
title: Initialize Rot
weight: 20
---

In this guide, we'll go over initializing a new Rot configuration.

## Initialization Process

Out of the box, Rot doesn't know about any keys or values.  Rot can perform ad-hoc encryption and decryption using [`rot encrypt`](../../references/cli#encrypt) and [`rot decrypt`](../../references/cli#decrypt), but long term storage and sharing of secrets must be done using a shared configuration.

Rot is initialized using [`rot init`](../../references/cli#init).  This command requires a name for a new key, and can optionally be provided with an existing public key.  This command performs these actions:

- Generate the [Rot Private Key](../../references/cryptography#rot-private-key) and [Rot Public Key](../../references/cryptography#rot-public-key).
- If a public key is not specified, generate a [User Private Key](../../references/cryptography#user-private-key) and [User Public Key](../../references/cryptography#user-public-key).  Rot will also save the User Private Key to the [keyPath](../../references/config#keypath), defaulting to `.rot-keys`
- Encrypt the Rot Private Key with the User Public Key and create a [signature](../../references/config#signature) of the User Public Key.
- Populate the [config values](../../references/config) and save the configuration to [`configPath`](../../references/config#configpath).

If `rot init` is ran again, a warning will appear to prevent accidental overwriting.
