---
categories:
- explanations
description: An overview of Rot's security
title: Security
---

## Bug Bounties

Rot does not yet have an established bug bounty program.  
Please [contact us](mailto:info@candid.dev?subject=Rot%20Bug) if you think you've found a bug or security issue with Rot.

## CVEs

Rot does not have any CVEs.  When a CVE is reported, it will be listed on this page.

## Code

Rot is developed using Go.  Here are some of the methods we use to help keep Rot's code free of vulnerabilities:

- **Auth Test Suite**: We use an extensive authentication and authorization test suite for every pull request and build.
- **Limit Third-Party Libraries**: We try and use as few third-party libraries as possible, and when we do select a third-party library, we review the codebase to ensure it's something we are comfortable maintaining.
- **Secure Software Supply Chain**: We require a clean `govulncheck` for every pull request and build.

## Cryptography

Rot leverages boring, secure cryptographic functions.  See [Cryptography]({{< ref "/docs/references/cryptography" >}}) for more information.
