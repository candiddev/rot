---
categories:
- guide
description: How to manage keyrings using Rot
title: Manage Keyrings
weight: 20
---

In this guide, we'll go over adding and removing keyrings in a Rot configuration.

## Adding Keyrings

After Rot has been initialized, keyrings can be added using [`rot add-keyring`]({{< ref "/docs/references/cli#add-keyring" >}}).  This command requires a name for the new keyring and a list of [Decrypt Keys]({{< ref "/docs/references/config#decryptkeys" >}}) to grant access to the new Keyring.  This command performs these actions

- Geenrate Private and Public Keys for the Keyring
- Encrypt the Keyring Private Key with the Decrypt Public Keys specified and create [signatures]({{< ref "/docs/references/config#signature" >}}) of the Decrypt Public Key.

Repeating this command can be used to grant additional access to the Keyring.  Alternatively, you can run [`rot add-keyprv`]({{< ref "/docs/references/cli#add-keyprv" >}}) to add Decrypt Keys individually to a Keyring, and [`rot remove-keyprv`]({{< ref "/docs/references/cli#remove-keyprv" >}}) to remove Decrypt Keys from a Keyring.

## Removing Keyrings

Keyrings can be removed by running [`rot remove-keyring`]({{< ref "/docs/references/cli#remove-keyring" >}}), or editing the configuration and removing the keyring and all references.  Decrypt Keys will be removed from the Keyring as well.  Rot will log errors if it discovers Decrypt Keys with access to unknown Keyrings.
