package main

import (
	"context"
	"errors"
	"time"

	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
	"github.com/candiddev/shared/go/types"
)

var errNotInitialized = errs.ErrReceiver.Wrap(errors.New("rot not initialized"))

// decryptKeys will decrypt all cfg.keysEncrypted and use all keys in cfg.keys to decrypt cfg.DecryptKeys until cfg.privateKey is found.
func (c *cfg) decryptKeys(ctx context.Context) errs.Err {
	if !c.privateKey.IsNil() {
		return nil
	}

	if len(c.DecryptKeys) == 0 {
		return errNotInitialized
	}

	for i := range c.keysEncrypted {
		out, err := c.keysEncrypted[i].Decrypt(c.keys)
		if err == nil {
			var k cryptolib.Key[cryptolib.KeyProvider]

			k, err = cryptolib.ParseKey[cryptolib.KeyProvider](string(out))
			if err == nil {
				c.keys = append(c.keys, k.Key)
			}
		}

		if err != nil {
			logger.Debug(ctx, errs.ErrReceiver.Wrap(err).Error())
		}
	}

	c.keysEncrypted = nil

	var err error

	var out []byte

	for i := range c.DecryptKeys {
		out, err = c.DecryptKeys[i].PrivateKey.Decrypt(c.keys)
		if err == nil {
			k, err := cryptolib.ParseKey[cryptolib.KeyProviderDecryptAsymmetric](string(out))
			if err == nil {
				c.privateKey = k

				break
			}
		}
	}

	if err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
	}

	return nil
}

// decryptValue will lookup a cfg.Values using value, decrypt Key, and use Key to decrypt Value.
func (c *cfg) decryptValue(ctx context.Context, value string) ([]byte, errs.Err) {
	if v, ok := c.Values[value]; ok {
		kb, err := c.privateKey.Key.DecryptAsymmetric(v.Key)
		if err != nil {
			return nil, logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}

		k, err := cryptolib.ParseKey[cryptolib.KeyProviderEncryptSymmetric](string(kb))
		if err != nil {
			return nil, logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}

		out, err := k.Key.DecryptSymmetric(v.Value)
		if err != nil {
			return nil, logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}

		return out, logger.Error(ctx, nil)
	}

	return nil, logger.Error(ctx, errs.ErrReceiver.Wrap(errors.New("value not found")))
}

// encryptValue will encrypt a value and add it to the config.
func (c *cfg) encryptvalue(ctx context.Context, value []byte, name, comment string) errs.Err {
	if err := types.EnvValidate(name); err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(errors.New("invalid name"), err))
	}

	key, err := cryptolib.NewKeyEncryptSymmetric(c.Algorithms.Symmetric)
	if err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
	}

	ev, err := key.Key.EncryptSymmetric(value, key.ID)
	if err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
	}

	ek, err := c.PublicKey.Key.EncryptAsymmetric([]byte(key.String()), c.PublicKey.ID, c.Algorithms.Symmetric)
	if err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
	}

	c.Values[name] = cfgValue{
		Comment:  comment,
		Key:      ek,
		Modified: time.Now(),
		Value:    ev,
	}

	return logger.Error(ctx, nil)
}
