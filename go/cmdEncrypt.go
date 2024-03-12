package main

import (
	"context"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdEncrypt() cli.Command[*cfg] {
	return cli.Command[*cfg]{
		ArgumentsOptional: []string{
			"recipient key, optional",
		},
		Flags: cli.Flags{
			"d": {
				Default: []string{`\n`},
				Usage:   "Delimiter",
			},
		},
		Usage: "Encrypt a value and print it to stdout.  Can specify a recipient key to use asymmetric encryption.",
		Run: func(ctx context.Context, args []string, f cli.Flags, c *cfg) errs.Err {
			r := ""
			delimiter := ""

			if len(args) >= 2 {
				r = args[1]
			}

			if v, ok := f.Value("d"); ok {
				delimiter = v
			}

			v, err := cli.Prompt("Value:", delimiter, true)
			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			var ev cryptolib.EncryptedValue

			if r == "" {
				ev, err = cryptolib.KDFSet(cryptolib.Argon2ID, "decrypt", v, c.Algorithms.Symmetric)
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}
			} else {
				key, err := cryptolib.ParseKey[cryptolib.KeyProviderPublic](r)
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}

				ev, err = key.Key.EncryptAsymmetric(v, key.ID, c.Algorithms.Symmetric)
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}
			}

			logger.Raw(ev.String() + "\n")

			return logger.Error(ctx, nil)
		},
	}
}
