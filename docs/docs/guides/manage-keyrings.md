---
categories:
- guide
description: How to manage Keyrings using Rot
title: Manage Keyrings
weight: 20
---

In this guide, we'll go over adding and removing Keyrings in a Rot configuration.

## Adding Keyrings

After Rot has been initialized, Keyrings can be added using {{% cli keyring-add %}}.  This command requires a name for the new keyring and a list of {{% config decryptKeys decryptkeys %}} to grant access to the new Keyring.  This command performs these actions

- Geenrate private and public keys for the Keyring.
- Encrypt the Keyring Private Key with the Decrypt Public Keys specified and create the {{% config decryptKeys_privateKeys_signature signature %}} of the Decrypt Public Key.

Repeating this command can be used to grant additional access to the Keyring.  Alternatively, you can run {{% cli key-add-private %}} to add Decrypt Keys individually to a Keyring, and {{% cli key-delete-private %}} to remove Decrypt Keys from a Keyring.

### Encrypted Value Names and Metadata

{{% alert title="License Required" color="warning" %}}
This requires an [Unlimited License]({{< ref "/pricing" >}})
{{% /alert %}}

By default, Value names, modified dates, and metadata are stored in plaintext within a Keyring's {{% config values %}}.  This is extremely convenient when using version control to see what has changed, however it does leak details about your secrets.

Keyrings can be configured to encrypt this data, however you'll need to be able to decrypt the Keyring to view the names and metadata within it.  Run `rot keyring-set -e <your keyring> && rot keyring-rekey <your keyring>` to change an existing Keyring to encrypted, or `rot keyring-add -e <new name>` to create a new, encrypted Keyring.

## Removing Keyrings

Keyrings can be removed by running [`rot remove-keyring`]({{< ref "/docs/references/cli#remove-keyring" >}}), or editing the configuration and removing the Keyring and all references.  Decrypt Keys will be removed from the Keyring as well.  Rot will log errors if it discovers Decrypt Keys with access to unknown Keyrings.

## Rekeying Keyrings

Periodic rekeying of Rot's Keyrings is a non-destructive way to reduce the risk of storing your secrets using Rot.  This process generates new cryptographic keys for Values and Keyrings, and re-encrypts them using the Decrypt Public Keys.  Users can specific new Algorithms for use with this process, allowing easy upgrades to newer cryptography.

### Rekey Process

After Rot has been initialized, rekeying can be performed on a per-Keyring basis by running {{% cli keyring-rekey %}}.  This command performs these actions:

- Decrypt the existing [Keyring Private Key]({{< ref "/docs/references/cryptography#keyring-private-key" >}}) using the [Decrypt Private Key]({{< ref "/docs/references/cryptography#decrypt-private-key" >}}).
- Generate a new [Keyring Private Key]({{< ref "/docs/references/cryptography#keyring-private-key" >}}) and [Keyring Public Key]({{< ref "/docs/references/cryptography#keyring-public-key" >}}).
- Decrypt each Value and generate a new [Value Key]({{< ref "/docs/references/cryptography#value-key" >}})
- Encrypt the Value Key using the new [Keyring Public Key]({{< ref "/docs/references/cryptography#keyring-public-key" >}})
- Encrypt the new Keyring Private Key with each Decrypt Public Key and create a {{% config decryptKeys_privateKeys_signature signature %}} of the Decrypt Public Key.
- Populate the [config](../../references/config) with the new keys and values, and save the configuration to {{% config configPath configpath %}}.
