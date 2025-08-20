---
categories:
- guide
description: How to manage values using Rot
title: Manage Values
weight: 30
---

In this guide, we'll go over adding and removing values in a Rot Keyring.

## Adding Values

After Rot has been initialized, keys can be added using {{% cli val-add %}}.  This command requires a name for a new key, and can optionally be provided with metadata and delimiter for the value.  This command performs these actions:

- Generate a [Value Key]({{< ref "/docs/references/cryptography#value-key" >}})
- Encrypt the provided Value using the Value Key
- Encrypt the Value Key using the [Keyring Public Key]({{< ref "/docs/references/cryptography#keyring-public-key" >}})
- Populate the {{% config keyrings_values value %}} with the new value, and save the configuration to {{% config cli_configPath %}}.

### Value Paths

All Values live underneath a Keyring in Rot.  Values can have subpaths within the Keyring--this is useful for commands like {{% cli run %}}, which will convert all Values within a path into environment variables.

Value names must be valid Environment Variable names.  That means they must start with a letter or `_`, and can only have letters, numbers, and underscores in their names.

A valid path looks like this: `keyring1/a/path/VALUE`.

### Value Filtering

Rot can display and filter Values using `rot val-lst`.  This command supports using regular expressions, so you can easily filter and display Values using regex:

{{< highlight bash >}}
$ rot val-lst ^[prd|stg]/postgresql`
{{< /highlight >}}

This command would display all Values that start with `prd/postgresql` or `stg/postgresql`.

### Meta

Rot Values can optionally be provided with metadata--key/value pairs that can store useful information about the Value, such as where or how the secret was generated, or who owns the secret.

Metadata can be provided using {{% cli val-add %}} and {{% cli val-set-meta %}}.

Certain Rot commands may add data to a Value's meta:

- {{% cli val-add-prv %}} will store the public key in the `publicKey` meta field.

## Versions

Values can be modified after being added to Rot, either by adding a new version of the Value using {{% cli val-add %}} or modifying metadata using {{% cli val-set-meta %}}.  Rot will create a new timestamped version of a Value when a new version is created using {{% cli val-add %}}.

You can view all values and versions using {{% cli val-lst %}}.  Additionally, you can display a specific version in {{% cli val-dis %}} via `@`:

{{< highlight bash >}}
$ rot val-lst mykeyring
{
  "mykeyring/test": [
    "2024-04-26T12:32:52Z",
    "2024-04-26T10:32:44Z"
  ]
}
# Display the latest version
$ rot val-dis mykeyring -v mykeyring/test
123
# Display a specific version
$ rot val-dis mykeyring -v mykeyring/test@2024-04-26T10:32:44Z
456
# Use longest match
$ rot val-dis mykeyring -v 'mykeyring/test@2024-04-26T10'
456
{{< /highlight >}}

## Editing Values in an Editor

Values can be edited as decrypted Jsonnet text using {{% cli edit %}}.  This command will decrypt an entire Keyring, nicely format it into a nested Jsonnet string, and open it in your `$EDITOR` (or whatever editor you specify).

Any changes made to the Keyring will be committed--Values removed will be removed, Values changed will have a new version added, and Values added will be encrypted and added to the Keyring.

## Generating Values

Rot can generate random, cryptographically secure strings for you, instead of having you provide a value, via {{% cli val-add %}}.

## Removing Values

Values can be removed by running {{% cli val-del %}}, or editing the configuration and removing the value.
