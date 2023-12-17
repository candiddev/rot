// rot is a CLI tool for managing secrets.
package main

import (
	"os"

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
			"add-value": {
				ArgumentsRequired: []string{
					"name",
				},
				ArgumentsOptional: []string{
					`comment, default: ""`,
					`delimiter, default: \n`,
				},
				Run:   cmdAddValue,
				Usage: "Add a value to a configuration",
			},
			"decrypt": {
				ArgumentsRequired: []string{
					"value",
				},
				Run:   cmdDecrypt,
				Usage: "Decrypt a value and print it to stdout",
			},
			"encrypt": {
				ArgumentsRequired: []string{
					"recipient key",
				},
				ArgumentsOptional: []string{
					`delimiter, default: \n`,
				},
				Run:   cmdEncrypt,
				Usage: "Encrypt a value and print it to stdout without adding it to the config",
			},
			"generate-key": cryptolib.GenerateKeys[*cfg](),
			"generate-value": {
				ArgumentsRequired: []string{
					"name",
				},
				ArgumentsOptional: []string{
					"length, default=20",
					"comment",
				},
				Run:   cmdGenerateValue,
				Usage: "Generate a random string and encrypt it",
			},
			"init": {
				ArgumentsRequired: []string{
					"initial key name",
				},
				ArgumentsOptional: []string{
					"initial public key, default: generate a PBKDF symmetric key",
				},
				Run:   cmdInit,
				Usage: "Initialize a new Rot configuration",
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
				Usage: "Run a command and inject values as environment variables",
			},
			"show-algorithms": {
				Run:   cmdAlgorithms,
				Usage: "Show algorithms Rot understands",
			},
			"show-value": {
				ArgumentsRequired: []string{
					"name",
				},
				Run:   cmdShowValue,
				Usage: "Show a decrypted value",
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
