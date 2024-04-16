---
categories:
- guide
description: How to rekey secrets using Rot
title: Rekey Rot
weight: 40
---

In this guide, we'll go over rekeying a Rot Keyring.

## Rationale

Periodic rekeying of Rot's Keyrings is a non-destructive way to reduce the risk of storing your secrets using Rot.  This process generates new cryptographic keys for Values and Keyrings, and re-encrypts them using the Decrypt Public Keys.  Users can specific new Algorithms for use with this process, allowing easy upgrades to newer cryptography.

## Rekey Process

After Rot has been initialized, rekeying can be performed by running [`rot rekey`]({{< ref "/docs/references/cli#rekey" >}}).  This command performs these actions for the **current keyring**:

- Decrypt the existing [Keyring Private Key]({{< ref "/docs/references/cryptography#keyring-private-key" >}}) using the [Decrypt Private Key]({{< ref "/docs/references/cryptography#decrypt-private-key" >}}).
- Generate a new [Keyring Private Key]({{< ref "/docs/references/cryptography#keyring-private-key" >}}) and [Keyring Public Key]({{< ref "/docs/references/cryptography#keyring-public-key" >}}).
- Decrypt each Value and generate a new [Value Key]({{< ref "/docs/references/cryptography#value-key" >}})
- Encrypt the Value Key using the new [Keyring Public Key]({{< ref "/docs/references/cryptography#keyring-public-key" >}})
- Encrypt the new Keyring Private Key with each Decrypt Public Key and create a [signature]({{< ref "/docs/references/config#signature" >}}) of the Decrypt Public Key.
- Populate the [config](../../references/config) with the new keys and values, and save the configuration to [`configPath`]({{< ref "/docs/references/config#configpath" >}}).

By default, rekey will only affect the current [`keyring`]({{< ref "/docs/references/config#keyring" >}}).  Modify this configuration value permanently or temporarily to specify a different Keyring to rekey:

```bash
$ rot -x keyring=keyring2 rekey
```
