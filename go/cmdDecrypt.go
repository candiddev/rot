package main

import (
	"context"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdDecrypt() cli.Command[*cfg] {
	return cli.Command[*cfg]{
		ArgumentsRequired: []string{
			"value or path",
		},
		Usage: "Decrypt a value or unwrap a KDF value and print it to stdout.",
		Run: func(ctx context.Context, args []string, _ cli.Flags, c *cfg) errs.Err {
			c.decryptKeysEncrypted(ctx)

			value := args[1]

			if value == "-" {
				value = string(cli.ReadStdin())
			}

			ev, err := cryptolib.ParseEncryptedValue(value)
			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			v, err := ev.Decrypt(c.keys.Keys())
			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			logger.Raw(string(v))

			return logger.Error(ctx, nil)
		},
	}
}
