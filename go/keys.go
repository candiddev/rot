package main

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
	"github.com/candiddev/shared/go/types"
)

var errNotInitialized = errs.ErrReceiver.Wrap(errors.New("rot not initialized, run rot init"))
var errNotFound = errs.ErrReceiver.Wrap(errors.New("value not found"))

// decryptKeysEncrypted will decrypt all cfg.keysEncrypted.
func (c *cfg) decryptKeysEncrypted(ctx context.Context) {
	keys := c.keys.Keys()

	for i := range c.keysEncrypted {
		out, err := c.keysEncrypted[i].Decrypt(keys)
		if err == nil {
			var k cryptolib.Key[cryptolib.KeyProviderPrivate]

			k, err = cryptolib.ParseKey[cryptolib.KeyProviderPrivate](string(out))
			if err == nil {
				c.keys = append(c.keys, k)
			}
		}

		if err != nil {
			logger.Debug(ctx, errs.ErrReceiver.Wrap(err).Error())
		}
	}

	c.keysEncrypted = nil
}

// decryptPrivateKey will use all keys in cfg.keys to decrypt cfg.DecryptKeys until cfg.privateKey is found.
func (c *cfg) decryptPrivateKey(ctx context.Context) errs.Err {
	if !c.privateKey.IsNil() {
		return nil
	}

	if len(c.DecryptKeys) == 0 {
		return errNotInitialized
	}

	c.decryptKeysEncrypted(ctx)

	var err error

	var out []byte

	keys := c.keys.Keys()

	for i := range c.DecryptKeys {
		out, err = c.DecryptKeys[i].PrivateKey.Decrypt(keys)
		if err == nil {
			k, err := cryptolib.ParseKey[cryptolib.KeyProviderPrivate](string(out))
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

		k, err := cryptolib.ParseKey[cryptolib.KeyProviderSymmetric](string(kb))
		if err != nil {
			return nil, logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}

		out, err := k.Key.DecryptSymmetric(v.Value)
		if err != nil {
			return nil, logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}

		return out, logger.Error(ctx, nil)
	}

	return nil, logger.Error(ctx, errNotFound)
}

// decryptValuePrivateKey will lookup a PrivateKey using value.
func (c *cfg) decryptValuePrivateKey(ctx context.Context, privateKey string) (cryptolib.Key[cryptolib.KeyProviderPrivate], errs.Err) {
	pk := cryptolib.Key[cryptolib.KeyProviderPrivate]{}

	if len(strings.Split(privateKey, ":")) == 1 {
		if err := c.decryptPrivateKey(ctx); err != nil {
			return pk, logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}

		v, err := c.decryptValue(ctx, privateKey)
		if err != nil {
			return pk, logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}

		privateKey = string(v)
	}

	pk, err := cryptolib.ParseKey[cryptolib.KeyProviderPrivate](privateKey)
	if err != nil {
		return pk, logger.Error(ctx, errs.ErrReceiver.Wrap(errors.New("error parsing private key"), err))
	}

	return pk, logger.Error(ctx, nil)
}

// encryptValue will encrypt a value and add it to the config.
func (c *cfg) encryptvalue(ctx context.Context, value []byte, name, comment string) errs.Err {
	if err := types.EnvValidate(name); err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(errors.New("invalid name"), err))
	}

	key, err := cryptolib.NewKeySymmetric(c.Algorithms.Symmetric)
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

func (c *cfg) publicKey(value string) (cryptolib.Key[cryptolib.KeyProviderPublic], error) {
	pk := value
	if !strings.Contains(pk, ":") {
		pk = c.Values[value].Comment

		if pk == "" {
			f, err := os.ReadFile(value)
			if err == nil {
				pk = string(f)
			}
		}
	}

	return cryptolib.ParseKey[cryptolib.KeyProviderPublic](pk)
}
