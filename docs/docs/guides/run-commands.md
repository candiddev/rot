---
categories:
- guide
description: How to run commands using Rot
title: Run Commands
weight: 40
---

In this guide, we'll go over running commands with injected secrets via environment variables.

## Run Process

After Rot has been [initialized]({{< ref "/docs/guides/initialize-rot" >}}) and [populated with values]({{< ref "/docs/guides/manage-values" >}}), Rot can be used to wrap commands with secrets via environment variables.

The command to perform this is {{% cli run %}}.  This command performs these actions:

- Decrypt the existing [Rot Private Key]({{< ref "/docs/references/cryptography#rot-private-key" >}}) using the [User Private Key]({{< ref "/docs/references/cryptography#user-private-key" >}}).
- Decrypt each value
- Run the specified command with an environment variable set for each value within the specified path.  The environment variable name will be the values name by default, e.g. the key `secrets/key1` with a secret value `123` will be set to `key1=123`.  You can specify a custom suffix for the environment variable using the meta key, `_`--setting that key for the last example to `abc` would change the environment variable to `key1_abc=123`.
- For each value with {{% config keyrings_values_meta meta %}} values, the Values will be set as environment variables using the Value's name.  Given this Value:
```json
"TEST": {
  "meta": {
    "test1": "123",
    "test2": "456",
  }
}
```
The following variables will be set (in addition to the `TEST` variable):
```bash
TEST_test1=123
TEST_test2=456
```

You can specify the path to look for values when using the run command:

{{< highlight bash >}}
$ rot run keyring1/path1,keyring2/path2 <command>
{{< /highlight >}}

These paths accept regular expressions as well.
