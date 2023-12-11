---
categories:
- reference
description: Reference documentation for Rot's CLI
title: CLI
---

## Arguments

Arguments must be entered before commands.

### `-c [path]`

Path to the JSON/Jsonnet [configuration file](../config).

### `-f [format]`

Set log format (human, kv, raw, default: human).

### `-l [level]`

Set minimum log level (none, debug, info, error, default: info).

### `-n`

Disable colored log output.

### `-x [key=value]`

Set config key=value (can be provided multiple times)

## Commands

### `add-key [key name] [public key, default: generate a PBKDF-protected asymmetric key]` {#add-key}

Add a key to a configuration.  See [Manage Keys](../../guides/manage-keys) for more information.

### `add-value [name] [comment, default: ""] [delimiter, default: \n]` {#add-value}

Add a value to a configuration.  See [Manage Values](../../guides/manage-values) for more information.

### `decrypt [value]` {#decrypt}

Perform ad-hoc decryption of a value using the User Private Keys.

### `encrypt [recipient key] [delimiter, default: \n]` {#encrypt}

Perform ad-hoc encryption of a value using the recipient's key.

### `generate-keys [name] [algorithm, default: best]` {#generate-keys}

Generate ad-hoc cryptographic keys.

### `generate-value [name] [length, default=20] [comment]` {#generate-value}

Generate a random string of the length provided, encrypt it, and add it to the configuration.  See [Manage Values](../../guides/manage-values) for more information.

### `init [initial key name] [initial public key, default: generate a PBKDF symmetric key]` {#init}

Initialize a new Rot configuration.  See [Initialize Rot](../../guides/initialize-rot) for more information.

### `jq [jq query options]` {#jq}

Query JSON from stdin using jq.  Supports standard JQ queries, and the `-r` flag to render raw values.

### `rekey`

Rekey a Rot configuration.  See [Rekey Rot](../../guides/rekey-rot) for more information.

### `remove-key [name]` {#remove-key}

Remove a key from a Rot configuration.  See [Manage Keys](../../guides/manage-keys) for more information.

### `remove-value [name]` {#remove-value}

Remove a value from a Rot configuration.  See [Manage Values](../../guides/manage-values) for more information.

### `run [command]` {#run}

Run a command and inject secrets into it via environment variables.  See [Run Commands](../../guides/run-commands) for more information.

### `show-algorithms`

Show algorithms Rot understands

### `show-config`

Show the rendered config from all sources (file, environment variables, and command line arguments).

### `show-value [name]`

Show a decrypted value from the Rot configuration.  See [Manage Values](../../guides/manage-values) for more information.

### `version`

Print the current version of Rot.
