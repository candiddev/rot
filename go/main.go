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
			"add-key":      cmdAddKey(),
			"add-pk":       cmdAddPK(),
			"add-value":    cmdAddValue(),
			"base64":       cmdBase64(),
			"decrypt":      cmdDecrypt(),
			"encrypt":      cmdEncrypt(),
			"gen-crt":      cmdGenCrt(),
			"gen-jwt":      cmdGenJWT(),
			"gen-key":      cryptolib.GenerateKeys[*cfg](),
			"gen-sig":      cmdGenSig(),
			"gen-ssh":      cmdGenSSH(),
			"init":         cmdInit(),
			"pem":          cmdPEM(),
			"rekey":        cmdRekey(),
			"remove-key":   cmdRemove(true),
			"remove-value": cmdRemove(false),
			"run":          cmdRun(),
			"show-alg":     cmdShowAlg(),
			"show-crt":     cmdShowCrt(),
			"show-jwt":     cmdShowJWT(),
			"show-keys":    cmdShowKeysValues(true),
			"show-pk":      cmdShowPK(),
			"show-value":   cmdShowValue(),
			"show-values":  cmdShowKeysValues(false),
			"ssh":          cmdSSH(),
			"verify-sig":   cmdVerifySig(),
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
