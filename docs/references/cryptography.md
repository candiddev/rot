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

Rot uses either EC, Ed25519, ML-DSA, or RSA private keys, in PKCS #8 DER form, standard base64 encoded.  Rot can also generate post-quantum hybrid private keys, like Ed25519/ML-KEM-768.  These keys are concatenated and standard base64 encoded.

Example private keys:

- `ed25519private:MC4CAQAwBQYDK2VwBCIEINr4+TmCQY+rhiEvBJnYhntrfiI2DXfZBx5PQcfHXC2I:alice`
- `ed25519mlkem768private:c18A9ENji5DoI3D5EmUYzrA3ElqhLOFiJPdRsUCWVlM+ES7z5D9HFmBs0lzGhqNx0IhJ/j+i8/4QsYdZ6wCTQh7QyJHhEOz0V5B4+V+nF0qmTgexKsg75h4XDw500xNKeudUDmWlOKHl5/e5tbvovV/FeWwxnB6ioDzWkf08L48=:bob`
- `mldsa44private:MDQCAQAwCwYJYIZIAWUDBAMRBCKAIEmPlXI6kH772Qb3FcEeVvDPWK9UADXJOvr7zZLIbm7u:435JAMvzeW`

### Public Key Format

Rot uses EC, Ed25519, ML-DSA, or RSA public keys, in PKIX DER form, standard base64 encoded.  Rot can also generate post-quantum hybrid public keys, like Ed25519/ML-KEM-768.  These keys are concatenated and standard base64 encoded.

Example public keys:

- `ed25519public:MCowBQYDK2VwAyEArBbmsC9/nzSGTRqYfBRr2gju+pL6XlO6moJOTU+6flE=:alice`
- `ed25519mlkem768public:PhEu8+Q/RxZgbNJcxoajcdCISf4/ovP+ELGHWesAk0L6GoyiYQtepVeDHMzyZ2SvrIJDNDqpoIIVWKWFowRuOxr+QWwWez0++m5+iTAZKoYoRTUijDQD5HL1qyzbklqOqUWt2UKJto2zMwAFQAYgcspu5HhGwqfAgJfDWFth+JNZbHyi/LVfRh8t10zimqgZx7FKfDMWETylIRAVerN4pLHaCqvyUJG71ySf9F8A4Xf/BUbelWMGEpyOUKIF2MyNYzUJa37p4Yd/sE87mEmyIWabEm3l2k2nuC03S7uCq1nb6HzRCabRcZqAKFa+Kwmawa/2+RryerkKuUHW4wjrcjxfqFenc1KZkU1rU5WIUBT/2oasWoSqNiMsoMKAenpZ9LFEaASrhZdC83rIU5XZ2y7hxmPbIy3tcpXUJRlYk3Oe1GVrpIAeG2+WGnoerK49MXOIlwP+YxhByZIv5TgPBHZNwjRn87oAyUhtK0B+nIyDZc6HQ6bVs6l+1SMWgZtEYszjF5j5NC756KLcinvl6Fv9CGQHlq3TJxOD5ST7BE3GeY5MKD2kUKtAMVkQ8h1kgVCduJSbyXU79RprtLWs4VTKs0nZtmoiGg/dFDrnYkG+yDIBppzvcXaLmBl/qK/HWLTz1jJCWw1LpgZpVb8jgMLRPLqPo8VEa6G7i8EirLNg5aLsRRdphRsxZkjkOAr8ELKBeAb12xg0ogmEQGCUCFrYOjqcXMEQEzkrZ7lxpM+7hxWtYGH/YjkurHLsKIhySrO9sEnWqGhpzCt/WW6boGRi2FzWvCaAaLpLIT8lFkhplIxrBTnvNXn1USRWMspXEjjtynx2oD9SNr8es6EHPJVp9WTYMcFb8Rko3LCdyxBY5H+TIb9LYwkhYilIBEzS5JR6MbRpi5S5iCqx43zY61ORozY4RpnZt8OVqUnozHIvZ3TmZ4OT+RQGBiW8KwHd1m/9MqA6w0sHYHZ1osgCA6hnMLFZm1qKhR4BoE97hV0DxKmmshVx6T1/qnCJBXA8Kc35ppA9CaGU1ClZkm7+t3DGsLmUSSLKDC4Md3Jfo7/YTHQMwcGAmWfxBnhtPJiiuFtadi2xWLFmEXvchcumRH9xNoD/+1l/Wb+5dVE1aYaVikoOI82BXBhpsiZKBndjopTnhKF5tozL2XJVuJqYoZD9mFI/HJ396qGkJRgTicz0UqcFSWj696qBYYWJErQtawcucbxAJxHNJDjCqig4/FBfoxNJSJZ9k14xwLbCKLFQ2qZ3sRnkMWOwijwAd4AnhkgMNZADIoo7sQIQkTOZmkfL0ypfFYg0xKtl5829hqoI9ZgO+7Pz9Gf+5w/jZyROLLtHSc7rF7JpW1OA4oa1TL+m9wrvbFZNMGJtejM452V1BKgjo4PuebPitbuORTIgimWCgxDNkX6yIqTGVGxuoVKoAJPNaSG8FoveZhOT+Ieo40k6c2g4g5UY2Vi3qY+0+jrQ1VbsZVV9A1PcN1615yeT8ZcxWi7AHEoc4IKRWldH20a8RcLc9TYHyl91IICNwJ/1x0uvg1nRpCK8+AsL2nqVx4G4VXOuyyrakUB9N3ToXzqTXqRl3TlZI88wc5iVWtLebrj31UCM1DkxIC5mVg==:bob`
- `mldsa44public:MIIFMjALBglghkgBZQMEAxEDggUhAKp6abP5FrBxIsagZBb5ZYCH4fHn5AdC4KLjDUFb6sCHqnpdGZU1cdAaDB5fVjw7OVWxYkV/wLmfPJbUbpvAfiPJbmtHO84oNAsJDRjysY67q1RO16ab2x4tz07kPb6ItX+IXMtK0xjLcRPK1tHgg+hedG8ohnuN4fPFZ/AjxJ7zUwWOYNvqzqEp8jhDM+ph4XIb2dw4cg3QyPY4JZsoh7BGY7VGa6bufwfv/cpq0+mJTBt/fzSIn04oOvxi8BWY909h7zL0LA54jee35utkEGPrfqUSI5/G80wHC4uUwCgD84Vg2tENQ+krk6mCUYZYPJ5DNxCkmWDou7dJc/rtV0jpPVRWGanLPOKScycHeeOQ4GyYzbTaDH6JBzOV7jQVVqD6xDOwHASzHk3y4q7vU/z7VYDom1tM+nVenjuItsmi8cZplWJ6u/dutYES+33ZqGd9kLjQIwlo6Oh+FxGsCFhQEurIXfdY/tH88Zco9LYyCMmo+DKHeObzIAqpL6XATfkWN6wv8YdkS3IqeSpOdnXCWqfpXV1r31UZfhdU8OSXqlkjm1mY7F5sDqWef4qSJNiOcE6gMhrj72omU/LQNd8ohnBZvZxwUbp+IWllqnkdeTIWkkCtHc1mawzF5IjMSrkDFjj44qR4htT1OW3UQVacy/npu4GrVIZPyiWfWvobxFiXxGnEkgeYsiwZeaJEV+Oi0JQJM9CH5Vc3zhEqoHASNtExHmKaKyl6EJPvBAAnYVhBlP5rt5QOITm7icwl0oJ9XE4MQIXOQH8tAxFwhd9ufb9TiLptlsuovcW+zJM02f2y/uVUBMqfABc4I5ZL77OD3LhNPXvyZWOhoJmJ8KxvxCkF7whEWiyfhr7HLLu3epLmGGouw+rcKFcpNP4vFMt6OlKGtB9xb22uWveb7FYWdLY2OFGRvkVY0/J6SgilRfXOqbRV72ZtsupUp3hKhwpYFw8pdSih/mFAAKus9i/coEzIGB3fPSdeXFeyeSWCYmSi2NrGNVKhG6EoCPHzwIcL09aF+YVlhPaFRhvzrd/B2SJUmWY6GaMUvSmXG4hTjf+52qSfc2IWpYdq5kC0bx4YOXTLtA2Y45XtEbkTjtM9FGIjyDzuQwJcuc6j876X2JtSvTh2s0BjBRrAn39iNfnFor1Vodt9zShwLjgfoeSHwKxucsQmscuf4Ro6ET6yPoGyJDlcTRRtYOe6PfyUKTKvNC8mWTsvbtPSvqd437iH37DyIB5/qfLH+f1GMNb1x5hnVgndBjvRo1tR6bWCJA1xaxze60G1Zqw7dX/iKBmx1+Q/elZLVfUPWeNmCg2EODWKVln1kaQEfCy96NJw9l6u46WXC/uKpw1jEpc+TKaAT8ilI7IGGfMn+8rR4WiL24bBKx2LnYbnplaJ7CYp/pDJBVQr9kmL0CRqc62jVHTIJpHxHZwLHc94IiRQwkNFbCBoRetP0G+eZRKE9U6eOPLjvPFJGdjVJ8sO2udcp+E+Q7A+8ov5ABs55lHHcDjeiWvWXMcIvNOxxpqNLJWHZXY9ToRrrGD2kCpL4sGfj2LOOxgZDhnukGS7HPW5w8PKqIgUbsbCMrLaW6r0gP4bHJ00T5gXk6tWP70db9koAYvcjEdfavxGG43Cu/v0lirXminutw4lBuFBsjoC2F7LWFSdsBTN3bvTtVULLKNNquMaFDIKugePuq8YhvPvny1Sg7u53f5F4uGv3AgqJoRyIoDR/fU7ZWPDYqCXUnCFGLQ=:435JAMvzeW`

{{% alert title="Post-Quantum Keys are Huge" color="info" %}}
Post-quantum cryptographic keys tend to be very large, especially with larger parameters.  Keep this in mind when thinking about distributing keys and the impact the added key size can have on latency.
{{% /alert %}}


## Library

All Rot cryptographic functions are provided by the Go standard library `golang.org/x/crypto`.

## Keys

Rot uses multiple sets of cryptographic keys to protect secrets.  Throughout this guide, we refer to a few of them by name:

### Decrypt Private Key

This is the asymmetric key used to decrypt the [Keyring Private Key](#keyring-private-key).  These are typically keys generated by end users and devices.

### Decrypt Public Key

This is the asymmetric key used to encrypt the [Keyring Private Key](#keyring-private-key).  A signature of this key is created using the [Keyring Private Key](#keyring-private-key) and verified by the [Keyring Public Key](#keyring-public-key) to prevent tampering.

### Keyring Private Key

This is the asymmetric key used to decrypt all [Value Keys](#value-key) and sign all [Decrypt Public Keys](#decrypt-public-key).

### Keyring Public Key

This is the asymmetric key used to encrypt all [Value Keys](#value-key).

### Value Key

This is the symmetric key used to encrypt and decrypt secrets.  It is encrypted using the [Keyring Public Key](#keyring-public-key)
