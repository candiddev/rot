package main

import (
	"context"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdEncrypt(ctx context.Context, args []string, c *cfg) errs.Err {
	r := ""
	delimiter := ""

	if len(args) >= 2 {
		r = args[1]
	}

	if len(args) == 3 {
		delimiter = args[2]
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
}
