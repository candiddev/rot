package main

import (
	"context"

	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdShowPublicKey(ctx context.Context, args []string, c *cfg) errs.Err {
	c.decryptKeysEncrypted(ctx)

	n := args[1]

	for i := range c.keys {
		if c.keys[i].ID == n {
			key, err := c.keys[i].Key.Public()
			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			logger.Raw(cryptolib.Key[cryptolib.KeyProviderPublic]{
				ID:  c.keys[i].ID,
				Key: key,
			}.String() + "\n")

			return logger.Error(ctx, nil)
		}
	}

	return logger.Error(ctx, errs.ErrReceiver.Wrap(cryptolib.ErrNoPrivateKey))
}
