---
categories:
- reference
description: Reference documentation for Rot's CLI
title: CLI
---

{{% snippet cli_arguments %}}

{{% snippet cli_commands Rot %}}

### `alg`

Show algorithms Rot understands

{{% cli_autocomplete %}}

### `base64`

Import and export base64 strings in various encodings.  See command line for specific usage.

{{% snippet cli_config %}}

### `decrypt`

Perform ad-hoc decryption of a value.

### `decrypt-keys`

A convenience command to decrypt and source all encrypted keys to avoid typing them in constantly.  Best used within scripts.

{{% snippet cli_docs %}}

### `edit`

Interactively edit a Keyring's values in an editor like VSCode or vi.  Values will be displayed, decrypted, in Jsonnet, and any changes (add/change/delete) will be committed to the keyring upon saving.

### `encrypt`

Encrypt a value and print it to stdout without adding it to the config.  Can specify a recipient key to use asymmetric encryption.

{{% snippet cli_eula Rot %}}

### `hash-new`

Generate a hash of a file or string.

### `hash-ver`

Verify the hash of a file or string.  Will return an error/exit non-zero if the hash doesn't match, otherwise silently exits 0 on success.

### `init`

Initialize a new Rot configuration.  See [Initialize Rot]({{< ref "/docs/guides/initialize-rot" >}}) for more information.

{{% snippet cli_jq %}}

### `jwt-dis`

Show the contents of a JWT and optionally verify it against a public key.

### `jwt-new`

Generate JSON Web Tokens (JWTs).  Visit [Generate JWTs]({{< ref "/docs/guides/generate-jwts" >}}) for more information.

### `key-add-prv`

Add an existing Decrypt Key to a Keyring.  See [Manage Keys]({{< ref "/docs/guides/manage-keys" >}}) for more information.

### `key-add-pub`

Add a new or existing key to Decrypt Keys.  See [Manage Keys]({{< ref "/docs/guides/manage-keys" >}}) for more information.

### `key-del-prv`

Remove a Decrypt Key from a Keyring.  See [Manage Keys]({{< ref "/docs/guides/manage-keys" >}}) for more information.

### `key-del-pub`

Remove a Decrypt Key from Rot.  See [Manage Keys]({{< ref "/docs/guides/manage-keys" >}}) for more information.

### `key-dis`

Show the details of a {{% config decryptKey decryptkeys %}}.

### `key-dis-pub`

Display the public key of a private key.

### `key-lst`

List {{% config decryptKeys decryptkeys %}} in a configuration.

### `key-new`

Generate ad-hoc cryptographic keys.

### `keyring-add`

Add a new or modify an existing Keyring.  See [Manage Keyrings]({{< ref "/docs/guides/manage-keyrings" >}}) for more information.
Generate and add a private key to a configuration with the specified name.

### `keyring-del`

Remove a Keyring from Rot.  See [Manage Keyrings]({{< ref "/docs/guides/manage-keyrings" >}}) for more information.

### `keyring-lst`

Show the names of {{% config keyrings %}} in the configuration, optionally showing `decryptKeys` that can access them.

### `keyring-rekey`

Rekey a Keyring configuration.  See [Manage Keyrings]({{< ref "/docs/guides/manage-keyrings" >}}) for more information.

### `lic-add`

Add or replace the Rot license key.

### `pem`

Convert a Rot key to PEM or a PEM key to Rot.  Can specify an ID for the key when converting from PEM to Rot.

### `rand`

Generate random strings.

### `run`

Run a command and inject secrets into it via environment variables.  See [Run Commands]({{< ref "/docs/guides/run-commands" >}}) for more information.  By default, any Value written to stderr/stdout will be masked with `***`.  Values can be unmasked using the {{% config unmask %}} config.

### `sig-new`

Generate signatures.  Visit [Generate Signatures]({{< ref "/docs/guides/generate-signatures" >}}) for more information.

### `sig-ver`

Verify the signature of a message and public key.  Will return an error/exit non-zero if the signature doesn't match, otherwise silently exits 0 on success.

### `ssh-dis`

Convert a Rot key to SSH or a SSH key to Rot.

### `ssh-new`

Generate SSH certificates.  Visit [Generate SSH]({{< ref "/docs/guides/generate-ssh" >}}) for more information.

### `val-add`

Add a value to a configuration.  Can specify an optional length to have Rot randomly generate a value instead of prompting for it.  See [Manage Values]({{< ref "/docs/guides/manage-values" >}}) for more information.


### `val-add-prv`

Generate and add a private key to a configuration with the specified name.

### `val-cpy`

Copy a value between Keyrings.

### `val-del`

Remove a value from a Rot configuration.  See [Manage Values]({{< ref "/docs/guides/manage-values" >}}) for more information.

### `val-dis`

Show a decrypted value from the Rot configuration.  See [Manage Values]({{< ref "/docs/guides/manage-values" >}}) for more information.

### `val-lnk`

Link a {{% config keyrings_values value %}} to other values.  When the value is modified, the changes will be copied to the linked values.

### `val-lst`

Show the names of {{% config keyrings_values values %}} in a Keyring.

### `val-mov`

Move a value between Keyrings.

### `val-set-meta`

Modify metadata for an existing Value.

{{% cli_version %}}

### `x509-dis`

Show the contents of an X.509 certificate and optionally verify it against a CA certificate.

### `x509-new`

Generate X.509 certificates.  Visit [Generate Certificates]({{< ref "/docs/guides/generate-certificates" >}}) for more information.
