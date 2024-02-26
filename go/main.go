// rot is a CLI tool for managing secrets.
package main

import (
	"os"
	"strings"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
)

func m() errs.Err {
	c := defaultCfg()

	return (cli.App[*cfg]{
		Commands: map[string]cli.Command[*cfg]{
			"add-key": {
				ArgumentsRequired: []string{
					"key name",
				},
				ArgumentsOptional: []string{
					"public key, default: generate a PBKDF-protected asymmetric key",
				},
				Run:   cmdAddKey,
				Usage: "Add a key to a configuration",
			},
			"add-private-key": {
				ArgumentsRequired: []string{
					"name",
				},
				Run:   cmdAddPrivateKey,
				Usage: "Add a private key to a configuration",
			},
			"add-value": {
				ArgumentsRequired: []string{
					"name",
				},
				ArgumentsOptional: []string{
					`comment, default: ""`,
				},
				Flags: cli.Flags{
					"d": {
						Default: []string{`\n`},
						Usage:   "Delimiter",
					},
				},
				Run:   cmdAddValue,
				Usage: "Add a value to a configuration",
			},
			"decrypt": {
				ArgumentsRequired: []string{
					"value, or - for stdin",
				},
				Run:   cmdDecrypt,
				Usage: "Decrypt a value and print it to stdout",
			},
			"encrypt": {
				ArgumentsOptional: []string{
					"recipient key, optional",
				},
				Flags: cli.Flags{
					"d": {
						Default: []string{`\n`},
						Usage:   "Delimiter",
					},
				},
				Run:   cmdEncrypt,
				Usage: "Encrypt a value and print it to stdout without adding it to the config.  Can specify a recipient key to use asymmetric encryption.",
			},
			"generate-certificate": {
				ArgumentsRequired: []string{
					"private key value, name, or - for stdin",
				},
				ArgumentsOptional: []string{
					"public key",
					"ca certificate or path",
				},
				Flags: cli.Flags{
					"c": {
						Usage: "Create a CA certificate",
					},
					"d": {
						Placeholder: "hostname",
						Usage:       "DNS hostname (can be used multiple times)",
					},
					"e": {
						Default:     []string{"31536000"},
						Placeholder: "seconds",
						Usage:       "Expiration in seconds",
					},
					"eu": {
						Default:     []string{"clientAuth", "serverAuth"},
						Placeholder: "extended key usage",
						Usage:       "Extended key usage, valid values: " + strings.Join(cryptolib.ValidX509ExtKeyUsages(), ", "),
					},
					"i": {
						Placeholder: "address",
						Usage:       "IP address (can be used multiple times)",
					},
					"ku": {
						Default:     []string{"digitalSignature"},
						Placeholder: "key usage",
						Usage:       "Key usage, valid values: " + strings.Join(cryptolib.ValidX509KeyUsages(), ", "),
					},
					"n": {
						Placeholder: "name",
						Usage:       "Common Name (CN)",
					},
				},
				Run:   cmdGenerateCertificate,
				Usage: "Generate X.509 certificates using Private Keys.  Must specify the private key of the signer (for CA signed certificates) or the private key of the certificate (for self-signed certificates).  A public key can be specified, otherwise the public key of the private key will be used.  Outputs a PEM certificate.",
			},
			"generate-key": cryptolib.GenerateKeys[*cfg](),
			"generate-value": {
				ArgumentsRequired: []string{
					"name",
				},
				ArgumentsOptional: []string{
					"comment",
				},
				Flags: cli.Flags{
					"l": {
						Default: []string{"20"},
						Usage:   "Length of value",
					},
				},
				Run:   cmdGenerateValue,
				Usage: "Generate a random string and encrypt it",
			},
			"init": {
				ArgumentsOptional: []string{
					"key name or id of an existing key",
					"initial public key, default: generate a PBKDF symmetric key",
				},
				Run:   cmdInit,
				Usage: "Initialize a new Rot configuration.  Will look for a .rot-keys file and use the first available key if none specified as the initial user key.",
			},
			"pem": {
				ArgumentsRequired: []string{
					"key, or - for stdin",
				},
				Flags: cli.Flags{
					"i": {
						Placeholder: "id",
						Usage:       "Add id to the key imported from a PEM",
					},
				},
				Run:   cmdPEM,
				Usage: "Convert a key to PEM, or a PEM file to a key",
			},
			"rekey": {
				Run:   cmdRekey,
				Usage: "Rekey all keys and values",
			},
			"remove-key": {
				ArgumentsRequired: []string{
					"name",
				},
				Run:   cmdRemove,
				Usage: "Remove a key from the configuration",
			},
			"remove-value": {
				ArgumentsRequired: []string{
					"name",
				},
				Run:   cmdRemove,
				Usage: "Remove a value from the configuration",
			},
			"run": {
				ArgumentsRequired: []string{
					"command",
				},
				Run:   cmdRun,
				Usage: "Run a command and inject values as environment variables.  Values written to stderr/stdout will be masked with ***.",
			},
			"show-algorithms": {
				Run:   cmdAlgorithms,
				Usage: "Show algorithms Rot understands",
			},
			"show-keys": {
				Run:   cmdShowKeysValues,
				Usage: "Show decrypt key names",
			},
			"show-private-key": {
				Run:   cmdShowPrivateKey,
				Usage: "Show the decrypted private key",
			},
			"show-public-key": {
				ArgumentsRequired: []string{
					"name, private key, or - for stdin",
				},
				Run:   cmdShowPublicKey,
				Usage: "Show the public key of a private key using a name, the private key contents, or - for stdin",
			},
			"show-value": {
				ArgumentsRequired: []string{
					"name",
				},
				Flags: cli.Flags{
					"c": {
						Usage: "Show the comment only",
					},
					"v": {
						Usage: "Show the value only",
					},
				},
				Run:   cmdShowValue,
				Usage: "Show a decrypted value",
			},
			"show-values": {
				Run:   cmdShowKeysValues,
				Usage: "Show available value names",
			},
		},
		Config:      c,
		Description: "Rot encrypts and decrypts secrets",
		HideConfigFields: []string{
			"key",
			"keys",
			"privateKey",
		},
		Name: "Rot",
	}).Run()
}

func main() {
	if err := m(); err != nil {
		os.Exit(1)
	}
}
