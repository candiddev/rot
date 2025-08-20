---
categories:
- guide
description: How to create X.509 certificates using Rot
title: Generate Certificates
weight: 50
---

In this guide, we'll go over managing X.509 certificates using Rot.

## X.509 Introduction

An X.509 certificate is basically a signed hash of a public key and other fields.  A self-signed X.509 certificate (and a root Certificate Authority (CA)) will use the same public/private key pair for generating the certificate, while an intermediate CA or signed certificate will use the public key of the certificate to be signed, and the private key of the parent CA.

## Add Private Keys

You'll need to generate a private key for every certificate, including the CA.  The easiest way to do this is using {{% cli val-add-prv %}} (encrypting the keys into Rot) or {{% cli key-new %}} (printing the keys to stdout).

Rot will store the public key in the comment of the encrypted value, we can grab the public key from the comment when we generate certificates.

## Create a certificate

You can generate a certificate using a private key with {{% cli x509-new %}}.  This command generates an X.509 certificate using the options you provide and prints a PEM file to stdout for you to save.

It supports the following flags:

- `-c`: Set the IsCA flag to true and adds the keyUsage crlSign and certSign
- `-d <hostname>`: Add a DNS hostname entry to the certificate (can be specified multiple times)
- `-e <expiration>`: Sets the expiration in seconds for the certificate (default: 1 year)
- `-eu <extended key usage>`: Set an extended key usage for the certificate (can be specified multiple times).  Default: clientAuth, serverAuth
- `-i <address>`: Add an IP address entry to the certificate (can be specified multiple times)
- `-ku <key usage>`: Set a key usage for the certificate (can be specified multiple times).  Default: digitalSignature
- `-n <common name>`: Set the common name for the certificate.

Here are some example usages:

### Self Signed

{{< highlight bash >}}
$ rot x509-new -d localhost -i 127.0.0.1 -n localhost ed25519private:MC4CAQAwBQYDK2VwBCIEIAw1E0///GuHuAsxK/2gAGRRwZkrJD/mxk0HUS1VSN1a:1CjPxcEvDy
{{< /highlight >}}

### Certificate Authority

{{< highlight bash >}}
$ rot val-add-prv rot/ca
$ rot x509-new -c -n 'Rot CA' rot/ca
{{< /highlight >}}

### Intermediate Certificate Authority

{{< highlight bash >}}
$ rot val-add-prv rot/ca
$ rot val-add-prv rot/intermediate_ca
$ rot x509-new -c -n 'Rot CA' rot/ca > ca.pem
$ rot x509-new -c -n 'Rot Intermediate CA' rot/ca rot/intermediate_ca ca.pem
{{< /highlight >}}

### Host Certificate

{{< highlight bash >}}
$ rot val-add-prv rot/ca
$ rot val-add-prv rot/example_com
$ rot x509-new -c -n 'Rot CA' rot/ca > ca.pem
$ rot x509-new -d www.example.com -n www.example.com rot/ca rot/example_com ca.pem
{{< /highlight >}}

## View Certificates

You can view the contents of an existing X.509 Certificate as JSON using {{% cli x509-dis%}}, optionally providing a list of CA certificates to verify it against.
