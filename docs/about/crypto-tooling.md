---
categories:
- feature
description: Rot is a cryptographic swiss army knife
title: Crypto Tooling
type: docs
---

Rot can read, convert, and display various forms of cryptographic keys:

- JSON Web Tokens (JWTs)
- PEM Keys and Certificates
- SSH Keys and Certificates
- X.509 Certificates

Using Rot, you can easily manage all of your cryptography needs without additional tooling:

{{< highlight bash >}}
$ rot x509-dis ca.pem
$ rot jwt-dis eyJhbGciOiJFZERTQSIsImtpZCI6IjlzY0lrOW1TaHIiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJleGFtcGxlIiwiZXhwIjoxNzEwMjEyMjE4LCJpYXQiOjE3MTAyMDgyMTgsImlzcyI6Ik15SXNzdWVyIiwianRpIjoiMTIzIiwibmJmIjoxNzEwMjA4MjE4LCJzdWIiOiJFeGFtcGxlIiwidGVzdCI6InllcyJ9.aSPcgRUEmm0g4ak-OjEyyPSn0-_AxRxpFsir_f64UJ_lntR8o6Q3zulUi1IDHDtIYF4hhyutMCzMVIFkS1ufCA ed25519public:MCowBQYDK2VwAyEASI/qzkRrx2hy3GGX1ereMpSw9+Z8KpGJ1HHjv6H+EXs=:9scIk9mShr
$ cat ca.pem | rot pem - > ca.rot
$ cat ca.rot | rot pem - > ca-new.pem
$ cat ~/.ssh/id_rsa.pub | rot ssh-dis - > ssh.rot
{{< /highlight >}}
