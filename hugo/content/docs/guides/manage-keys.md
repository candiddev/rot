---
categories:
- guide
description: How to manage keys using Rot
title: Manage Keys
weight: 20
---

In this guide, we'll go over adding and removing keys in a Rot configuration.

## Adding Keys

After Rot has been initialized, keys can be added using [`rot add-keypub`]({{< ref "/docs/references/cli#add-keypub" >}}).  This command requires a name for a new key, and can optionally be provided with an existing public key.  If a public key is not specified, Rot will generate a [Decrypt Private Key]({{< ref "/docs/references/cryptography#decrypt-private-key" >}}) and [Decrypt Public Key]({{< ref "/docs/references/cryptography#decrypt-public-key" >}}).  Rot will also save the Decrypt Private Key to the [keyPath]({{< ref "/docs/references/config#keypath" >}}), defaulting to `.rot-keys`

By default, keys added to Rot will not have any access to Keyrings.  Access must be explicitly granted to Keyrings by existing Keyring users.  See [Manage Keyrings]({{< ref "/docs/guides/manage-keyrings" >}}) for more information.

## Removing Keys

Keys can be removed by running [`rot remove-keypub`]({{< ref "/docs/references/cli#remove-keypub" >}}), or editing the configuration and removing the key.
