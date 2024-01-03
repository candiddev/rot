package main

import (
	"context"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdDecrypt(ctx context.Context, args []string, c *cfg) errs.Err {
	c.decryptKeysEncrypted(ctx)

	value := args[1]

	if value == "-" {
		value = cli.ReadStdin()
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
}
