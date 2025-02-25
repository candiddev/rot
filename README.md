# <img alt=logo src=rot.png width=40px> Rot

> Secure Secrets Management for the Modern Sysadmin

**Latest Release Notes: [v2025.02](https://rotx.dev/blog/whats-new-202502/)**

[:speech_balloon: Mailing Lists](https://rotx.dev/docs/references/mailing-lists/)\
[:book: Docs](https://rotx.dev/docs/)\
[:arrow_down: Download](https://rotx.dev/docs/guides/install-rot/)\
[:eyes: Open Source Details](https://candid.dev/open-source)

Rot is a command line (CLI) tool for managing cryptographic values.

Rot makes cryptography easy:

- Generate keys and values using current best encryption
- Version and track your secrets
- Secure secrets among individuals using Shamir Secret Sharing
- Rekey encrypted values to the latest encryption standards
- Share your encrypted values with other users and devices
- Perform one-way encryption for zero-knowledge secrets
- Run commands and scripts with encrypted values injected through environment variables
- Store your encrypted values securely in git with human-readable diffs
- Generate and view X.509 certificates and Certificate Authorities
- Generate and view JWTs
- Generate and verify signatures
- Generate SSH keys and certificates

Visit https://rotx.dev for more information.

## Why Rot?

Aside from the infamous [ROT13 Caeser cipher](https://en.wikipedia.org/wiki/ROT13), cryptographic keys have a tendency to "rot" rather quickly.  Frequent use of keys inevitably leads to leakage and/or compromise, and the underlying encryption algorithms may not be secure in the future.  Cryptographic material doesn't `age` well in general.
