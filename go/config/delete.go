package config

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/candiddev/shared/go/logger"
)

var ErrDecryptKeyPrivateKeyNotFound = errors.New("decryptKey privateKey not found")

// DeleteDecryptKey will remove a DecryptKey from the config.
func (c *Config) DeleteDecryptKey(ctx context.Context, name string) error {
	ctx = c.lock(ctx)
	defer c.unlock(ctx)

	if _, err := c.getDecryptKey(ctx, name); err != nil {
		return logger.Error(ctx, err)
	}

	delete(c.cfg.DecryptKeys, name)

	return logger.Error(ctx, c.save(ctx, nil))
}

// DeleteDecryptKeyPrivateKey will remove a DecryptKey.PrivateKey from the config.
func (c *Config) DeleteDecryptKeyPrivateKey(ctx context.Context, keyring KeyringName, name string) error {
	ctx = c.lock(ctx)
	defer c.unlock(ctx)

	d, err := c.getDecryptKey(ctx, name)
	if err != nil {
		return logger.Error(ctx, err)
	}

	if _, ok := d.PrivateKeys[keyring]; !ok {
		return logger.Error(ctx, fmt.Errorf("%w: decryptKey: %s, keyring: %s", ErrDecryptKeyPrivateKeyNotFound, name, keyring))
	}

	delete(d.PrivateKeys, keyring)

	d.Modified = time.Now()
	c.cfg.DecryptKeys[name] = d

	return logger.Error(ctx, c.save(ctx, nil))
}

// DeleteKeyring will remove a Keyring from the config.
func (c *Config) DeleteKeyring(ctx context.Context, keyring KeyringName) error {
	ctx = c.lock(ctx)
	defer c.unlock(ctx)

	_, err := c.getKeyring(ctx, keyring)
	if err != nil {
		return logger.Error(ctx, err)
	}

	for k := range c.cfg.DecryptKeys {
		if err := c.DeleteDecryptKeyPrivateKey(ctx, keyring, k); err != nil && !errors.Is(err, ErrDecryptKeyPrivateKeyNotFound) {
			return logger.Error(ctx, err)
		}
	}

	delete(c.cfg.Keyrings, keyring)

	return logger.Error(ctx, c.save(ctx, nil))
}

// DeleteKeyringValue will remove a Keyring.Value from the config.
func (c *Config) DeleteKeyringValue(ctx context.Context, keyring KeyringName, name string) error {
	ctx = c.lock(ctx)
	defer c.unlock(ctx)

	k, err := c.getKeyring(ctx, keyring)
	if err != nil {
		return logger.Error(ctx, err)
	}

	if _, ok := k.Values[name]; !ok {
		return logger.Error(ctx, ErrValueNotFound)
	}

	delete(k.Values, name)

	c.cfg.Keyrings[keyring] = k

	return logger.Error(ctx, c.save(ctx, nil))
}
