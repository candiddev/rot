---
categories:
- reference
description: Reference documentation for Etcha's configuration
title: Config
---

Etcha can be configured using a JSON/Jsonnet configuration file, environment variables, and [command line arguments](../cli#-x-keyvalue).  Configurations from one source will override previous sources, i.e. environment variables override configuration files, command line arguments override environment variables.

The configuration is divided into these sections:

- <a href="#algorithms">Algorithms</a>
- <a href="#cli">CLI</a>
- <a href="#decryptkeys">DecryptKeys</a>
- <a href="#keypath">KeyPath</a>
- <a href="#keys">Keys</a>
- <a href="#privatekey">PrivateKey</a>
- <a href="#publickey">PublicKey</a>
- <a href="#unmask">Unmask</a>
- <a href="#values">Values</a>

All configuration keys are camelCase.  Configuration values can be:

- **boolean** `true` or `false`
- **integer** `1`
- **list** `["a","b","c"]`
- **map** `{"key": {}}`
- **string** `"a"`

**For command line values**, every configuration key can be set using `-x <a_config_key1>="a value" -x <a_config_key2>="another value"`, i.e. `-x cli_logLevel=debug -x algorithms_asymmetric=best`.  Config values can also be set using JSON, i.e. `-x algorithms='{"asymmetric": "best"}'`

**For environment variables**, every configuration key can be set using `ROT_section_key=a value`, i.e. `ROT_cli_logLevel=debug`

**For configuration files**, they can be formatted using JSON or Jsonnet.  Rot will look for `rot.jsonnet` by default, ascending the directory tree to find it.  See [the Jsonnet reference](../jsonnet/) for more information.  **Configuration files are rendered at startup**, allowing you to use [dynamic Jsonnet functions](../jsonnet#native-functions) to dynamically alter the config, i.e.:

```
local getRecord(type, name, fallback=null) = std.native('getRecord')(type, name, fallback);
local verifyKey = getRecord('txt', 'server1.candid.dev');

{
  decryptKeys: [
    server1: {
      publicKey: verifyKey,
    },
  ]
}
```

You can view the rendered configuration by running [`rotshow-config`](../cli#show-config).

## `algorithms`

### `asymmetric`

Specify the asymmetric encryption algorithm to use.  See [`rot show-algorithms`](../cli#show-algorithms) for options.

### `pbkdf`

Specify the Password Based Key Derivation Function (PBKDF) encryption algorithm to use.  See [`rot show-algorithms`](../cli#show-algorithms) for options.

**Default:** `"best"`

### `symmetric`

Specify the symmetric encryption algorithm to use.  See [`rot show-algorithms`](../cli#show-algorithms) for options.

**Default:** `"best"`

## `cli`

### `configPath`

String, path to the configuration file.  If a filename without a path is specified, Rot will search parent directories for the filename and use the first one found.

**Default:** `"rot.jsonnet"`

### `logFormat`

String, log format to use for logging: human, kv, or raw.

**Default:** `"human"`

### `logLevel`

String, log level to use for logging: none, debug, info, or error.

**Default:** `"info"`

### `noColor`

Boolean, disables colored log output.

**Default:** `false`

## `decryptKeys`

A map of key names to key configurations.

**Default:** `{}`

### `modified`

String, the UTC time the key was last modified.

**Default:** `""`

### `privateKey` {#userprivatekey}

String, the Rot Private Key encrypted using the [User's Public Key](#userpublickey).

**Default:** `""`

### `publicKey` {#userpublickey}

String, the User's Public Key.

**Default:** `""`

### `signature`

String, a signature created by the [Rot Private Key](#rotprivatekey) of the User's Public Key.  This is used to prevent tampering of public keys.

**Default:** `""`

## `keys`

A list of strings containing User Private Keys.  This is mostly used to pass User Private Keys via environment variables.  Rot will never save this value to disk.

**Default:** `""`

## `keyPath`

String, the path to a file containing User Private Keys, one per line.  If a filename without a path is specified, Rot will search parent directories for the filename and use the first one found.

**Default:** `".rot-keys"`

## `privateKey`

String, the decrypted Rot Private Key.  This is used to pass the Rot Private Key via environment variables, such as from a external script that decrypts a [decryptKey.privateKey](#userprivatekey) using a HSM or KMS.  Rot will never save this value to disk.

**Default:** `""`

## `publicKey` {#rotpublickey}

String, the Rot Public Key.

**Default:** `""`

## `unmask`

A list of Value names to unmask.

**Default:** `[]`

## `values`

A map of value names to value configurations.

### `comment`

String, a user-supplied comment for the value.

**Default:** `""`

### `key`

String, the Rot Public Key encrypted symmetric key, used to encrypt/decrypt the [`value`](#value)

**Default:** `""`

### `modified`

String, the UTC time the value was last modified.

**Default:** `""`

### `value`

String, a value encrypted using the [`key`](#key).

**Default:** `""`
