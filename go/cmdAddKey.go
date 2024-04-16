package main

import (
	"context"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/logger"
)

func cmdAddKey(private bool) cli.Command[*config.Config] {
	opt := []string{
		"public key, default: generate a PBKDF-protected asymmetric key",
	}
	usage := "Add a new or existing key to Decrypt Keys."

	if private {
		opt = nil
		usage = "Add an existing Decrypt Key to a Keyring"
	}

	return cli.Command[*config.Config]{
		ArgumentsRequired: []string{
			"key name",
		},
		ArgumentsOptional: opt,
		Usage:             usage,
		Run: func(ctx context.Context, args []string, _ cli.Flags, c *config.Config) error {
			name := args[1]

			if args[0] == "add-keyprv" {
				return logger.Error(ctx, c.SetDecryptKeyPrivateKey(ctx, c.GetKeyringName(ctx), name))
			}

			var err error

			var pub cryptolib.Key[cryptolib.KeyProviderPublic]

			if len(args) == 3 {
				pub, err = cryptolib.ParseKey[cryptolib.KeyProviderPublic](args[2])
			} else {
				pub, err = c.NewKeyPathKey(ctx, name)
			}

			if err != nil {
				return logger.Error(ctx, err)
			}

			return logger.Error(ctx, c.SetDecryptKey(ctx, name, pub))
		},
	}
}
