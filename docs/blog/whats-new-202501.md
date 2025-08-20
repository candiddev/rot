---
author: Candid Development
date: 2025-01-21
description: Release notes for Rot v2025.01.
tags:
  - release
title: "What's New in Rot: v2025.01"
type: blog
---

## Features

### Mailing Lists

We are moving all of our discussions from GitHub to email-based mailing lists.  See [Mailing Lists]({{< ref "/docs/references/mailing-lists" >}}) for more information.

## Enhancements

- Added a `docs` command to quickly view the Rot documentation website.
- Command line arguments now support partial command matching.  Rot will attempt to match a partial CLI command (like `cfg`) to the longest matching command or macro (in this instance, `config`).
