---
categories:
- guide
description: How to run commands using Rot
title: Run Commands
weight: 50
---

In this guide, we'll go over running commands with injected secrets via environment variables.

## Run Process

After Rot has been [initialized]({{< ref "/docs/guides/initialize-rot" >}}) and [populated with values]({{< ref "/docs/guides/manage-values" >}}), Rot can be used to wrap commands with secrets via environment variables.

The command to perform this is [`rot run <command>`]({{< ref "/docs/references/cli#run" >}}).  This command performs these actions:

- Decrypt the existing [Rot Private Key]({{< ref "/docs/references/cryptography#rot-private-key" >}}) using the [User Private Key]({{< ref "/docs/references/cryptography#user-private-key" >}}).
- Decrypt each value
- Run the specified command with an environment variable set for each value.  The environment variable name will be the values name.
