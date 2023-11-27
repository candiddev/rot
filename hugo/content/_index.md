---
title: Rot - Future proof secrets management
---

{{< blocks/section color="white" height=full >}}
<h1>Future Proof Secrets Management.</h1>
<h2>No footguns.  Just good, boring cryptography.</h2>

{{< highlight bash >}}
# Initialize rot.jsonnet and a key for Alice
$ ./rot init alice
New Password (empty string skips PBKDF):
Confirm Password (empty string skips PBKDF):

# View alice's new key
$ cat .rot-keys
ed25519private:MC4CAQAwBQYDK2VwBCIEINtf2nnktooZJTHfFU5SV0Ten6cmr9Qi/JRy3zAaoHVP:alice

# Add Bob's public key
$ ./rot add-key bob ed25519public:MCowBQYDK2VwAyEAAYkJzjQGb+4I7bfcaq6TnkI6nWJXolUdYSQDKSZIDZU=:bob

# Add a secret value
$ ./rot add-value MYSQL_ROOT 'MySQL root password'
Value:

# Retrieve the secret value
$ ./rot show-value MYSQL_ROOT
{
  "comment": "MySQL root password",
  "modified": "2023-12-09T08:45:15.844039874-06:00",
  "value": "hunter123"
}

# See the encrypted value filtered using JSON
$ ./rot show-config | ./rot jq '.values.MYSQL_ROOT'
{
  "comment": "MySQL root password",
  "key": "ecdhx25519:MCowBQYDK2VwAyEArXBuejuUQ...",
  "modified": "2023-12-09T08:45:15.844039874-06:00",
  "value": "xchacha20poly1305:WwJycFa/DEdJVmupdEjXURrRJmdwBcy0wMd5fEWfkZEO3E37bfxDIMACBh+jVrJ6sQ==:2y5yMhNEsO"
}

# Rekey the entire file using EC and AES instead of Ed25519 and ChaCha20
$ ./rot -x algorithms_asymmetric=ecdhp256 -x algorithms_symmetric=aes128gcm rekey

# See the encrypted value using new encryption
$ ./rot show-config | ./rot jq '.values.MYSQL_ROOT'
{
  "comment": "MySQL root password",
  "key": "ecdhp256:EwLIPFZVXT93JcueT66r...",
  "modified": "2023-12-09T08:47:59.964373017-06:00",
  "value": "aes128gcm:3GzMm6FniYrik+npa/noAZ+j18/qZzRDX7lY6W3UnUN2fKZxqg==:Nr36uOzU89"
}
{{< /highlight >}}
{{< /blocks/section >}}
