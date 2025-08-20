---
categories:
- explanation
description: How Rot wraps secrets
title: Secret Wrapping
weight: 10
---

Rot uses multiple layers of secrets wrapping to limit key-reuse and enable one-way encryption:

```mermaid
flowchart TD
  decryptPrv[Decrypt Private Key]
  decryptPub[Decrypt Public Key]
  keyringPrv[Keyring Private Key]
  keyringPub[Keyring Public Key]
  valuePrv[Value Symmetric Key]
  value[Value]

  decryptPub -- Encrypts --> keyringPrv
  keyringPub -- Encrypts --> valuePrv
  valuePrv -- Encrypts/Decrypts --> value

  keyringPrv -- Signs --> decryptPub
  decryptPrv -- Decrypts --> keyringPrv
  keyringPrv -- Decrypts --> valuePrv
```

With this setup, users can encrypt **Values** without decrypting the **Keyring Private Key**.  Additionally, users can easily rekey all of the **Values** using the registered **Decrypt Public Keys**.  Rot will also check the signature of all **Decrypt Public Keys**, ensuring no one can tamper or add new keys out of band.
