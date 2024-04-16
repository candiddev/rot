package main

import (
	"context"
	"strings"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdShowPK() cli.Command[*config.Config] {
	return cli.Command[*config.Config]{
		ArgumentsRequired: []string{
			"name, private key, or - for stdin",
		},
		Usage: "Show the public key of a private key.",
		Run: func(ctx context.Context, args []string, _ cli.Flags, c *config.Config) error {
			var err error

			var key cryptolib.Key[cryptolib.KeyProviderPrivate]

			var s string

			switch {
			// Stdin
			case args[1] == "-":
				s = string(logger.ReadStdin())

				fallthrough
			// commandline
			case len(strings.Split(args[1], ":")) >= 3:
				if s == "" {
					s = args[1]
				}

				key, err = cryptolib.ParseKey[cryptolib.KeyProviderPrivate](s)
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}
			}

			if key.Key == nil {
				keys := c.GetKeys(ctx)
				n := args[1]

				for i := range keys {
					if keys[i].ID == n {
						key = keys[i]

						break
					}
				}
			}

			if key.Key != nil {
				pub, err := key.Key.Public()
				if err == nil {
					logger.Raw(cryptolib.Key[cryptolib.KeyProviderPublic]{
						ID:  key.ID,
						Key: pub,
					}.String() + "\n")

					return nil
				}
			}

			return logger.Error(ctx, errs.ErrReceiver.Wrap(cryptolib.ErrNoPrivateKey))
		},
	}
}
