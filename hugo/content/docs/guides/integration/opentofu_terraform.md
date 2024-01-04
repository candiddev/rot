---
categories:
- guide
description: How to use OpenTofu/Terraform with Rot
title: OpenTofu/Terraform
---

Rot can integrate with [OpenTofu](https://opentofu.org/) and [Terraform](https://www.terraform.io/) by wrapping the respective commands using [`rot run`](../../../references/cli#run).

## Prerequisites

- An existing Rot configuration setup with secret values
- OpenTofu or Terraform installed (these examples will use OpenTofu, but it should apply to Terraform, too)

## Using Rot for Provider and Terraform Settings

OpenTofu can be provided environment variables for accessing remote state storage and providers.  These environment variables can be provided by Rot, enabling a secure and seamless way for accessing your remote state storage.

The [s3 backend](https://opentofu.org/docs/language/settings/backends/s3) can use these environment variables:

- `AWS_DEFAULT_REGION`
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`

We can also use these same values for the [AWS provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs).

You'll have to generate these values on your own.  Once you have them, you can add them to your Rot configuration:

```bash
$ rot add-value AWS_DEFAULT_REGION
$ rot add-value AWS_ACCESS_KEY_ID
$ rot add-value AWS_SECRET_ACCESS_KEY
```
{{% alert title="Note" color="info" %}}
This is a completely contrived example.  Static AWS should be avoided.
{{% /alert %}}

Our Terraform backend will look like this:

```hcl
terraform {
  backend "s3" {
    bucket = "mybucket"
    key    = "my/key.tfstate"
  }
}
```

We can now initialize our remote state by wrapping OpenTofu with Rot:

```bash
$ rot run opentofu init
$ rot run opentofu apply
```

`rot run` will decrypt and pass the values we set above as environment variables automatically.

## Using Rot for OpenTofu Variables

We can define variables within OpenTofu that will have data set at runtime.  They look like this:

```hcl
variable "secret" {
  type = string
}
```

OpenTofu will prompt us to enter this value when running `opentofu apply`.  We can also set this variable using environemnt variables, like `TF_VAR_secret=something`, so Rot can set these too.

{{% alert title="Warning" color="warning" %}}
Any secret used this way will be added in plaintext to the underlying state file.  See https://github.com/hashicorp/terraform/issues/516 for more information.
{{% /alert %}}

Lets add this variable to Rot:

```bash
$ rot add-value TF_VAR_secret
```

And now we can run our apply:

```bash
$ rot run opentofu apply
```
