package main

import (
	"context"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/logger"
)

func cmdDecrypt() cli.Command[*config.Config] {
	return cli.Command[*config.Config]{
		ArgumentsRequired: []string{
			"value or - for stdin",
		},
		Usage: "Decrypt a value or unwrap a KDF value and print it to stdout.",
		Run: func(ctx context.Context, args []string, _ cli.Flags, c *config.Config) error {
			value := args[1]

			if value == "-" {
				value = string(logger.ReadStdin())
			}

			ev, err := cryptolib.ParseEncryptedValue(value)
			if err != nil {
				return logger.Error(ctx, err)
			}

			v, err := c.GetDecrypted(ctx, ev)
			if err != nil {
				return logger.Error(ctx, err)
			}

			logger.Raw(string(v))

			return logger.Error(ctx, nil)
		},
	}
}
