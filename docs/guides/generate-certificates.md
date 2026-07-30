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

You'll need to generate a private key for every certificate, including the CA.  The easiest way to do this is using {{% cli value-add-private %}} (encrypting the keys into Rot) or {{% cli key-new %}} (printing the keys to stdout).

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
$ rot value-add-private rot/ca
$ rot x509-new -c -n 'Rot CA' rot/ca
{{< /highlight >}}

### Intermediate Certificate Authority

{{< highlight bash >}}
$ rot value-add-private rot/ca
$ rot value-add-private rot/intermediate_ca
$ rot x509-new -c -n 'Rot CA' rot/ca > ca.pem
$ rot x509-new -c -n 'Rot Intermediate CA' rot/ca rot/intermediate_ca ca.pem
{{< /highlight >}}

### Host Certificate

{{< highlight bash >}}
$ rot value-add-private rot/ca
$ rot value-add-private rot/example_com
$ rot x509-new -c -n 'Rot CA' rot/ca > ca.pem
$ rot x509-new -d www.example.com -n www.example.com rot/ca rot/example_com ca.pem
{{< /highlight >}}

## View Certificates

You can view the contents of an existing X.509 Certificate as JSON using {{% cli x509-get %}}, optionally providing a list of CA certificates to verify it against.

## Using Authenticated Certificate Management Environment (ACME)

Rot can act as an ACME client, allowing you to generate certificates from an ACME server like Let's Encrypt.

Rot currently understands these ACME challenge types:

- `DNS-01`: ACME server requests a token periodically from a TXT record associated to your domain name.
- `HTTP-01`: ACME server requests a token periodically from a URL on your HTTP server.
- `DNS-PERSIST-01`: ACME server requests an authorization for your ACME account periodically from a TXT record associated to your domain name.

You can view the various ACME commands available in Rot in the [CLI > Commands]({{% ref "/docs/references/cli#commands" %}}) documentation.

A workflow for request ACME certificates might look like this:

{{< highlight bash >}}
# Get a new Private Key to associate to your ACME account, must be ECP256, ECP384, ECP521, or RSA2048
$ rot key-new -a ecp256
New Password (empty string skips PBKDF):
Confirm Password (empty string skips PBKDF): 
[
  {
    "privateKey": "ecp256private:MIGH...",
    "publicKey": "ecp256public:MFkw..."
  }
]

# Create a new ACME account on Lets Encrypt staging
$ rot acme-account-new -s -p ecp256private:MIGH...
{
  "URI": "https://acme-staging-v02.api.letsencrypt.org/acme/acct/123",
  "Contact": null,
  "Status": "valid",
  "OrdersURL": "",
  "AgreedTerms": "",
  "CurrentTerms": "",
  "Authz": "",
  "Authorizations": "",
  "Certificates": "",
  "ExternalAccountBinding": null
}

# Create a new ACME order for a domain name
$ rot acme-order-new -s -p ecp256private:MIGH... acme.example.com
{
  "authzIDs": {
    "dns:acme.example.com": {
      "authzURI": "https://acme-v02.api.letsencrypt.org/acme/authz/123/456",
      "dns01TXTRecordName": "_acme-challenge.acme.example.com",
      "dns01TXTRecordValue": "secret-value",
      "dns01Status": "pending",
      "dns01URL": "https://acme-v02.api.letsencrypt.org/acme/chall/123/456/789",
      "expires": "2026-07-27T19:22:46Z",
      "http01ResponsePath": "/.well-known/acme-challenge/secret-token",
      "http01ResponseValue": "secret-value",
      "http01Status": "pending",
      "http01URL": "https://acme-v02.api.letsencrypt.org/acme/chall/123/456/111",
      "identifierType": "dns",
      "identifierValue": "acme.example.com"
    }
  },
  "certificate": "",
  "commonName": "acme.example.com",
  "dnsNames": [
    "acme.example.com"
  ],
  "expires": "2026-07-27T19:22:46Z",
  "finalizeURL": "https://acme-v02.api.letsencrypt.org/acme/finalize/123/456",
  "ipAddresses": [],
  "orderURI": "https://acme-v02.api.letsencrypt.org/acme/order/123/456",
  "status": "pending"
}

# Create the HTTP challenge response on a HTTP server, or add one of the DNS challenge response records.  Then test it:
$ rot acme-authz-test -s -p ecp256private:MIGH... https://acme-v02.api.letsencrypt.org/acme/authz/123/456
ERROR Error testing ACME authorization: error testing HTTP-01: GET http://acme.example.com/.well-known/acme-challenge/secret-token: got , want secret-value
{
  "DNS-01": true,
  "DNS-PERSIST-01": true,
  "HTTP-01": false
}

# Trigger the challenge URL for validation--if this fails, you must restart the order process
$ rot acme-challenge-validate -s -p ecp256private:MIGH... https://acme-v02.api.letsencrypt.org/acme/chall/123/456/789

# Finalize the certificate once an authorization passes
$ rot acme-order-finalize -s -p ecp256private:MIGH... https://acme-v02.api.letsencrypt.org/acme/order/123/456
{
  "crt": "-----BEGIN CERTIFICATE-----...",
  "key": "-----BEGIN PRIVATE KEY-----..."
}
{{< /highlight >}}
