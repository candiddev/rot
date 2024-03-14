package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdGenSSH() cli.Command[*cfg] {
	return cli.Command[*cfg]{
		ArgumentsRequired: []string{
			"private key value, encrypted value name, or - for stdin",
			"public key value, encrypted value name, or path",
		},
		Flags: cli.Flags{
			"c": {
				Placeholder: "option=value",
				Usage:       "Critical options to add to SSH certificate, can be specified multiple times",
			},
			"e": {
				Placeholder: "extension=value",
				Usage:       "Extensions to add to SSH certificate, can be specified multiple times",
			},
			"h": {
				Usage: "Create a host certificate (default: user certificate)",
			},
			"i": {
				Placeholder: "id",
				Usage:       "Key ID",
			},
			"p": {
				Placeholder: "principal",
				Usage:       "Valid principals to add to SSH certificate, can be specified multiple times",
			},
			"v": {
				Default:     []string{"3600"},
				Placeholder: "seconds",
				Usage:       "Seconds until certificate expires",
			},
		},
		Usage: "Generate SSH certificate and output a SSH formatted certificate. Must specify a Private Key, a Public Key, and a list of principals.  Can provide additional fields for the certificate",
		Run: func(ctx context.Context, args []string, f cli.Flags, c *cfg) errs.Err {
			pk := args[1]
			if pk == "-" {
				pk = string(cli.ReadStdin())
			}

			privateKey, errr := c.decryptValuePrivateKey(ctx, pk)
			if errr != nil {
				return logger.Error(ctx, errr)
			}

			publicKey, err := c.publicKey(args[2])
			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			opts := cryptolib.SSHSignOpts{
				CriticalOptions: map[string]string{},
				Extensions:      map[string]string{},
			}

			criticalOptions, _ := f.Values("c")

			for _, o := range criticalOptions {
				s := strings.Split(o, "=")
				v := ""
				if len(s) > 1 {
					v = strings.Join(s[1:], "=")
				}

				opts.CriticalOptions[s[0]] = v
			}

			extensions, _ := f.Values("e")

			for _, o := range extensions {
				s := strings.Split(o, "=")
				v := ""
				if len(s) > 1 {
					v = strings.Join(s[1:], "=")
				}

				opts.Extensions[s[0]] = v
			}

			_, opts.TypeHost = f.Values("h")
			opts.KeyID, _ = f.Value("i")
			opts.ValidPrincipals, _ = f.Values("p")

			valid, _ := f.Value("v")
			i, err := strconv.Atoi(valid)
			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(fmt.Errorf("error parsing v flag: %w", err)))
			}

			opts.ValidBeforeSec = i

			crt, err := cryptolib.SSHSign(privateKey, publicKey, opts)
			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			logger.Raw(string(crt))

			return nil
		},
	}
}
