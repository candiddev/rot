---
categories:
- guide
description: How to manage Trusted Platform Modules (TPMs) using Rot
title: Manage TPMs
weight: 20
---

In this guide, we'll go over managing TPMs using Rot.

## TPM Basics

Interacting with a TPM using Rot typically requires your user to be a member of the Linux `tss` group or root access.

### Hash Banks

TPMs contain data typically stored within a Hash Bank--a set of registers that all rely on the same hashing function.  Most TPM functions within Rot will target a specific Hash Bank on a TPM, by default it's `sha256`.

### Platform Configuration Registers (PCRs)

Hash Banks primarily store Platform Configuration Registers (PCRs), these are indexes numbered 0 - 23 and contain a hash digest using the Hash Bank's hashing function.  When a TPM first starts up, these values are all 0s.  As a system boots, various events will trigger a PCR Extend event, where an additional hash digest is added to the PCR.  The function looks something like this:

newPCR = sha256(oldPCRHash, newHash)

#### Getting PCR Values

Rot can retrieve all PCR values or a specific index using {{% cli tpm-pcr-get %}}:

```bash
$ rot tpm-pcr-get
{
  "00": "sha256:9a18f1edfc1007e534cacc1e84338aa57302f3fb2e332508012ac5886f192647",
  "01": "sha256:5257d7b1f9453f8e3dc294403765584afcc58e2d931f11f047ae7d23fc176604",
  "02": "sha256:ef674cd26ad78dd8240a4723b1a527de84218773adeaaabf68a9a71b64e1dba8",
  "03": "sha256:3d458cfe55cc03ea1f443f1562beec8df51c75e14a9fcf9a7234a13f198e7969",
  "04": "sha256:6da146705b029e0d75305d55ea45379045bfab730510772b8abee7c81290aaee",
  "05": "sha256:d8f6423c72a713a6905b09aa392bbace9328fa2aea6ca715955600441ddb95df",
  "06": "sha256:3d458cfe55cc03ea1f443f1562beec8df51c75e14a9fcf9a7234a13f198e7969",
  "07": "sha256:3ce1bc35f76916b16b1164186d0eb3a932f65e211100142394525a03f9bb3555",
  "08": "sha256:0000000000000000000000000000000000000000000000000000000000000000",
  "09": "sha256:e0097cc7960419323c260c971a41a40ab1b5b1c310d6265189ad93a0b351005e",
  "10": "sha256:4a5e312990d29715b6d59cca3d907be7e6109d0e263f1a5ed7cba7ed0c9ddc19",
  "11": "sha256:0000000000000000000000000000000000000000000000000000000000000000",
  "12": "sha256:ed66f7feea751d58058028c984f8117b6fe5aa6bb1f53a67969c806888b7dc07",
  "13": "sha256:0000000000000000000000000000000000000000000000000000000000000000",
  "14": "sha256:1dd2c3c78e629ad7ecc3789c06f7b41feaf3cc9a6d5ee4256e41c64c89910f22",
  "15": "sha256:0000000000000000000000000000000000000000000000000000000000000000",
  "16": "sha256:2a3c73e026a49579cdd2a644f845bcb8296420f35ca193a35ca65b075811283b",
  "17": "sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
  "18": "sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
  "19": "sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
  "20": "sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
  "21": "sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
  "22": "sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
  "23": "sha256:0000000000000000000000000000000000000000000000000000000000000000"
}
```

#### Predicting PCRs

These extend events are typically deterministic, so if you know the strings or files that are used by them, Rot can predict what the PCR value will be using {{% cli tpm-pcr-predict %}}:

```bash
$ rot tpm-pcr-predict hello world
sha256:98d128df384d428ffe76af3c0198ff1e8945ef71e741ba440bafff0510da8f22
```

#### Extending PCRs

Rot can extend PCRs using {{% cli tpm-pcr-extend %}}, allowing you to add your own additional hash digests to existing PCR values.  Typically, only PCR 15 and 16 can be extended:

```bash
$ rot tpm-pcr-extend 16 sha256:1dd2c3c78e629ad7ecc3789c06f7b41feaf3cc9a6d5ee4256e41c64c89910f22
```

#### Resetting PCRs

While experimenting with predicting and extending PCRs, it may be helpful to reset the PCR to 0 using {{% cli tpm-pcr-reset %}}.  This typically only works for PCR 16 (and 23):

```bash
$ rot tpm-pcr-reset 16
$ rot -p tpm-pcr-get 16
sha256:0000000000000000000000000000000000000000000000000000000000000000
```

### Event Log

On Linux, modifications to the TPM's PCRs will be written in the a binary event log (`/sys/kernel/security/tpm0/binary_bios_measurements`).  Rot can display these events using {{% cli tpm-event-log %}}:

```bash
$ rot event-log
...
  {
    "event": "EFI Platform Firmware Blob 2",
    "digest": "sha256:a5eea1abe74e31d37bfcad463e3711937c289086070b92809d005f9f32d6ec89",
    "pcr": 0
  },
  {
    "event": "EFI Variable Driver Config",
    "digest": "sha256:115aa827dbccfb44d216ad9ecfda56bdea620b860a94bed5b7a27bba1c4d02d8",
    "pcr": 7
  },
  {
    "event": "EFI Variable Driver Config",
    "digest": "sha256:dea7b80ab53a3daaa24d5cc46c64e1fa9ffd03739f90aadbd8c0867c4a5b4890",
    "pcr": 7
  },
...
```

## Encrypting and Decrypting Secrets

TPMs encrypt secrets by sealing the secret using cryptographic keys within the TPM that cannot be exported.  These keys are called Storage Root Keys (SRKs).  You can use predefined SRKs (some operating systems create them), or use transient SRKs that are deterministically created on-demand.  You can additionally add a policy around the unsealing process, such as requiring a specific set of PCRs to have a specific digest.  If the PCRs don't align, the value cannot be retrieved.

Sealing a value using ({{% cli tpm-seal-value %}}) will generate an encrypted value string that can be used to unseal and decrypt the value:

```bash
$ rot tpm-seal-value -r
tpm2:sha256@pAE4ACAALAAAAkgAgU9FlRLu0j4eQoyzWqY/ok9e/WbnIc3HtQD0jeFs+jcwAEAAgcYFHk4Y1eJtKiXWsWkkCZZ7ij2EpwO31eGAQbD/0PF4=@i7,11:ACD9jrJczAFT6HY9ox8pQioaqJj0CXaT0iSnaLEvkfJ08AAQIWPsNcaW811ufkVX9ji4SYDVi/uapoNetVRjiwtDyPbxbQxrXFgIfz2BaO/BewwxP26zJ53SKJ9hIM0nagK6yJeuFk9EIpMQYrAM+78bmYCu0ciAfi5mn9Qh/JcmATGJUPWRHhqpvVUZzccj/0AFzKDLjbnaD8mTuETQVACOJnI/OhBzkw==
```

In this example, Rot generated a random string (`-r`) and sealed it against the default/recommend PCRs, 7 and 11.

Lets unseal it using {{% cli tpm-unseal-value %}}:

```bash
$ rot tpm-unseal-value tpm2:sha256@pAE4ACAALAAAAkgAgU9FlRLu0j4eQoyzWqY/ok9e/WbnIc3HtQD0jeFs+jcwAEAAgcYFHk4Y1eJtKiXWsWkkCZZ7ij2EpwO31eGAQbD/0PF4=@i7,11:ACD9jrJczAFT6HY9ox8pQioaqJj0CXaT0iSnaLEvkfJ08AAQIWPsNcaW811ufkVX9ji4SYDVi/uapoNetVRjiwtDyPbxbQxrXFgIfz2BaO/BewwxP26zJ53SKJ9hIM0nagK6yJeuFk9EIpMQYrAM+78bmYCu0ciAfi5mn9Qh/JcmATGJUPWRHhqpvVUZzccj/0AFzKDLjbnaD8mTuETQVACOJnI/OhBzkw==
f10jyV9CYxaMYMUY66rikK4X5xG6XUDpYJMJrKkhXSe
```
