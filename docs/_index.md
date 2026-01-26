---
description: Your command-line toolbox for managing cryptographic values. Easily generate, version, and share secrets using Shamir Secret Sharing. Rekey to latest standards, inject secrets into environment variables, and store securely in git. Generate X.509 certs, JWTs, signatures, and SSH keys. Rot makes cryptography effortless.
title: Rot | Boring Cryptography Tooling for Humans
---

{{% blocks/section color="white" %}}
<h1 style="border-bottom: 2px solid var(--bs-gray)"><b>Secure Your Build and Deployment Processes with Rot</b></h1>
<h1>Boring Cryptography Tooling for Humans</h1>
<div style="align-items: center; display: flex; justify-content: center; padding-top: 40px; width 100%">
  <a class="button button--gray" href="/docs/guides/install-rot">Download</a>
</div>
{{% /blocks/section %}}

{{< blocks/section color="white" type=row >}}
{{% blocks/feature icon="fa-lock" title="Always Use the Best Algorithms" %}}
Generate keys and values using the latest and greatest cryptography standards. Sleep soundly knowing your secrets are safe.
{{% /blocks/feature %}}

{{% blocks/feature icon="fa-code-branch" title="Version Control Your Secrets" %}}
Store your secrets in version control, track changes and revert to previous versions with ease. No more scrambling to remember what went wrong.
{{% /blocks/feature %}}

{{% blocks/feature icon="fa-users" title="Share and Rotate Secrets with Ease" %}}
Securely share secrets among your team. No single point of failure, just peace of mind.
{{% /blocks/feature %}}

{{% blocks/feature icon="fa-clock" title="Achieve Crypto Agility" %}}
Rekey your encrypted values to the latest standards with minimal effort. Your secrets stay safe, even as technology evolves.
{{% /blocks/feature %}}
{{< /blocks/section >}}

{{< blocks/section color="white" height=full >}}
<h2 style="border-bottom: 2px solid var(--bs-gray)"><b>Effortless Collaboration</b></h2>
<h3>Share secrets securely within your team and control access with granular permissions.</h3>
{{< highlight bash >}}
# Initialize rot.jsonnet and a key for Alice
$ rot init mykeyring alice
New Password (empty string skips PBKDF):
Confirm Password (empty string skips PBKDF):

# Add Bob's public key
$ rot key-add-public bob ed25519public:MCowBQYDK2VwAyEAAYkJzjQGb+4I7bfcaq6TnkI6nWJXolUdYSQDKSZIDZU=:bob

# Grant Bob access to our Keyring
$ rot key-add-private bob mykeyring
{{< /highlight >}}

<h2 style="border-bottom: 2px solid var(--bs-gray); padding-top: 50px"><b>Keep It Safe, and Versioned</b></h2>
<h3>Values are added with zero-knowledge of the encryption key, and versioned by timestamp</h3>
{{< highlight bash >}}
# Add a Value
$ rot value-add mykeyring/prd/database/secret database=backend1
Value:

# Replace it
$ rot value-add mykeyring/prd/database/secret database=backend2
Value:

# List the Values
$ rot value-list mykeyring
{
  "mykeyring/prd/database/secret": [
    "2024-04-26T12:32:44Z",
    "2024-04-26T10:32:52Z"
  ]
}

# View it--need access to decrypt!
$ rot value-display mykeyring/prd/database/secret
Password for alice:
{
  "meta": {
    "database": "backend2"
  },
  "modified": "2024-04-26T12:32:44Z",
  "value": "hunter2"
}
{{< /highlight >}}

<h2 style="border-bottom: 2px solid var(--bs-gray); padding-top: 50px"><b>Convert Everything</b></h2>
<h3>Manage all kinds of cryptographic keys: PEM, SSH, X509, JWT, and more</h3>
{{< highlight bash >}}
# Generate a public/private keypair
$ rot key-new
{
  "privateKey": "ed25519private:MC4CAQAwBQYDK2VwBCIEIHZm1QM2sOFMWwJpmSptFkb2KsndJzebZZ8V8ogaMB9z:q5cnabsnVa",
  "publicKey": "ed25519public:MCowBQYDK2VwAyEAYes5wEXIHi1qaQNCNRKbqcgi3qQ07QAbSx2e7LgOhVo=:q5cnabsnVa"
}

# Convert it to PEM
$ rot pem ed25519public:MCowBQYDK2VwAyEAYes5wEXIHi1qaQNCNRKbqcgi3qQ07QAbSx2e7LgOhVo=:q5cnabsnVa
-----BEGIN PUBLIC KEY-----

MCowBQYDK2VwAyEAYes5wEXIHi1qaQNCNRKbqcgi3qQ07QAbSx2e7LgOhVo=
-----END PUBLIC KEY-----

# Convert it to SSH
$ rot ssh-display ed25519public:MCowBQYDK2VwAyEAYes5wEXIHi1qaQNCNRKbqcgi3qQ07QAbSx2e7LgOhVo=:q5cnabsnVa
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGHrOcBFyB4tamkDQjUSm6nIIt6kNO0AG0sdnuy4DoVa

# Generate an X.509 CA
$ rot x509-new -c -n "My CA" ed25519private:MC4CAQAwBQYDK2VwBCIEIHZm1QM2sOFMWwJpmSptFkb2KsndJzebZZ8V8ogaMB9z:q5cnabsnVa ed25519public:MCowBQYDK2VwAyEAYes5wEXIHi1qaQNCNRKbqcgi3qQ07QAbSx2e7LgOhVo=:q5cnabsnVa | rot x509-display -
...
  "CommonName": "My CA",
  "IsCA": true,
{{< /highlight >}}
{{< /blocks/section >}}
