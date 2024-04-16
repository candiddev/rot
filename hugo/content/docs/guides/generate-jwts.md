---
categories:
- guide
description: How to create JSON Web Tokens (JWTs) using Rot
title: Generate JWTs
weight: 50
---

In this guide, we'll go over managing JWTs using Rot.

## JWT Introduction

A JSON Web Token (JWT) is a string containing three parts:

- A JSON header containing algorithm details about the JWT, base64 raw URL encoded
- A JSON payload containing key value pairs of user specified and standard claims, base64 raw URL encoded
- A cryptographic signature of the first two parts

## Add Private Keys

You'll need to generate a private key to sign the JWT.  The easiest way to do this is using [`rot add-pk`]({{< ref "/docs/references/cli#add-pk" >}}) (encrypting the keys into Rot) or [`rot gen-key`]({{< ref "/docs/references/cli#gen-key" >}}) (printing the keys to stdout).

Rot will store the public key in the comment of the encrypted value, we can grab the public key from the comment when we verify the JWT.

## Generate a JWT

You can generate a JWT using a private key with [`rot gen-jwt`]({{< ref "/docs/references/cli#gen-jwt" >}}).  This command generates a JWT using the options you provide and prints the token to stdout.

It supports the following flags:

- `-a <audience>`: The audience (aud) for the JWT.  Can be provided multiple times.
- `-e <expiration>`: The expiration (exp) in seconds from now for the JWT.  Defaults to 3600/one hour.
- `-f <key=value>`: Add a key and value to the JWT.  Will attempt to parse bools and ints unless they are quoted.  Can be provided multiple times.
- `-id <id>`: ID (jti) of the JWT, will generate a UUID if not provided
- `-is <issuer>`: Issuer (iss) of the JWT, defaults to Rot
- `-s <subject>`: Subject (sub) of the JWT

Here is an example usage:

```bash
$ rot gen-jwt -a example -e 4000 -f test=yes -id 123 -is MyIssuer -s Example ed25519private:MC4CAQAwBQYDK2VwBCIEIDp+bj8yxdPB7kSUjsqp4WNoHGnSFKeA9opbwGphFm+F:9scIk9mShr
eyJhbGciOiJFZERTQSIsImtpZCI6IjlzY0lrOW1TaHIiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJleGFtcGxlIiwiZXhwIjoxNzEwMjEyMjE4LCJpYXQiOjE3MTAyMDgyMTgsImlzcyI6Ik15SXNzdWVyIiwianRpIjoiMTIzIiwibmJmIjoxNzEwMjA4MjE4LCJzdWIiOiJFeGFtcGxlIiwidGVzdCI6InllcyJ9.aSPcgRUEmm0g4ak-OjEyyPSn0-_AxRxpFsir_f64UJ_lntR8o6Q3zulUi1IDHDtIYF4hhyutMCzMVIFkS1ufCA
```

## View JWT

You can view the contents of a JWT as JSON using [`rot show-jwt`]({{< ref "/docs/references/cli#show-jwt" >}}), optionally providing a list of public keys to verify it against:

```bash
$ rot show-jwt eyJhbGciOiJFZERTQSIsImtpZCI6IjlzY0lrOW1TaHIiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJleGFtcGxlIiwiZXhwIjoxNzEwMjEyMjE4LCJpYXQiOjE3MTAyMDgyMTgsImlzcyI6Ik15SXNzdWVyIiwianRpIjoiMTIzIiwibmJmIjoxNzEwMjA4MjE4LCJzdWIiOiJFeGFtcGxlIiwidGVzdCI6InllcyJ9.aSPcgRUEmm0g4ak-OjEyyPSn0-_AxRxpFsir_f64UJ_lntR8o6Q3zulUi1IDHDtIYF4hhyutMCzMVIFkS1ufCA ed25519public:MCowBQYDK2VwAyEASI/qzkRrx2hy3GGX1ereMpSw9+Z8KpGJ1HHjv6H+EXs=:9scIk9mShr
{
  "header": {
    "alg": "EdDSA",
    "kid": "9scIk9mShr",
    "typ": "JWT"
  },
  "payload": {
    "aud": "example",
    "exp": 1710212218,
    "iat": 1710208218,
    "iss": "MyIssuer",
    "jti": "123",
    "nbf": 1710208218,
    "sub": "Example",
    "test": "yes"
  },
  "signature": "aSPcgRUEmm0g4ak-OjEyyPSn0-_AxRxpFsir_f64UJ_lntR8o6Q3zulUi1IDHDtIYF4hhyutMCzMVIFkS1ufCA"
}

```
