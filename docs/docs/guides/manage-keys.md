---
categories:
- guide
description: How to manage keys using Rot
title: Manage Keys
weight: 20
---

In this guide, we'll go over adding and removing keys in a Rot configuration.

## Adding Keys

After Rot has been initialized, keys can be added using {{% cli key-add-pub %}}.  This command requires a name for a new key, and can optionally be provided with an existing public key.  If a public key is not specified, Rot will generate a [Decrypt Private Key]({{< ref "/docs/references/cryptography#decrypt-private-key" >}}) and [Decrypt Public Key]({{< ref "/docs/references/cryptography#decrypt-public-key" >}}).  Rot will also save the Decrypt Private Key to the {{% config keyPath keypath %}}, defaulting to `.rot-keys`

By default, keys added to Rot will not have any access to Keyrings.  Access must be explicitly granted to Keyrings by existing Keyring users.  See [Manage Keyrings]({{< ref "/docs/guides/manage-keyrings" >}}) for more information.

### Shamir Secret Sharing 

{{% alert title="License Required" color="warning" %}}
This requires an [Unlimited License]({{< ref "/pricing" >}})
{{% /alert %}}

Rot can generate Decrypt Keys using [Shamir Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing).  This algorithm will split the Decrypt Key into a configurable number of Key Shares, each encrypted using User Public Key or a Password-based Key Derivation Function (PBKDF, like Argon2id).  The algorithm then requires a configurable number of Key Shares (called the Threshold) to recover the Decrypt Key.

Shamir Secret Sharing is a great way to protect highly sensitive Keyrings.  You could create a Keyring with only one Decrypt Key that uses Shamir Secret Sharing split amongst trusted individuals.  The only way to decrypt Values within the Keyring is to have all individuals provide their credentials to Rot.

All Key Shares are persisted within Rot, so the Users only need to remember the password they used to encrypt their Key Share, or have access to their Private Key to decrypt the Key Share.

## Removing Keys

Keys can be removed by running {{% cli key-pub-del %}}, or editing the configuration and removing the key.
