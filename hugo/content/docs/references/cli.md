---
categories:
- reference
description: Reference documentation for Rot's CLI
title: CLI
---

## Arguments

Arguments must be entered before commands.

### `-c [path]`

Path to the JSON/Jsonnet [configuration file]({{< ref "/docs/references/config" >}}).

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

Add a key to a configuration.  See [Manage Keys]({{< ref "/docs/guides/manage-keys" >}}) for more information.

### `add-private-key [name]` {#add-private-key}

Generate and add a private key to Rot with the specified name.

### `add-value [-d delimiter, default: \n] [name] [comment, default: ""]` {#add-value}

Add a value to a configuration.  See [Manage Values]({{< ref "/docs/guides/manage-values" >}}) for more information.

### `decrypt [value]` {#decrypt}

Perform ad-hoc decryption of a value using the User Private Keys.

### `encrypt [-d delimiter, default: \n] [recipient key, optional]` {#encrypt}

Encrypt a value and print it to stdout without adding it to the config.  Can specify a recipient key to use asymmetric encryption.

### `generate-certificate [-c, create CA] [-d hostname, add DNS hostname] [-e expiration in seconds] [-eu extended key usage] [-i address, add IP address] [-ku key usage] [-n common name] [private key value, name, or - for stdin] [public key] [ca certificate or path]` {#generate-certificate}

Generate X.509 certificates.  Visit [Generating Certificates]({{< ref "/docs/guides/generate-certificates" >}}) for more information.

### `generate-keys [-a algorithm] [name]` {#generate-keys}

Generate ad-hoc cryptographic keys.

### `generate-value [-l length, default 20] [name] [comment]` {#generate-value}

Generate a random string of the length provided, encrypt it, and add it to the configuration.  See [Manage Values]({{< ref "/docs/guides/manage-values" >}}) for more information.

### `init [initial key name] [initial public key, default: generate a PBKDF symmetric key]` {#init}

Initialize a new Rot configuration.  See [Initialize Rot]({{< ref "/docs/guides/initialize-rot" >}}) for more information.

### `jq [jq query options]` {#jq}

Query JSON from stdin using jq.  Supports standard JQ queries, and the `-r` flag to render raw values.

### `pem [-i id] [key, or - for stdin]` {#pem}

Convert a Rot key to PEM or a PEM key to a Rot key.  Can specify an ID for the key when converting from PEM to Rot.

### `rekey`

Rekey a Rot configuration.  See [Rekey Rot]({{< ref "/docs/guides/rekey-rot" >}}) for more information.

### `remove-key [name]` {#remove-key}

Remove a key from a Rot configuration.  See [Manage Keys]({{< ref "/docs/guides/manage-keys" >}}) for more information.

### `remove-value [name]` {#remove-value}

Remove a value from a Rot configuration.  See [Manage Values]({{< ref "/docs/guides/manage-values" >}}) for more information.

### `run [command]` {#run}

Run a command and inject secrets into it via environment variables.  See [Run Commands]({{< ref "/docs/guides/run-commands" >}}) for more information.  By default, any Value written to stderr/stdout will be masked with `***`.  Values can be unmasked using the [`unmask`]({{< ref "/docs/references/config#unmask" >}}) config.

### `show-algorithms`

Show algorithms Rot understands

### `show-config`

Show the rendered config from all sources (file, environment variables, and command line arguments).

### `show-keys`

Show the names of [decryptKeys]({{< ref "/docs/references/config#decryptKeys" >}}) in the configuration.

### `show-private-key`

Show the decrypted Rot Private Key.  Useful for piping into [`rot encrypt`](#encrypt) using a temporary key or combined with the [`privateKey`]({{< ref "/docs/references/config#privatekey" >}}) configuration to provide containers and other tools access to a configuration without compromising User Private Keys.

### `show-public-key [name]`

Show the public key for a User Private Key.  Takes a name of a key that it will lookup from [`keys`]({{< ref "/docs/references/config#keys" >}}) or [`keyPath`]({{< ref "/docs/references/config#keyPath" >}}).  Will return the public key of the first key found that matches `name`.

### `show-value [-c, show comment only] [-v, show value only] [name]`

Show a decrypted value from the Rot configuration.  See [Manage Values]({{< ref "/docs/guides/manage-values" >}}) for more information.

### `show-values`

Show the names of [values]({{< ref "/docs/references/config#values" >}}) in the configuration.

### `version`

Print the current version of Rot.
