package main

import (
	"context"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdEncrypt(ctx context.Context, args []string, c *cfg) errs.Err {
	r := args[1]
	delimiter := ""

	if len(args) == 3 {
		delimiter = args[2]
	}

	key, err := cryptolib.ParseKey[cryptolib.KeyProviderEncryptAsymmetric](r)
	if err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
	}

	v, err := cli.Prompt("Value:", delimiter, true)
	if err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
	}

	ev, err := key.Key.EncryptAsymmetric(v[0], key.ID, c.Algorithms.Symmetric)
	if err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
	}

	logger.Raw(ev.String() + "\n")

	return logger.Error(ctx, nil)
}
