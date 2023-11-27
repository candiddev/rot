---
categories:
- feature
description: Rot makes sharing secrets easy and secure
title: Secret Sharing
type: docs
---

Users can easily share secrets using Rot either using a shared secrets file, or ad-hoc using a traditional GPG-style workflow:

{{< highlight bash >}}
# Using a shared secrets file
$ ./rot add-value MYSQL_ROOT 'MySQL root password'
Value:
$ ./rot show-value MYSQL_ROOT
{
  "comment": "MySQL root password",
  "modified": "2023-12-09T08:45:15.844039874-06:00",
  "value": "hunter123"
}

# Using ad-hoc keys as Bob
$ ./rot encrypt ed25519public:MCowBQYDK2VwAyEAAYkJzjQGb+4I7bfcaq6TnkI6nWJXolUdYSQDKSZIDZU=:alice
Value:
ecdhx25519:MCowBQYDK2VwAyEAInJFvsjoY1XdTPij+mRL72NT8cUjIRA11QCFRVaoHw4=@aes128gcm:9tj1/wHN4rwC0GJG4z2/MYYPdxQhyjfyaZPwZX+tA7gaTmPOGOLLTA==:alice

# Decrypting the secret as Alice
$ ./rot decrypt ecdhx25519:MCowBQYDK2VwAyEAInJFvsjoY1XdTPij+mRL72NT8cUjIRA11QCFRVaoHw4=@aes128gcm:9tj1/wHN4rwC0GJG4z2/MYYPdxQhyjfyaZPwZX+tA7gaTmPOGOLLTA==:alice
Hello World!
{{< /highlight >}}
