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
$ rot val-add mykeyring/MYSQL_ROOT
$ git diff rot.jsonnet
diff --git a/rot.jsonnet b/rot.jsonnet
--- a/rot.jsonnet
+++ b/rot.jsonnet
@@ -16,9 +16,9 @@
     values: {
       MYSQL_ROOT: {
+          '2024-04-26T15:04:11Z': {
+            key: 'ecdhx25519hkdfsha256:MCowBQYDK2VwAyEAVcW4A5ICz/f3d7bdOQTsX4SkXRHHwNTUcoDYccOCGhc=@xchacha20poly1305:9cgkADJcgH8RFVttKZBE+sPWqejGUtVSzic86kFCfEjjvAFZDqiQPFNr3vH2g77D97mjdoRa+DW+TYwyglfy6b2idxZUp+sujiEZZIDDEf1H3h9XwohS+WKd5Np899Idj4iHoE0pz+Y=:s1cAADoakP',
+            meta: {},
+            value: 'xchacha20poly1305:4XzCQq9AHJGgSYzHWUFG9sY3BbM5N7iqNBuZvBKenHZpnomCGgM1fjVilg==:tc34esW7B7',
+          },
       },
     },

# Add a new key
$ rot key-pub-add bob ed25519public:MCowBQYDK2VwAyEAAYkJzjQGb+4I7bfcaq6TnkI6nWJXolUdYSQDKSZIDZU=:bob
$ git diff rot.jsonnet
diff --git a/rot.jsonnet b/rot.jsonnet
index 9a30848..98c17ec 100644
--- a/rot.jsonnet
+++ b/rot.jsonnet
@@ -11,6 +11,12 @@
       publicKey: 'ed25519public:MCowBQYDK2VwAyEAOorjruB6HEyBUkxlz07lwhNJo4LXoSTrTbrM+0w+2i0=:alice',
     },
+    bob: {
+      modified: '2023-12-09T14:07:43.394974404-06:00',
+      publicKey: 'ed25519public:MCowBQYDK2VwAyEAAYkJzjQGb+4I7bfcaq6TnkI6nWJXolUdYSQDKSZIDZU=:bob',
+    },
   },
{{< /highlight >}}
