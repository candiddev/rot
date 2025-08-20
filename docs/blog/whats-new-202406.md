---
author: Candid Development
date: 2024-06-03
description: Release notes for Rot v2024.06.
tags:
  - release
title: "What's New in Rot: v2024.06"
type: blog
---

## Features

### Hashes

Rot now includes hash CLI helpers to generate and verify hashes.  Visit [`rot hash-new`]({{< ref "/docs/references/cli#hash-new" >}}) and [`rot hash-ver`]({{< ref "/docs/references/cli#hash-ver" >}}) for more information.


## Enhancements

- Added {{< config cli.macros macros >}} for creating repeatable CLI shortcuts.
- Change {{% cli run %}} to allow the base environment variable name to set for values using the meta key `_`, e.g. `secrets/mykey` with the meta key `_` set to `test` will be `mykey_test` instead of `mykey` when using {{% cli run %}}.
