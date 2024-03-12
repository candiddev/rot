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
			"add-key":              cmdAddKey(),
			"add-private-key":      cmdAddPrivateKey(),
			"add-value":            cmdAddValue(),
			"decrypt":              cmdDecrypt(),
			"encrypt":              cmdEncrypt(),
			"generate-certificate": cmdGenerateCertificate(),
			"generate-jwt":         cmdGenerateJWT(),
			"generate-key":         cryptolib.GenerateKeys[*cfg](),
			"generate-ssh":         cmdGenerateSSH(),
			"init":                 cmdInit(),
			"pem":                  cmdPEM(),
			"rekey":                cmdRekey(),
			"remove-key":           cmdRemove(true),
			"remove-value":         cmdRemove(false),
			"run":                  cmdRun(),
			"show-algorithms":      cmdAlgorithms(),
			"show-certificate":     cmdShowCertificate(),
			"show-jwt":             cmdShowJWT(),
			"show-keys":            cmdShowKeysValues(true),
			"show-public-key":      cmdShowPublicKey(),
			"show-value":           cmdShowValue(),
			"show-values":          cmdShowKeysValues(false),
			"ssh":                  cmdSSH(),
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
