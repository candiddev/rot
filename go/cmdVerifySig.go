package main

import (
	"context"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdVerifySig() cli.Command[*config.Config] {
	return cli.Command[*config.Config]{
		ArgumentsRequired: []string{
			"public key value, encrypted value name, or path",
			"message or - for stdin",
			"signature",
		},
		Usage: "Verify a signature for a message using a public key.  Signature must be in the form <hash>:<signature>:<optional key id>.",
		Run: func(ctx context.Context, args []string, _ cli.Flags, c *config.Config) error {
			pk, err := c.GetPublicKey(ctx, c.GetKeyringName(ctx), args[1])
			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			m := []byte(args[2])
			if string(m) == "-" {
				m = logger.ReadStdin()
			}

			s, err := cryptolib.ParseSignature(args[3])
			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			if err := s.Verify(m, cryptolib.Keys[cryptolib.KeyProviderPublic]{
				pk,
			}); err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			return nil
		},
	}
}
