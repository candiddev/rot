---
categories:
- guide
description: How to create SSH keys and certificates using Rot
title: Generate SSH
weight: 70
---

In this guide, we'll go over managing a SSH Certificate Authority (CA) using Rot.

## SSH Certificate Introduction

OpenSSH can use SSH certificate authorities (CA) to authorize user access to servers and devices without having to add individual public keys to the servers.  The user presents a valid certificate to the server which is signed by the CA trusted on the machine and is granted access based on the constraints within the certificate.

## Add Private Keys

You'll need to generate a private key to create a SSH keypair and an SSH CA to sign the certificates.  The easiest way to do this is using {{% cli val-add-prv %}} (encrypting the keys into Rot) or {{% cli key-new %}} (printing the keys to stdout).

Rot will store the public key in the comment of the encrypted value, we can grab the public key from the comment when we verify the JWT.

## Generate SSH private keys

We can use Rot to generate a SSH keypair, similar to `ssh-keygen`:

{{< highlight bash >}}
$ rot key-new | tee >(rot jq -r .privateKey | rot ssh - > id_ed25519 && chmod 0400 id_ed25519) | rot jq -r .publicKey | rot ssh - > id_ed25519.pub
{{< /highlight >}}

This will generate two files, id_ed25519 containing the SSH private key, and id_ed25519.pub containing the SSH public key.

## Generate SSH certificate private and public keys

Lets use Rot to generate another keypair, this time for use as the SSH CA:

{{< highlight bash >}}
$ rot val-add-prv rot/SSH_CA
$ rot val-dis -m publicKey rot/SSH_CA | rot ssh -
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIN+5rkhggPylubB7l9GNhrkuPX+da3iS0g5Vd9ZEhSTf
{{< /highlight >}}

The last value returned is the SSH CA.

## Add the SSH CA to a SSH Server

The addition of a SSH CA will vary depending on the operating system.  The basic steps are:

- Add the SSH CA value from the above step into a new or existing SSH CA file, such as `/etc/ssh/ssh_ca`
- Reference the SSH CA file in the server's sshd_config file with the line `TrustedUserCAKeys /etc/ssh/ssh_ca`
- Reload or restart the SSH service

## Sign the SSH public key

Now we can sign the public key we generated earlier using the SSH CA:

{{< highlight bash >}}
$ rot ssh-new -e permit-pty -p root rot/SSH_CA id_ed25519.pub > id_ed25519-cert.pub
{{< /highlight >}}

This command creates a special file, `id_ed25519-cert.pub`, that SSH will automatically present to our server for authentication.  Running this command should get us onto the server:

{{< highlight bash >}}
$ ssh -i id_ed25519 root@server
{{< /highlight >}}
