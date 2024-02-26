package main

import (
	"context"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
)

func cmdAlgorithms(_ context.Context, _ []string, _ cli.Flags, _ *cfg) errs.Err {
	return cli.Print(map[string]any{
		"asymmetric":     cryptolib.EncryptionAsymmetric,
		"asymmetricBest": cryptolib.BestEncryptionAsymmetric,
		"pbkdf":          cryptolib.ValidPBKDF,
		"pbkdfBest":      cryptolib.KDFArgon2ID,
		"symmetric":      cryptolib.EncryptionSymmetric,
		"symmetricBest":  cryptolib.BestEncryptionSymmetric,
	})
}
