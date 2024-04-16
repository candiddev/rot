package main

import (
	"context"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
)

func cmdShowAlg() cli.Command[*config.Config] {
	return cli.Command[*config.Config]{
		Usage: "Show algorithms Rot understands.",
		Run: func(_ context.Context, _ []string, _ cli.Flags, _ *config.Config) error {
			return cli.Print(map[string]any{
				"asymmetric":     cryptolib.EncryptionAsymmetric,
				"asymmetricBest": cryptolib.BestEncryptionAsymmetric,
				"pbkdf":          cryptolib.ValidPBKDF,
				"pbkdfBest":      cryptolib.KDFArgon2ID,
				"symmetric":      cryptolib.EncryptionSymmetric,
				"symmetricBest":  cryptolib.BestEncryptionSymmetric,
			})
		},
	}
}
