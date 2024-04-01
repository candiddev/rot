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

### `-d`

Disable [Jsonnet]({{< ref "/docs/references/jsonnet" >}}) native functions that can reach external sources like `getPath` and `getRecord`.

### `-f [format]`

Set log format (human, kv, raw, default: human).

### `-l [level]`

Set minimum log level (none, debug, info, error, default: info).

### `-n`

Disable colored log output.

### `-x [key=value]`

Set config key=value (can be provided multiple times)

## Commands

### `add-key`

Add a key to a configuration.  See [Manage Keys]({{< ref "/docs/guides/manage-keys" >}}) for more information.

### `add-pk`

Generate and add a private key to a configuration with the specified name.

### `add-value`

Add a value to a configuration.  Can specify an optional length to have Rot randomly generate a value instead of prompting for it.  See [Manage Values]({{< ref "/docs/guides/manage-values" >}}) for more information.

### `autocomplete`

{{< autocomplete name="Rot" >}}

### `base64`

Import and export base64 strings in various encodings.  See command line for specific usage.

### `decrypt`

Perform ad-hoc decryption of a value using the User Private Keys.

### `encrypt`

Encrypt a value and print it to stdout without adding it to the config.  Can specify a recipient key to use asymmetric encryption.

### `gen-crt`

Generate X.509 certificates.  Visit [Generate Certificates]({{< ref "/docs/guides/generate-certificates" >}}) for more information.

### `gen-jwt`

Generate JSON Web Tokens (JWTs).  Visit [Generate JWTs]({{< ref "/docs/guides/generate-jwts" >}}) for more information.

### `gen-keys`

Generate ad-hoc cryptographic keys.

### `gen-sig`

Generate signatures.  Visit [Generate Signatures]({{< ref "/docs/guides/generate-signatures" >}}) for more information.

### `gen-ssh`

Generate SSH certificates.  Visit [Generate SSH]({{< ref "/docs/guides/generate-ssh" >}}) for more information.

### `init`

Initialize a new Rot configuration.  See [Initialize Rot]({{< ref "/docs/guides/initialize-rot" >}}) for more information.

### `jq`

Query JSON from stdin using jq.  Supports standard JQ queries, and the `-r` flag to render raw values.

### `pem`

Convert a Rot key to PEM or a PEM key to Rot.  Can specify an ID for the key when converting from PEM to Rot.

### `rekey`

Rekey a Rot configuration.  See [Rekey Rot]({{< ref "/docs/guides/rekey-rot" >}}) for more information.

### `remove-key`

Remove a key from a Rot configuration.  See [Manage Keys]({{< ref "/docs/guides/manage-keys" >}}) for more information.

### `remove-value`

Remove a value from a Rot configuration.  See [Manage Values]({{< ref "/docs/guides/manage-values" >}}) for more information.

### `run`

Run a command and inject secrets into it via environment variables.  See [Run Commands]({{< ref "/docs/guides/run-commands" >}}) for more information.  By default, any Value written to stderr/stdout will be masked with `***`.  Values can be unmasked using the [`unmask`]({{< ref "/docs/references/config#unmask" >}}) config.

### `show-alg`

Show algorithms Rot understands

### `show-config`

Show the rendered config from all sources (file, environment variables, and command line arguments).

### `show-crt`

Show the contents of an X.509 certificate and optionally verify it against a CA certificate.

### `show-jwt`

Show the contents of a JWT and optionally verify it against a public key.

### `show-keys`

Show the names of [decryptKeys]({{< ref "/docs/references/config#decryptKeys" >}}) in the configuration.

### `show-pk`

Show the public key for a User Private Key.  Takes a name of a key that it will lookup from [`keys`]({{< ref "/docs/references/config#keys" >}}) or [`keyPath`]({{< ref "/docs/references/config#keyPath" >}}).  Will return the public key of the first key found that matches `name`.

### `show-value`

Show a decrypted value from the Rot configuration.  See [Manage Values]({{< ref "/docs/guides/manage-values" >}}) for more information.

### `show-values`

Show the names of [values]({{< ref "/docs/references/config#values" >}}) in the configuration.

### `ssh`

Convert a Rot key to SSH or a SSH key to Rot.

### `verify-sig`

Verify the signature of a message and public key.  Will return an error/exit non-zero if the signature doesn't match, otherwise silently exits 0 on success.

### `version`

Print the current version of Rot.
