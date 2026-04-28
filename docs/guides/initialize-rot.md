---
categories:
- guide
description: How to initialize Rot
title: Initialize Rot
weight: 20
---

In this guide, we'll go over initializing a new Rot configuration.

## Initialization Process

Out of the box, Rot doesn't know about any keys or values.  Rot can perform ad-hoc encryption and decryption using {{% cli encrypt %}} and {{% cli decrypt %}}, but long term storage and sharing of secrets must be done using a shared configuration.

Rot is initialized using {{% cli init %}}.  This command will create a new configuration within the current directory (or wherever {{% config configPath configpath %}} points to).  A Keyring name and initial public key can be provided, as well as enabling [Encrypted Vlaue Names and Metadata]({{< ref "/docs/guides/manage-keyrings#encrypted-value-names-and-metadata" >}})otherwise Rot will generate these.  This command performs these actions:

- Create a new Keyring
- Generate the [Keyring Private Key]({{< ref "/docs/references/cryptography#keyring-private-key" >}}) and [Keyring Public Key]({{< ref "/docs/references/cryptography#keyring-public-key" >}}).
- If a public key is not specified, generate a [Decrypt Private Key]({{< ref "/docs/references/cryptography#decrypt-private-key" >}}) and [Decrypt Public Key]({{< ref "/docs/references/cryptography#decrypt-public-key" >}}).  Rot will also save the User Private Key to the {{% config keyPath keypath %}}, defaulting to `.rot-keys`
- Geenrate private and public keys for the Keyring.
- Encrypt the Keyring Private Key with the Decrypt Public Key and create a {{% config decryptKeys_privateKeys_signature signature %}} of the Decrypt Public Key.
- Populate the [config values]({{< ref "/docs/references/config" >}}) and save the configuration to {{% config cli_configPath %}}.

If `rot init` is ran again, a warning will appear to prevent accidental overwriting.
