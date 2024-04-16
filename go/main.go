// rot is a CLI tool for managing secrets.
package main

import (
	"os"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/cli"
)

func m() error {
	c := config.Default()

	return (&cli.App[*config.Config]{
		Commands: map[string]cli.Command[*config.Config]{
			"add-keyprv":     cmdAddKey(true),
			"add-keypub":     cmdAddKey(false),
			"add-keyring":    cmdAddKeyring(),
			"add-pk":         cmdAddPK(),
			"add-value":      cmdAddValue(),
			"base64":         cmdBase64(),
			"copy-value":     cmdCopyValue(),
			"decrypt":        cmdDecrypt(),
			"encrypt":        cmdEncrypt(),
			"gen-crt":        cmdGenCrt(),
			"gen-jwt":        cmdGenJWT(),
			"gen-key":        cli.GenKeys[*config.Config](),
			"gen-sig":        cmdGenSig(),
			"gen-ssh":        cmdGenSSH(),
			"init":           cmdInit(),
			"pem":            cmdPEM(),
			"rekey":          cmdRekey(),
			"remove-keyprv":  cmdRemove(cmdRemoveKeyPrv),
			"remove-keypub":  cmdRemove(cmdRemoveKeyPub),
			"remove-keyring": cmdRemove(cmdRemoveKeyring),
			"remove-value":   cmdRemove(cmdRemoveValue),
			"run":            cmdRun(),
			"show-alg":       cmdShowAlg(),
			"show-crt":       cmdShowCrt(),
			"show-jwt":       cmdShowJWT(),
			"show-keys":      cmdShowKeysValues(cmdShowVluesTypeKeys),
			"show-keyings":   cmdShowKeysValues(cmdShowVluesTypeKeyrings),
			"show-pk":        cmdShowPK(),
			"show-value":     cmdShowValue(),
			"show-values":    cmdShowKeysValues(cmdShowVluesTypeValues),
			"ssh":            cmdSSH(),
			"verify-sig":     cmdVerifySig(),
		},
		Config:      c,
		Description: "Rot encrypts and decrypts secrets",
		HideConfigFields: []string{
			"key",
			"keys",
			"keyrings.*.privateKey",
		},
		Name:        "Rot",
		PricingLink: "https://rotx.dev/pricing",
	}).Run()
}

func main() {
	if err := m(); err != nil {
		os.Exit(1)
	}
}
