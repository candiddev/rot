---
author: Candid Development
date: 2024-04-15
description: Release notes for Rot v2024.04.
tags:
  - release
title: "What's New in Rot: v2024.04"
type: blog
---

## Features

### Edit Mode

Values can now be edited within a text editor for easier modifications.  Visit [Manage Values]({{< ref "/docs/guides/manage-values" >}}) for more information.

### Encrypted Value Names and Metadata

Keyring Value names and metadata can now be optionally encrypted to prevent leaking secret details.  Visit [Manage Keyrings]({{< ref "/docs/guides/manage-keyrings#encrypted-value-names-and-metadata" >}}) for more information.

### Keyrings

Rot now supports Keyrings: a way to logically separate your Keys and Values within one Rot configuration.  Visit [Manage Keyrings]({{< ref "/docs/guides/manage-keyrings" >}}) for more information.

### Shamir Secret Sharing

[Decrypt Private Keys]({{< ref "/docs/references/cryptography#decrypt-private-key" >}}) can now be split into secure key shares using Shamir Secret Sharing.  Visit [Manage Keys]({{< ref "/docs/guides/manage-keys#shamir-secret-sharing" >}}) for more infomration.

### Value Versions

Values will now have previous versions remembered.  You can retrieve the latest version, or a version from a specific point in time.

## Enhancements

- Added HKDF variants for ECDHX25519 and ECHP256 encryption.  This is considered the new best encryption, rekeying is recommended.
- Added {{% cli decrypt-keys %}} as a convenient way to cache keys without having to unlock them constantly.
- Added {{% cli val-cpy %}} to copy Values between Keyrings.
- Added Windows builds.
- Changed {{% cli run %}} to add all meta values as environment variables.  They will be appended to the Value's name, i.e. a value `TEST` with the meta key `publicKey=123` would have an environment variable `TEST_publicKey=123`
- Changed the Value lookup syntax.  All Values must reference the keyring: `rot val-dis mykeyring/a/path/to/value`.
- Changed Value comments to a key/value map called meta.  Visit [Manage Values]({{< ref "/docs/guides/manage-values#meta" >}}) for more information.
- `rot encrypt` and `rot decrypt` can now reference a path to a keyring/value.
- Updated the Rot CLI library for better paging navigation.
