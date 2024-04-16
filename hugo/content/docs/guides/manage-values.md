---
categories:
- guide
description: How to manage values using Rot
title: Manage Values
weight: 30
---

In this guide, we'll go over adding and removing values in a Rot Keyring.

## Adding Values

After Rot has been initialized, keys can be added using [`rot add-value`]({{< ref "/docs/references/cli#add-value" >}}).  This command requires a name for a new key, and can optionally be provided with a comment and delimiter for the value.  This command performs these actions:

- Generate a [Value Key]({{< ref "/docs/references/cryptography#value-key" >}})
- Encrypt the provided Value using the Value Key
- Encrypt the Value Key using the [Keyring Public Key]({{< ref "/docs/references/cryptography#keyring-public-key" >}})
- Populate the [values]({{< ref "/docs/references/config#decryptKeys" >}}) with the new value, and save the configuration to [`configPath`]({{< ref "/docs/references/config#configpath" >}}).

By default, Rot will generate Values in the current [`keyring`]({{< ref "/docs/references/config#keyring" >}}).  Modify this configuration value permanently or temporarily to specify a different Keyring:

```bash
$ rot -x keyring=keyring2 add-value
```

## Generating Values

Rot can generate random, cryptographically secure strings for you, instead of having you provide a value, via [`rot add-value`]({{< ref "/docs/references/cli#add-value" >}}).

## Removing Values

Values can be removed by running [`rot remove-value`]({{< ref "/docs/references/cli#remove-value" >}}), or editing the configuration and removing the value.
