package main

import (
	"context"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/jwt"
	"github.com/candiddev/shared/go/logger"
)

func cmdShowJWT() cli.Command[*cfg] {
	return cli.Command[*cfg]{
		ArgumentsRequired: []string{
			"JWT, or - for stdin",
		},
		ArgumentsOptional: []string{
			"public key value, can be specified multiple times",
		},
		Usage: "Show a JWT, optionally validating the signature with a public key.",
		Run: func(ctx context.Context, args []string, flags cli.Flags, config *cfg) errs.Err {
			j := args[1]
			if j == "-" {
				j = string(cli.ReadStdin())
			}

			keys := cryptolib.Keys[cryptolib.KeyProviderPublic]{}

			if len(args) > 2 {
				for i := range args[2:] {
					key, err := cryptolib.ParseKey[cryptolib.KeyProviderPublic](args[i+2])
					if err != nil {
						return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
					}

					keys = append(keys, key)
				}
			}

			t, _, err := jwt.Parse(j, keys)
			if err != nil {
				logger.Error(ctx, errs.ErrReceiver.Wrap(err)) //nolint:errcheck
			}

			if t != nil {
				s, errr := t.Values()
				logger.Raw(s + "\n")

				if errr != nil {
					err = errr
				}
			}

			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			return nil
		},
	}
}
