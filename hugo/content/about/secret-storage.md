---
categories:
- feature
description: Rot stores your secrets and keys securely
title: Secret Storage
type: docs
---

Rot stores secrets for offline use in a human readable, easily diff-able Jsonnet file.  Rot is designed to store secrets securely, even in untrusted locations.

{{< highlight bash >}}
# Replace an existing value
$ ./rot add-value MYSQL_ROOT
$ git diff rot.jsonnet
diff --git a/rot.jsonnet b/rot.jsonnet
--- a/rot.jsonnet
+++ b/rot.jsonnet
@@ -16,9 +16,9 @@
   values: {
     MYSQL_ROOT: {
       comment: 'MySQL root password',
-      key: 'ecdhx25519:MCowBQYDK2VwAyEAU8rmKnEy332m84BkGgUxkEeCCGv/m6vJ9wQ7ZNRNwMk=@xchacha20poly1305:R1D7uk9cuW4zgKyURr950O5DVbIXCUXKiODogEwWdrmCZVOOkBNdvhvU0zoi+yRybDeXFZRQlnpOjw7qeQTK/tYBVeEn8HjpAY3ZPywyDtIXZX+MpHF/7Ur2rI2LN5irm5G3DSKYC7A=:LbdyKg8Q3s',
-      modified: '2023-12-09T14:04:33.242792892-06:00',
-      value: 'xchacha20poly1305:O3dT88IPvNJ62q/PrXHkYuXwqyKecDC5y1ey6TfMpSUwz2DTeUOiKHJli/OZUlTLig==:j7CaSdbyuM',
+      key: 'ecdhx25519:MCowBQYDK2VwAyEArCxcLJyhRlpX0U/L58elyBR4KLW04zXDYaVC06fM8S4=@xchacha20poly1305:Lf7vIW/ZI4l+j6fH641E2F0yZoMBUioU4SssoCs0x3TISxh4dHX9gYhySdb4LtG8YnkHxDkckdzHAiHDTy8GZ2OCAAaRg6B2yi53QvBcSH6vJRY3eRY8PcMwHAmgxjLPabks9pLoObY=:LbdyKg8Q3s',
+      modified: '2023-12-09T14:04:43.308490659-06:00',
+      value: 'xchacha20poly1305:sX7UBOtlEYqycwPhBAlim4m4t9+m2TDbBJWpQaFmRS9wIMQ5SWt40kEEZA==:xSSl2068mO',
     },
   },
 }

# Add a new key
$ ./rot add-key bob ed25519public:MCowBQYDK2VwAyEAAYkJzjQGb+4I7bfcaq6TnkI6nWJXolUdYSQDKSZIDZU=:bob
$ git diff rot.jsonnet
diff --git a/rot.jsonnet b/rot.jsonnet
index 9a30848..98c17ec 100644
--- a/rot.jsonnet
+++ b/rot.jsonnet
@@ -11,6 +11,12 @@
       publicKey: 'ed25519public:MCowBQYDK2VwAyEAOorjruB6HEyBUkxlz07lwhNJo4LXoSTrTbrM+0w+2i0=:alice',
       signature: 'ed25519:TjnIMGwoV0fnS/PH12SEsg+jOqZtvp/QmBD7Iu4w7iMNwic7VgyFCVcYWbgKVtDrSlkyzGdT0aypYClmmepiAA==:LbdyKg8Q3s',
     },
+    bob: {
+      modified: '2023-12-09T14:07:43.394974404-06:00',
+      privateKey: 'ecdhx25519:MCowBQYDK2VwAyEAC4X7sjXGFXnt4B9r+qfOxaU0GJnkDsiOjCoW0qUsNrg=@xchacha20poly1305:d9eCFc+cwIKsPgYK3n8oF6Sl8/IMq8l/pCD0TAx5eoV8M/LQFoCbl/w/d1zFNHEUGSxAfIV2XWM/bKat1lb9FjJJLbFNlVplZ1qjyuo+YSv6JwhVBrktPCrz2CMNbK+Vuv6d6EDxbxq2khvFO/moEacLugNU9Z462jI7aiwQmb1ZJA==:LbdyKg8Q3s',
+      publicKey: 'ed25519public:MCowBQYDK2VwAyEAAYkJzjQGb+4I7bfcaq6TnkI6nWJXolUdYSQDKSZIDZU=:bob',
+      signature: 'ed25519:G1wel3K7isWH2xyhPY5oTCMuBc5Z/oIBFk9vx16lU+6UYJ6YHlflu1KodCTBemblMiAQ2E1P4kIwK9scm1EhCg==:LbdyKg8Q3s',
+    },
   },
   publicKey: 'ed25519public:MCowBQYDK2VwAyEAbhpjKAzmchYVA3Wm/lIiTdThYZr0rLkGXh0nG3oJ6Jw=:LbdyKg8Q3s',
   values: {
{{< /highlight >}}
