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
  userPrv[User Private Key]
  userPub[User Public Key]
  rotPrv[Rot Private Key]
  rotPub[Rot Public Key]
  valuePrv[Value Symmetric Key]
  value[Value]

  userPub -- Encrypts --> rotPrv
  rotPub -- Encrypts --> valuePrv
  valuePrv -- Encrypts/Decrypts --> value

  rotPrv -- Signs --> userPub
  userPrv -- Decrypts --> rotPrv
  rotPrv -- Decrypts --> valuePrv
```

With this setup, users can encrypt **Values** without decrypting the **Rot Private Key**.  Additionally, users can easily rekey all of the **Values** using the registered **User Public Keys**.  Rot will also check the signature of all **User Public Keys**, ensuring no one can tamper or add new keys out of band.
