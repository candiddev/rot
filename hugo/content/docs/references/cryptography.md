---
categories:
- reference
description: Reference documentation for Rot's cryptography
title: Cryptography
---

## Format

Rot formats keys and values like this:

`<encryption/algorithm>:<ciphertext/key>:<key id>`

- `encryption` is the Encryption used for an encrypted value
- `algorithm` is the underlying Algorithm for a key
- `ciphertext` is the encrypted value
- `key` is an unencrypted key
- `key id` is a free form string used to identify the key.

Some keys may be wrapped using a Key Derivation Function, these keys have this format:

`<kdf>:<kdf inputs>@<encryption>:<ciphertext>:<key id>`

An example KDF encrypted value looks like this:

`ecdhx25519:MCowBQYDK2VwAyEA8NpvTJLvgNNVhuy5NcN35hOzSxRQrDtoCXKMubIY1PM=@xchacha20poly1305:dfRf6/LfPzdd/LHjPDDtQOBnhUWXDI2BZMKjn1DeY7E0XMLCCtyhWBnF8w0qIgjaDMD+FT49ziKfqPxZaT+vxQCjXLD8QNToU+DXZmUhqSFP3kIrGhwAYrU/X2eHHR3NurkwZoo8eXZwwAoX5HLzNfEm8cmU6Ud5vsuR3QPmv0oe8A==:EcpHiptzqr`

### Private Key Format

Rot uses either EC or Ed25519 private keys, in PKCS #8 DER form, standard base64 encoded.

A example private key looks lke this:

`ed25519private:MC4CAQAwBQYDK2VwBCIEINr4+TmCQY+rhiEvBJnYhntrfiI2DXfZBx5PQcfHXC2I:alice`

### Public Key Format

Rot uses EC or Ed25519 public keys, in PKIX DER form, standard base64 encoded. 

A public key looks like this:

`ed25519public:MCowBQYDK2VwAyEArBbmsC9/nzSGTRqYfBRr2gju+pL6XlO6moJOTU+6flE=:EcpHiptzqr`

## Library

All Rot cryptographic functions are provided by the Go standard library.  Rot uses abstractions provided by [`cryptolib`](https://github.com/candiddev/shared/tree/main/go/cryptolib).

## Keys

Rot uses multiple sets of cryptographic keys to protect secrets.  Throughout this guide, we refer to a few of them by name:

### Rot Private Key

This is the asymmetric key used to decrypt all [Value Keys](#value-key) and sign all [User Public Keys](#user-public-keys).

### Rot Public Key

This is the asymmetric key used to encrypt all [Value Keys](#value-key).

### User Private Key

This is the asymmetric key used to decrypt the [Rot Private Key](#rot-private-key).

### User Public Key

This is the asymmetric key used to encrypt the [Rot Private Key](#rot-private-key).  A signature of this key is created using the [Rot Private Key](#rot-private-key) and verified by the [Rot Public Key](#rot-public-key) to prevent tampering.

### Value Key

This is the symmetric key used to encrypt and decrypt secrets.  It is encrypted using the [Rot Public Key](#rot-public-key)
