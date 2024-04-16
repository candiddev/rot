package main

import (
	"context"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdAddPK() cli.Command[*config.Config] {
	return cli.Command[*config.Config]{
		ArgumentsRequired: []string{
			"name",
		},
		ArgumentsOptional: []string{
			"key ID (default: random string)",
		},
		Usage: "Generate and add a cryptographic private key to Keyring Values.",
		Run: func(ctx context.Context, args []string, _ cli.Flags, c *config.Config) error {
			prv, pub, err := cryptolib.NewKeysAsymmetric(c.GetAlgorithms(ctx).Asymmetric)
			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			if len(args) == 3 {
				prv.ID = args[2]
				pub.ID = args[2]
			}

			return logger.Error(ctx, c.SetValue(ctx, c.GetKeyringName(ctx), []byte(prv.String()), args[1], pub.String()))
		},
	}
}
