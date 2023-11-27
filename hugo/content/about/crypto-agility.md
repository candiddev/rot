---
categories:
- feature
description: Rot provides future proof cryptographic primitives
title: Crypto Agility
type: docs
---

Rot is designed to prevent cryptographic key "rot":

- Constantly decrypting leading to leakage
- Encryption algorithms become insecure

Instead of generating secrets once, Rot encourages companies to rekey secrets and reissue keys by making the process as easy as possible:

{{< highlight bash >}}
$ ./rot add-key server1 ed25519public:MCowBQYDK2VwAyEAAYkJzjQGb+4I7bfcaq6TnkI6nWJXolUdYSQDKSZIDZU=:AVvPeIzIHg
$ ./rot -x algorithms_asymmetric=rsa2048oaepsha256 -x algorithms_symmetric=aes128gcm rekey
{{< /highlight >}}
