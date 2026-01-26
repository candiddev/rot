---
categories:
- reference
description: Reference documentation for Rot's configuration
title: Config
---

{{% snippet config_format Rot rot %}}

## Configuration Values

{{% snippet config_key "algorithms" %}}

Configuration values for setting the algorithms Rot will use by default.

{{% snippet config_key "algorithms_asymmetric" %}}

Specify the asymmetric encryption algorithm to use.  See {{< cli algorithms >}} for options.

**Default:** `"best"`

{{% snippet config_key "algorithms_pbkdf" %}}

Specify the Password Based Key Derivation Function (PBKDF) encryption algorithm to use.  See {{< cli algorithms >}} for options.

**Default:** `"best"`

{{% snippet config_key "algorithms_symmetric" %}}

Specify the symmetric encryption algorithm to use.  See {{< cli algorithms >}} for options.

**Default:** `"best"`

{{% snippet "config_cli" rot %}}

{{% snippet config_key "decryptKeys" %}}

A map of key names to key configurations.

**Default:** `{}`

### `decryptKeys_[name]_modified` {#decryptKeys_modified}

String, the UTC time the key was last modified.

**Default:** `""`

### `decryptKeys_[name]_privateKeySSS` {#decryptKeys_privateKeySSS}

List of encrypted values used by [Shamir Secret Sharing]({{< ref "/docs/guides/manage-keys#shamir-secret-sharing" >}}).

**Default:** `[]`

### `decryptKeys_[name]_privateKeys` {#decryptKeys_privateKeys}

A map of Keyring names to Private Keys and Signatures.

**Default:** `{}`


### `decryptKeys_[name]_privateKeys_[name]_signature` {#decryptKeys_privateKeys_signature}

String, a signature created by the [Keyring's Private Key](#keyringprivatekey) of the [Decrypt Public Key]({#decryptpublickey}).  This is used to prevent tampering of public keys.

**Default:** `""`

### `decryptKeys_[name]_privateKeys_[name]_publicKey` {#decryptKeys_privateKeys_publicKey}

String, the Decrypt Public Key.

**Default:** `""`

{{% snippet "config_httpClient" Rot %}}

{{% snippet "config_jsonnet" true %}}

{{% snippet config_key "keyrings" %}}

A map of Keyring names to Keyring details.

### `keyrings_[name]_encryptValues` {#keyrings_encryptValues}

Boolean, controls if the Keyring Values will have their Name and Meta properties encrypted too.  See [Manage Keyrings]({{< ref "/docs/guides/manage-keyrings" >}}) for more information.

**Default:** `false`

### `keyrings_[name]_privateKey` {#keyrings_privateKey}

String, the decrypted Keyring Private Key.  This is used to pass the Keyring Private Key via environment variables, such as from a external script that decrypts a [decryptKey.<keyring name>.privateKey](#decryptprivatekey) using a HSM or KMS.  Rot will never save this value to disk.

**Default:** `""`

### `keyrings_[name]_publicKey` {#keyrings_publicKey}

String, the Keyring Public Key.

**Default:** `""`

### `keyrings_[name]_values` {#keyrings_values}

A map of Value names to a map of Version Time and Value configurations.

{{< highlight json >}}
{
  "path/valueA": {
    "2024-04-24T00:00:00Z": {
      "key": "ecdhx25519hkdfsha256:MCowBQYDK2VwAyEAxNCdx0pHwQUh3f8QzhcYZ0qfmcvX1VF90iGfs+NWWUA=@xchacha20poly1305:A4f/zp076OopQaz8v1LOKqBLXH7QaXqSV190CaGwx0sAp3ah/ToFYdR
aAkobxojV4zCQtV7EQPwBrQ0rpNLLwNvzGNe8VNEV41KSPz9gcBACDZIz6cxpfCwZmz2HqvSTVyN+pDlix0Y=:s1cAADoakP",
      "meta": {
        "comment": "it's a value!"
      },
      "value": "xchacha20poly1305:D5hi10kxIiLH1URXJHlLscNeRBwfUR6q8YYvlRogAQfbReV/wErcskLebCsY3e0NJyX0YOlalEmMRSr+ncUbXqfyTYpBXWYoV/6qXWzMMlRQt5c0WIyaS/r9KoOa54IyWcHm32e
rgnkKo/0IdvXJHerUxusItlGhQns4G7ww+YlNSFDgyrq7UaZFROxAoqMxfe6n9h6HaSrXKxAn9bHdybV2ruUOOSrStwIVMyZdY97RYeyGYBZX5cqkcweE1HYoUO/r:cR5faafTQA",
    }
  }
}
{{< /highlight >}}

### `keyrings_[name]_values_[name]_key` {#keyrings_values_key}

String, the Rot Public Key encrypted symmetric key, used to encrypt/decrypt the [`value`](#value)

**Default:** `""`

### `keyrings_[name]_values_[name]_meta` {#keyrings_values_meta}

A map of strings containing various metadata about the Value.

{{< highlight json >}}
{
  "comment": "a comment!"
}
{{< /highlight >}}

Some special purpose meta keys:

- `_` Changes the name of the environment variable for the value in [run]({{% ref "/docs/guides/run-commands" %}}).
- `_links_dst` A comma separated list of paths that will also be updated when this value changes.
- `_links_src` The path that updates this value.
- `publicKey` {{% cli value-add-private %}} will store the public key here.

**Default:** `{}`

### `keyrings_[name]_values_[name]_value` {#keyrings_values_value}

String, a value encrypted using the [`key`](#key).

**Default:** `""`

### `keyrings_[name]_valuesEncrypted` {#keyrings_valuesEncrypted}

A map of encrypted Value Names to Value properties.

{{< highlight json >}}
{
  "xchacha20poly1305:6WnHSGlNLOYqiyGb1TGr/R3rb2mQFroSU7NyM4smsehUhnSPvb6yoXn7DAo=:fJveNApidl": {
    "key": "ecdhx25519hkdfsha256:MCowBQYDK2VwAyEAcQ+04/QsBhzONfYGq/99IlHWVeQV5Y+7h7lBKmwPz20=@xchacha20poly1305:+YC8IlIbpzH2Qs5GBtIId1gB+V3+ehMHdkugB+ARHUHYEC1ciEckP9VMaqcVDfywmUd78Mf23Jnn/G/mEDnn341lJlYgq7fgvW7TFd2fBtNloxCDBrs6JZryoVaIn5DnVeM9x1C7v/A=:JsZ8NTCNYS",
    "meta": {},
    "modified": "2024-04-25T00:00:00Z",
    "value": "xchacha20poly1305:bidJE4tsmHA4BrR58VWncRLXAeMD8kssOguy8kb8Yt3lNTlSFlrnMERxFOLPckE=:fJveNApidl"
  }
}
{{< /highlight >}}

### `keyrings_[name]_valuesEncrypted_[name]_key` {#keyrings_valuesEncrypted_key}

String, the Rot Public Key encrypted symmetric key, used to encrypt/decrypt the [`value`](#value)

**Default:** `""`

### `keyrings_[name]_valuesEncrypted_[name]_meta` {#keyrings_valuesEncrypted_meta}

A map of encrypted string keys to encrypted string values.

**Default:** `{}`

### `keyrings_[name]_valuesEncrypted_[name]_modified` {#keyrings_valuesEncrypted_modified}

String, the encrypted date the value was modified.

**Default:** `{}`

### `keyrings_[name]_valuesEncrypted_[name]_value` {#keyrings_valuesEncrypted_value}

String, a value encrypted using the [`key`](#valuesencryptedkey).

**Default:** `""`

{{% snippet config_key "keys" %}}

A list of strings containing Decrypt Private Keys.  This is mostly used to pass Decrypt Private Keys via environment variables.  Rot will never save this value to disk.

**Default:** `""`

{{% snippet config_key "keyPath" %}}

String, the path to a file containing Decrypt Private Keys, one per line.  If a filename without a path is specified, Rot will search parent directories for the filename and use the first one found.

**Default:** `".rot-keys"`

{{% snippet config_licenseKey Rot %}}

{{% snippet config_key "unmask" %}}

A list of Value names to unmask.

**Default:** `[]`

{{% snippet config_key "version" %}}

String, the version of the Rot configuration.

**Default:** `"<current Rot version>"`
