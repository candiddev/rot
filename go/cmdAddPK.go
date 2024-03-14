package main

import (
	"context"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdAddPK() cli.Command[*cfg] {
	return cli.Command[*cfg]{
		ArgumentsRequired: []string{
			"name",
		},
		Usage: "Generate and add a cryptographic private key to the configuration values.",
		Run: func(ctx context.Context, args []string, _ cli.Flags, c *cfg) errs.Err {
			if c.PublicKey.IsNil() {
				return logger.Error(ctx, errNotInitialized)
			}

			prv, pub, err := cryptolib.NewKeysAsymmetric(c.Algorithms.Asymmetric)
			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			prv.ID = args[1]
			pub.ID = args[1]

			if err := c.encryptvalue(ctx, []byte(prv.String()), args[1], pub.String()); err != nil {
				return logger.Error(ctx, err)
			}

			return logger.Error(ctx, c.save(ctx))
		},
	}
}
