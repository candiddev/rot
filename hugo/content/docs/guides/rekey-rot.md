---
categories:
- guide
description: How to rekey secrets using Rot
title: Rekey Rot
weight: 40
---

In this guide, we'll go over rekeying the Rot configuration.

## Rationale

Periodic rekeys of Rot's configuration is a non-destructive way to reduce the risk of storing your secrets using Rot.  This process generates new cryptographic keys for Values and Rot, and re-encrypts them using the User Public Keys.  Users can specific new Algorithms for use with this process, allowing easy upgrades to newer cryptography.

## Rekey Process

After Rot has been initialized, rekeying can be performed by running [`rot rekey`](../../references/cli#rekey).  This command performs these actions:

- Decrypt the existing [Rot Private Key](../../references/cryptography#rot-private-key) using the [User Private Key](../../references/cryptography#user-private-key).
- Generate a new [Rot Private Key](../../references/cryptography#rot-private-key) and [Rot Public Key](../../references/cryptography#rot-public-key).
- Decrypt each Value and generate a new [Value Key](../../references/cryptography#value-key)
- Encrypt the Value Key using the new [Rot Public Key](../../references/cryptography#rot-public-key)
- Encrypt the new Rot Private Key with each User Public Key and create a [signature](../../references/config#signature) of the User Public Key.
- Populate the [config](../../references/config) with the new keys and values, and save the configuration to [`configPath`](../../references/config#configpath).
