package main

import (
	"context"
	"strings"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdShowPublicKey() cli.Command[*cfg] {
	return cli.Command[*cfg]{
		ArgumentsRequired: []string{
			"name, private key, or - for stdin",
		},
		Usage: "Show the public key of a private key.",
		Run: func(ctx context.Context, args []string, _ cli.Flags, c *cfg) errs.Err {
			var err error

			var key cryptolib.Key[cryptolib.KeyProviderPrivate]

			var s string

			switch {
			// Stdin
			case args[1] == "-":
				s = string(cli.ReadStdin())

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
				c.decryptKeysEncrypted(ctx)

				n := args[1]

				for i := range c.keys {
					if c.keys[i].ID == n {
						key = c.keys[i]

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
