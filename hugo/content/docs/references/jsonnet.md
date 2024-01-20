---
categories:
- reference
description: Reference documentation for Jsonnet libraries.
title: Jsonnet
---

Jsonnet is a configuration language for building JSON.  It provides capabilities for defining functions, pulling in external data, and other useful features to create dynamic JSON configuration files.

## Format

Jsonnet looks very similar to JSON:

{{< tabpane text=true >}}
{{< tab header="Jsonnet" >}}
{{< highlight jsonnet >}}
{
  person1: {
    name: 'Alice',
    welcome: 'Hello ' + self.name + '!',
  },
  person2: self.person1 { name: 'Bob' },
}
{{< /highlight >}}
{{< /tab >}}
{{< tab header="JSON" >}}
{{< highlight json>}}
{
  "person1": {
    "name": "Alice",
    "welcome": "Hello Alice!"
  },
  "person2": {
    "name": "Bob",
    "welcome": "Hello Bob!"
  }
}
{{< /highlight >}}
{{< /tab >}}
{{< /tabpane >}}

You can learn more about the formatting for Jsonnet on the Jsonnet website: https://jsonnet.org/.

## Key Concepts

### Standard Library

Jsonnet ships with a number of standard functions.  You can read about them here: https://jsonnet.org/ref/stdlib.html.

### Importing Jsonnet

Jsonnet can import other jsonnet files (typically with the file extension ".libsonnet").  These files can contain functions, other JSON data, and even other imports:


{{< tabpane text=true >}}
{{< tab header="main.jsonnet" >}}
{{< highlight jsonnet >}}
local func = import 'func.libsonnet';

{
  person1: {
    name: 'Alice',
    welcome: 'Hello ' + self.name + '!',
  },
  person2: func('Bob'),
}
{{< /highlight >}}
{{< /tab >}}
{{< tab header="func.libsonnet">}}
{{< highlight jsonnet>}}
function(name)
  {
    name: name,
    welcome: 'Hello ' + self.name + '!',
  }
{{< /highlight >}}
{{< /tab >}}
{{< /tabpane >}}

### Immutable

Jsonnet values are immutable.  You cannot change them, you must define a new value.  Here is how to handle common scenarios where you might want to change something:

{{< tabpane text=true >}}
{{< tab header="Objects" >}}
{{< highlight jsonnet >}}
local object1 = {
  hello: 'world'
}

object1 + {
  hello: 'person'
}
{{< /highlight >}}
{{< /tab >}}
{{< tab header="Ternary" >}}
{{< highlight jsonnet >}}
local something = if true then 'else' else 'other';
{{< /highlight >}}
{{< /tab >}}
{{< /tabpane >}}

## Native Functions

Jsonnet can be extended with custom functions via Native Functions.  These functions are non-standard and are supported only within the current Jsonnet implementation.  We provide a few of these functions:

### `getArch() string` {#getArch}

This function returns the current architecture based on GOARCH.

{{< highlight jsonnet >}}
local getArch() = std.native('getArch')();

getArch()
{{< /highlight >}}

### `getConfig() object` {#getConfig}

This function returns the current configuration for the application as a Jsonnet object.

{{< highlight jsonnet >}}
local getConfig() = std.native('getConfig')();

config()
{{< /highlight >}}

### `getEnv(key, fallback=null) string` {#getEnv}

This function returns the string value of the environment variable.  If the environment variable is not defined or does not exist, it returns an empty string or a fallback value if provided.

{{< highlight jsonnet >}}
local getEnv(key) = std.native('getEnv')(key);

getEnv('PWD')
{{< /highlight >}}

### `getFile(path, fallback=null) string`  {#getFile}

This function returns the string value of a `path` (local or http/https via GET).  The results are cached for repeat lookups within the current render cycle.  For HTTP or HTTPS paths, you can set headers for your request using a `#`, the header as a `k:v`, and deliminiated by a newline `\r\n`, e.g. `getEnv('https://example.com/api#myHeader:myValue\r\nmyOtherHeader:myOtherValue'`.

If the path is unreachable, an error will be thrown and rendering will halt.  You can optionally provide a fallback value to prevent this, this value will be returned instead on failure.

{{< highlight jsonnet >}}
local getFile(path, fallback=null) = std.native('getFile')(path, fallback);

getFile('~/.bashrc', 'fallback')
{{< /highlight >}}

### `getOS () string` {#getArch}

This function returns the current architecture based on GOOS.

{{< highlight jsonnet >}}
local getOS() = std.native('getOS')();

getOS()
{{< /highlight >}}


### `getPath() string` {#getFile}

This function returns the string value of the full directory path containing the target jsonnet file.  This value may be an empty string if the exact value cannot be determined.

{{< highlight jsonnet >}}
local getPath() = std.native('getPath')();

getPath()
{{< /highlight >}}

### `getRecord(type, name, fallback=null) []string` {#getRecord}

This function returns a list of sorted string values of a DNS record with `type` and `name`.  The results are cached for repeat lookups within the current render cycle.  The currently supported values for  `type` are `a`, `aaaa`, `cname`, and `txt`.

{{< highlight jsonnet >}}
local getRecord(type, name, fallback=null) = std.native('getRecord')(type, name, fallback);

getRecord('a', 'candid.dev', 'fallback')
{{< /highlight >}}

### `randStr(length) string` {#randStr}

This function returns a random string of length `length`.  Will panic if it cannot generate a cryptographically secure value.

{{< highlight jsonnet >}}
local randStr(length) = std.native('randStr')(length);

randStr(10)
{{< /highlight >}}

### `regexMatch(regex, string) bool` {#regexMatch}

This function returns a bool if `string` matches `regex`.  Will throw an error if `regex` is invalid/doesn't compile.

{{< highlight jsonnet >}}
local regexMatch(regex, string) = std.native('regexMatch')(regex, string);

regexMatch('^hello world$', 'hello world')
{{< /highlight >}}

### `render(string) object` {#render}

This function renders `string` using Jsonnet.

{{< highlight jsonnet >}}
std.native('render')(|||
  local regexMatch(regex, string) = std.native('regexMatch')(regex, string);

  regexMatch('^hello world$', 'hello world')
|||)
{{< /highlight >}}

## Best Practices

### Always Wrap Your Ifs

`if`s can be wrapped with parenthesis in Jsonnet or not wrapped.  By keeping `if`s always wrapped, it makes it easier to understand where they end:

{{< highlight jsonnet >}}
local noWrap = if 'this' then 'that' else 'that' + 'yep'
local withWrap = (if 'this' then 'that' else 'that') + 'yep'
{{< /highlight >}}

### Formatting

"Proper" Jsonnet format/linting recommends:

- Single quotes for strings
- Two spaces, no tabs
