package main

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdAddKey(ctx context.Context, args []string, c *cfg) errs.Err {
	e := c.decryptPrivateKey(ctx)
	if e != nil {
		if e.Is(errNotInitialized) {
			return cmdInit(ctx, args, c)
		}

		return e
	}

	var err error

	var p string

	var pub cryptolib.Key[cryptolib.KeyProviderPublic]

	n := args[1]

	if len(args) == 3 {
		pub, err = cryptolib.ParseKey[cryptolib.KeyProviderPublic](args[2])
		if err != nil {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}
	} else {
		var prv cryptolib.Key[cryptolib.KeyProviderPrivate]

		prv, pub, err = cryptolib.NewKeysAsymmetric(c.Algorithms.Asymmetric)
		if err != nil {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}

		prv.ID = n
		pub.ID = n

		v, err := cryptolib.KDFSet(cryptolib.Argon2ID, prv.ID, []byte(prv.String()), c.Algorithms.Symmetric)
		if err != nil {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}

		if v.Ciphertext == "" {
			p = prv.String()
		} else {
			p = v.String()
		}

		f, err := os.OpenFile(c.KeyPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(errors.New("error opening keys file"), err))
		}

		defer f.Close()

		if _, err := f.WriteString(p + "\n"); err != nil {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(errors.New("error writing keys file"), err))
		}
	}

	v, err := pub.Key.EncryptAsymmetric([]byte(c.privateKey.String()), pub.ID, c.Algorithms.Symmetric)
	if err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
	}

	sig, err := cryptolib.NewSignature(c.privateKey, []byte(pub.String()))
	if err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
	}

	v.KeyID = c.privateKey.ID

	c.DecryptKeys[n] = cfgDecryptKey{
		Modified:   time.Now(),
		PrivateKey: v,
		PublicKey:  pub,
		Signature:  sig,
	}

	return logger.Error(ctx, c.save(ctx))
}
