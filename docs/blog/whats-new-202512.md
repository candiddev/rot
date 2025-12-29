---
author: Candid Development
date: 2025-12-01
description: Release notes for Rot v2025.12.
tags:
  - release
title: "What's New in Rot: v2025.12"
type: blog
---

## Deprecations

## Features

### Post-Quantum Cryptography

We added initial support for post-quantum cryptography via Ed25519/ML-KEM-768 hybrid encryption.  Rot can generate private and public keys (essential concatenated, base64 encoded Ed25519 keys and ML-KEM-768 decapsulation/encapsulation keys) for use with encryption only.  Signing and verifying are not supported, yet.  Additionally, this key format may change as ASN.1 DER support evolves for them--Rot will maintain compatibility with this initial key format, and be able to convert to a future, supported ASN.1 DER format in the future.

Rot will not default to using this encryption type for now.  You can set Rot to default to this format by setting {{% config algorithms_asymmetric %}} to `ecdhx25519mlkem768hkdfsha256`.

## Enhancements

- Update Go to 1.25.5

## Fixes

## Removals
