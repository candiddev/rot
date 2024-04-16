package main

import (
	"context"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdEncrypt() cli.Command[*config.Config] {
	return cli.Command[*config.Config]{
		ArgumentsOptional: []string{
			"recipient key, optional",
		},
		Flags: cli.Flags{
			"d": {
				Default:     []string{`\n`},
				Placeholder: "delimiter",
				Usage:       "Value delimiter",
			},
		},
		Usage: "Encrypt a value and print it to stdout.  Can specify a recipient key to use asymmetric encryption.",
		Run: func(ctx context.Context, args []string, f cli.Flags, c *config.Config) error {
			r := ""
			delimiter := ""

			if len(args) == 2 {
				r = args[1]
			}

			if v, ok := f.Value("d"); ok {
				delimiter = v
			}

			v, err := logger.Prompt("Value:", delimiter, true)
			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			var ev cryptolib.EncryptedValue

			if r == "" {
				ev, err = cryptolib.KDFSet(cryptolib.Argon2ID, "decrypt", v, c.GetAlgorithms(ctx).Symmetric)
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}
			} else {
				key, err := cryptolib.ParseKey[cryptolib.KeyProviderPublic](r)
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}

				ev, err = key.Key.EncryptAsymmetric(v, key.ID, c.GetAlgorithms(ctx).Symmetric)
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}
			}

			logger.Raw(ev.String() + "\n")

			return logger.Error(ctx, nil)
		},
	}
}
