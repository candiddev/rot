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

func cmdPEM() cli.Command[*config.Config] {
	return cli.Command[*config.Config]{
		ArgumentsRequired: []string{
			"key, or - for stdin",
		},
		Flags: cli.Flags{
			"i": {
				Placeholder: "id",
				Usage:       "Add id to the key imported from a PEM",
			},
		},
		Usage: "Convert Rot keys to PEM, or a PEM keys to Rot, and print it to stdout.",
		Run: func(ctx context.Context, args []string, f cli.Flags, _ *config.Config) error {
			s := []byte(args[1])
			if args[1] == "-" {
				s = logger.ReadStdin()
			}

			var out string

			if strings.HasPrefix(string(s), "--") {
				s, err := cryptolib.PEMToKey[cryptolib.KeyProvider](s)
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}

				s.ID, _ = f.Value("i")

				out = s.String() + "\n"
			} else {
				k, err := cryptolib.ParseKey[cryptolib.KeyProvider](string(s))
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}

				out = string(cryptolib.KeyToPEM(k))
			}

			logger.Raw(out)

			return nil
		},
	}
}
