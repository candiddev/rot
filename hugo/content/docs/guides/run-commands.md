---
categories:
- guide
description: How to run commands using Rot
title: Run Commands
weight: 50
---

In this guide, we'll go over running commands with injected secrets via environment variables.

## Run Process

After Rot has been [initialized](../initialize-rot) and [populated with values](../manage-values), Rot can be used to wrap commands with secrets via environment variables.

The command to perform this is [`rot run <command>`](../../references/cli#run).  This command performs these actions:

- Decrypt the existing [Rot Private Key](../../references/cryptography#rot-private-key) using the [User Private Key](../../references/cryptography#user-private-key).
- Decrypt each value
- Run the specified command with an environment variable set for each value.  The environment variable name will be the values name.
