---
categories:
- guide
description: How to run commands using Rot
title: Run Commands
weight: 40
---

In this guide, we'll go over running commands with injected secrets via environment variables.

## Run Process

After Rot has been [initialized]({{< ref "/docs/guides/initialize-rot" >}}) and [populated with values]({{< ref "/docs/guides/manage-values" >}}), Rot can be used to wrap commands with secrets via environment variables.

The command to perform this is [`rot run <command>`]({{< ref "/docs/references/cli#run" >}}).  This command performs these actions:

- Decrypt the existing [Rot Private Key]({{< ref "/docs/references/cryptography#rot-private-key" >}}) using the [User Private Key]({{< ref "/docs/references/cryptography#user-private-key" >}}).
- Decrypt each value
- Run the specified command with an environment variable set for each value.  The environment variable name will be the values name.

By default, Rot use Values in the current [`keyring`]({{< ref "/docs/references/config#keyring" >}}).  Modify this configuration value permanently or temporarily to specify a different Keyring:

```bash
$ rot -x keyring=keyring2 run <command>
```
