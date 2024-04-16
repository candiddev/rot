package config

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/candiddev/shared/go/config"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/jsonnet"
	"github.com/candiddev/shared/go/logger"
)

var (
	ErrDecryptKeyNotFound = errors.New("decryptKey not found")
	ErrKeyringNotFound    = errors.New("keyring not found")
	ErrNotInitialized     = errors.New("rot not initialized, run rot init")
	ErrValueNotFound      = errors.New("value not found")
)

// clear will remove decrypted private keys from keyrings.  Should only be used by tests.
func (c *Config) clear(ctx context.Context) {
	ctx = c.lock(ctx)
	defer c.unlock(ctx)

	for k, v := range c.cfg.Keyrings {
		v.privateKey = cryptolib.Key[cryptolib.KeyProviderPrivate]{}
		c.cfg.Keyrings[k] = v
	}
}

// decryptKeysEncrypted will decrypt all cfg.keysEncrypted.
func (c *Config) decryptKeysEncrypted(ctx context.Context) {
	ctx = c.lock(ctx)
	defer c.unlock(ctx)

	keys := c.cfg.keys.Keys()

	for i := range c.cfg.keysEncrypted {
		out, err := c.cfg.keysEncrypted[i].Decrypt(keys)
		if err == nil {
			var k cryptolib.Key[cryptolib.KeyProviderPrivate]

			k, err = cryptolib.ParseKey[cryptolib.KeyProviderPrivate](string(out))
			if err == nil {
				c.cfg.keys = append(c.cfg.keys, k)
			}
		}

		if err != nil {
			logger.Debug(ctx, errs.ErrReceiver.Wrap(err).Error())
		}
	}

	c.cfg.keysEncrypted = nil
}

// decryptPrivateKey will use all keys in cfg.keys to decrypt Keyring private keys until one is found is found.
func (c *Config) decryptKeyring(ctx context.Context, keyring KeyringName) (*Keyring20240410, error) {
	ctx = c.lock(ctx)
	defer c.unlock(ctx)

	if len(c.cfg.Keyrings) == 0 {
		// Don't log this, it may be expected
		return nil, ErrNotInitialized
	}

	k, err := c.getKeyring(ctx, keyring)
	if err != nil {
		return nil, logger.Error(ctx, err)
	}

	// Already decrypted?
	if !k.privateKey.IsNil() {
		return &k, nil
	}

	var out []byte

	keys := c.cfg.keys.Keys()

	// Try to decrypt without decrypting encrypted keys.
	for i := range c.cfg.DecryptKeys {
		if p, ok := c.cfg.DecryptKeys[i].PrivateKeys[keyring]; ok {
			out, err = p.PrivateKey.Decrypt(keys)
			if err == nil {
				key, err := cryptolib.ParseKey[cryptolib.KeyProviderPrivate](string(out))
				if err == nil {
					k.privateKey = key

					break
				}
			}
		}
	}

	// That didn't work, decrypt the keys and try again.
	if k.privateKey.IsNil() {
		c.decryptKeysEncrypted(ctx)
		keys := c.cfg.keys.Keys()

		for i := range c.cfg.DecryptKeys {
			if p, ok := c.cfg.DecryptKeys[i].PrivateKeys[c.cfg.Keyring]; ok {
				out, err = p.PrivateKey.Decrypt(keys)
				if err == nil {
					key, err := cryptolib.ParseKey[cryptolib.KeyProviderPrivate](string(out))
					if err == nil {
						k.privateKey = key

						break
					}
				}
			}
		}
	}

	c.cfg.Keyrings[keyring] = k

	if err != nil {
		return nil, logger.Error(ctx, err)
	}

	return &k, logger.Error(ctx, nil)
}

func (c *Config) getDecryptKey(ctx context.Context, name string) (DecryptKey20240410, error) {
	c.rlock(ctx)
	defer c.runlock(ctx)

	d, ok := c.cfg.DecryptKeys[name]
	if !ok {
		return DecryptKey20240410{}, logger.Error(ctx, fmt.Errorf("%w: %s", ErrDecryptKeyNotFound, name))
	}

	if d.PrivateKeys == nil {
		d.PrivateKeys = map[KeyringName]DecryptKeyPrivateKey20240410{}
	}

	return d, logger.Error(ctx, nil)
}

type ctxKey string

func (c *Config) getKeyring(ctx context.Context, keyring KeyringName) (Keyring20240410, error) {
	c.rlock(ctx)

	defer c.runlock(ctx)

	k, ok := c.cfg.Keyrings[keyring]
	if !ok {
		return Keyring20240410{}, logger.Error(ctx, fmt.Errorf("%w: %s", ErrKeyringNotFound, keyring))
	}

	if k.Values == nil {
		k.Values = map[string]Value20231210{}
	}

	return k, logger.Error(ctx, nil)
}

func addLock(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey("lock"), getLocks(ctx)+1)
}

func getLocks(ctx context.Context) int {
	if v, ok := ctx.Value(ctxKey("lock")).(int); ok {
		return v
	}

	return 0
}

func (c *Config) lock(ctx context.Context) context.Context {
	if getLocks(ctx) == 0 {
		c.mu.Lock()
	}

	return addLock(ctx)
}

func (c *Config) rlock(ctx context.Context) {
	if getLocks(ctx) > 0 {
		return
	}

	c.mu.RLock()
}
func (c *Config) runlock(ctx context.Context) {
	if getLocks(ctx) > 0 {
		return
	}

	c.mu.RUnlock()
}

func (c *Config) save(ctx context.Context, err error) error {
	if err != nil {
		return logger.Error(ctx, err)
	}

	if getLocks(ctx) > 1 {
		return logger.Error(ctx, nil)
	}

	c.rlock(ctx)

	defer c.runlock(ctx)

	out, err := config.Mask(ctx, c.cfg, []string{"cli", "keys", "keyrings.*.privateKey", "keyPath", "unmask"})
	if err != nil {
		return logger.Error(ctx, err)
	}

	s, err := jsonnet.Convert(ctx, out)
	if err != nil {
		return logger.Error(ctx, err)
	}

	if err := os.WriteFile(c.cfg.CLI.ConfigPath, []byte(s), 0600); err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(errors.New("error writing file"), err))
	}

	return logger.Error(ctx, nil)
}

func (c *Config) unlock(ctx context.Context) {
	if getLocks(ctx) > 1 {
		return
	}

	c.mu.Unlock()
}

func (k *Keyring20240410) getValueDecrypted(ctx context.Context, name string) ([]byte, error) {
	if v, ok := k.Values[name]; ok {
		kb, err := k.privateKey.Key.DecryptAsymmetric(v.Key)
		if err != nil {
			return nil, logger.Error(ctx, err)
		}

		k, err := cryptolib.ParseKey[cryptolib.KeyProviderSymmetric](string(kb))
		if err != nil {
			return nil, logger.Error(ctx, err)
		}

		out, err := k.Key.DecryptSymmetric(v.Value)
		if err != nil {
			return nil, logger.Error(ctx, err)
		}

		return out, logger.Error(ctx, nil)
	}

	return nil, logger.Error(ctx, fmt.Errorf("%w: %s", ErrValueNotFound, name))
}
