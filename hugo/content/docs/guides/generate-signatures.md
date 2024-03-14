---
categories:
- guide
description: How to create signatures using Rot
title: Generate Signatures
weight: 50
---

In this guide, we'll go over managing signatures using Rot.

## Signatures Introduction

Cryptographic signatures are signed hashes of a message created by a private key.  A public key will verify these hashes against a message, ensuring the message has not been modified and is signed by the underlying private key.

## Add Private Keys

You'll need to generate a private key for every certificate, including the CA.  The easiest way to do this is using [`rot add-pk`]({{< ref "/docs/references/cli#add-pk" >}}) (encrypting the keys into Rot) or [`rot gen-key`]({{< ref "/docs/references/cli#gen-key" >}}) (printing the keys to stdout).

Rot will store the public key in the comment of the encrypted value, we can grab the public key from the comment when we generate certificates.

## Create a signature

You can generate a signature using a private key with [`rot gen-sig`]({{< ref "/docs/references/cli#gen-sig" >}}).  This command generates a cryptographic hash from a private for the data provided as an argument or from stdin.

It supports the following flags:

- `-s`: Display just the signature

Example usage:

```bash
$ rot add-pk signer
$ rot gen-sig signer HelloWorld
ed25519:C6bcvoOqW+sROArW5EZqoG0W+4/sAKTQzqzsOueGgPBIDpc2uenu4TIXVp2eG2PCzZsavbjlxLaAnWfRZR/6CA==:signer
```

## Signature Format

Rot creates signatures with this format:

`<hash>:<signature>:<key id>`

- `hash` is the cryptographic hash that generated the signature.  Typically this will be `sha256` or `ed25519`.
- `signature` is the cryptographic signature, base64 standard encoded with padding.
- `key id` is the ID of the key used to generate the private key.

## Converting Signature to Base64 URL

Some applications require a different base64 encoding, such as JWTs (though [Rot can generate those too]({{< ref "/docs/guides/generate-jwts" >}})).  Rot has a base64 utility that can switch the signature format, [`rot base64`]({{< ref "/docs/references/cli#base64" >}}):

```bash
$ rot add-pk signer
$ rot gen-sig -s signer HelloWorld | rot base64 -d - | rot base64 -u -r -
C6bcvoOqW-sROArW5EZqoG0W-4_sAKTQzqzsOueGgPBIDpc2uenu4TIXVp2eG2PCzZsavbjlxLaAnWfRZR_6CA
```


### Verify Signature

Rot can verify signatures using [`rot verify-sig`]({{< ref "/docs/references/cli#verify-sig" >}}):

```bash
$ rot add-pk goodSigner
$ rot add-pk badSigner
$ rot gen-sig goodSigner HelloWorld
ed25519:s5AFx9ohilblnb1Fu5hlRiHB3qCYkL+rD3vyOF4xgG3CIbura2lcLnmNihVI4zmEyPtat4y4zq3rMv7o+kFDDw==:goodSigner
$ rot verify-sig goodSigner HelloWorld ed25519:s5AFx9ohilblnb1Fu5hlRiHB3qCYkL+rD3vyOF4xgG3CIbura2lcLnmNihVI4zmEyPtat4y4zq3rMv7o+kFDDw==:goodSigner
$ rot verify-sig badSigner HelloWorld ed25519:s5AFx9ohilblnb1Fu5hlRiHB3qCYkL+rD3vyOF4xgG3CIbura2lcLnmNihVI4zmEyPtat4y4zq3rMv7o+kFDDw==:goodSigner
ERROR error verifying signature against message
```
