---
categories:
- guide
description: How to manage keys using Rot
title: Manage Keys
weight: 20
---

In this guide, we'll go over adding and removing keys in a Rot configuration.

## Adding Keys

After Rot has been initialized, keys can be added using [`rot add-key`](../../references/cli#add-key).  This command requires a name for a new key, and can optionally be provided with an existing public key.  This command performs these actions:

- Decrypt the existing [Rot Private Key](../../references/cryptography#rot-private-key) using the [User Private Key](../../references/cryptography#user-private-key).
- If a public key is not specified, generate a [User Private Key](../../references/cryptography#user-private-key) and [User Public Key](../../references/cryptography#user-public-key).  Rot will also save the User Private Key to the [keyPath](../../references/config#keypath), defaulting to `.rot-keys`
- Encrypt the Rot Private Key with the User Public Key and create a [signature](../../references/config#signature) of the User Public Key.
- Populate the [decryptKeys](../../references/config#decryptKeys) with the new key, and save the configuration to [`configPath`](../../references/config#configpath).

## Removing Keys

Keys can be removed by running [`rot remove-key`](../../references/cli#remove-key), or editing the configuration and removing the key.
