package main

import (
	"context"
	"time"

	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdRekey(ctx context.Context, _ []string, c *cfg) errs.Err {
	e := c.decryptKeys(ctx)
	if e != nil {
		return e
	}

	prv, pub, err := cryptolib.NewKeysEncryptAsymmetric(c.Algorithms.Asymmetric)
	if err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
	}

	c.PublicKey = pub

	for k, v := range c.DecryptKeys {
		p, err := c.DecryptKeys[k].PublicKey.Key.EncryptAsymmetric([]byte(prv.String()), c.DecryptKeys[k].PublicKey.ID, c.Algorithms.Symmetric)
		if err != nil {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}

		sig, err := cryptolib.NewSignature(prv.Key, c.privateKey.ID, []byte(v.PublicKey.String()))
		if err != nil {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}

		v.Modified = time.Now()
		v.PrivateKey = p
		v.Signature = sig

		c.DecryptKeys[k] = v
	}

	for k := range c.Values {
		v, err := c.decryptValue(ctx, k)
		if err != nil {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}

		if err := c.encryptvalue(ctx, v, k, c.Values[k].Comment); err != nil {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}
	}

	return logger.Error(ctx, c.save(ctx))
}
