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

func cmdSSH() cli.Command[*config.Config] {
	return cli.Command[*config.Config]{
		ArgumentsRequired: []string{
			"key, or - for stdin",
		},
		Flags: cli.Flags{
			"i": {
				Placeholder: "id",
				Usage:       "Add id to the key imported from a SSH key",
			},
		},
		Usage: "Convert Rot keys to SSH, or SSH keys to Rot, and print it to stdout.",
		Run: func(ctx context.Context, args []string, f cli.Flags, _ *config.Config) error {
			s := []byte(args[1])
			if args[1] == "-" {
				s = logger.ReadStdin()
			}

			var out []byte

			if strings.HasPrefix(string(s), "--") || strings.HasPrefix(string(s), "ssh-") {
				s, err := cryptolib.SSHToKey[cryptolib.KeyProvider](s)
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}

				s.ID, _ = f.Value("i")

				out = []byte(s.String() + "\n")
			} else {
				k, err := cryptolib.ParseKey[cryptolib.KeyProvider](string(s))
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}

				out, err = cryptolib.KeyToSSH(k)
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}
			}

			logger.Raw(string(out))

			return nil
		},
	}
}
