package main

import (
	"context"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
)

func cmdShowAlg() cli.Command[*cfg] {
	return cli.Command[*cfg]{
		Usage: "Show algorithms Rot understands.",
		Run: func(_ context.Context, _ []string, _ cli.Flags, _ *cfg) errs.Err {
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
