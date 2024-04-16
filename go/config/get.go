package config

import (
	"context"
	"errors"
	"os"
	"sort"
	"strings"

	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
	"github.com/candiddev/shared/go/types"
)

// GetAlgorithms returns all Algorithms.
func (c *Config) GetAlgorithms(ctx context.Context) Algorithms20231210 {
	c.rlock(ctx)
	defer c.runlock(ctx)

	return c.cfg.Algorithms
}

// GetDecryptKeys returns a list of decryptKeys.
func (c *Config) GetDecryptKeys(ctx context.Context) []string {
	c.rlock(ctx)
	defer c.runlock(ctx)

	s := []string{}

	for i := range c.cfg.DecryptKeys {
		s = append(s, i)
	}

	sort.Strings(s)

	return s
}

// GetDecryptKeysKeyring returns a list of decryptKeys with access to a keyring.
func (c *Config) GetDecryptKeysKeyring(ctx context.Context, keyring KeyringName) ([]string, error) {
	c.rlock(ctx)
	defer c.runlock(ctx)

	s := []string{}

	if _, err := c.getKeyring(ctx, keyring); err != nil {
		return nil, logger.Error(ctx, err)
	}

	for i := range c.cfg.DecryptKeys {
		if _, ok := c.cfg.DecryptKeys[i].PrivateKeys[keyring]; ok {
			s = append(s, i)
		}
	}

	sort.Strings(s)

	return s, logger.Error(ctx, nil)
}

// GetDecrypted will decrypt an EncryptedValue using Keys.
func (c *Config) GetDecrypted(ctx context.Context, ev cryptolib.EncryptedValue) ([]byte, error) {
	if ev.KDF == "" {
		c.decryptKeysEncrypted(ctx)
	}

	o, err := ev.Decrypt(c.cfg.keys.Keys())

	return o, logger.Error(ctx, err)
}

// GetKeys returns all of the decrypt User Private Keys.
func (c *Config) GetKeys(ctx context.Context) cryptolib.Keys[cryptolib.KeyProviderPrivate] {
	c.decryptKeysEncrypted(ctx)

	return c.cfg.keys
}

// GetKeyringName will return the default keyring name.
func (c *Config) GetKeyringName(ctx context.Context) KeyringName {
	c.rlock(ctx)
	defer c.runlock(ctx)

	return c.cfg.Keyring
}

// GetKeyrings returns all of the keyrings.
func (c *Config) GetKeyrings(ctx context.Context) []string {
	c.rlock(ctx)
	defer c.runlock(ctx)

	v := []string{}

	for k := range c.cfg.Keyrings {
		if k == c.cfg.Keyring {
			k += " (current)"
		}

		v = append(v, string(k))
	}

	sort.Strings(v)

	return v
}

// GetKeyringValues will return value names from a keyring.
func (c *Config) GetKeyringValues(ctx context.Context, name KeyringName) ([]string, error) {
	k, err := c.getKeyring(ctx, name)
	if err != nil {
		return nil, logger.Error(ctx, err)
	}

	n := []string{}

	for v := range k.Values {
		n = append(n, v)
	}

	sort.Strings(n)

	return n, logger.Error(ctx, nil)
}

// GetPrivateKey will return a private key from a string, value, or file.
func (c *Config) GetPrivateKey(ctx context.Context, keyring KeyringName, value string) (cryptolib.Key[cryptolib.KeyProviderPrivate], error) {
	pk := value
	if len(strings.Split(pk, ":")) == 1 {
		v, err := c.GetValueDecrypted(ctx, keyring, pk)
		if err == nil {
			pk = string(v)
		} else {
			f, err := os.ReadFile(value)
			if err == nil {
				pk = string(f)
			}
		}
	}

	key, err := cryptolib.ParseKey[cryptolib.KeyProviderPrivate](pk)
	if err != nil {
		return key, logger.Error(ctx, errs.ErrReceiver.Wrap(errors.New("error parsing private key"), err))
	}

	return key, logger.Error(ctx, nil)
}

// GetPublicKey will return a public key from a string, value or file.
func (c *Config) GetPublicKey(ctx context.Context, keyring KeyringName, value string) (cryptolib.Key[cryptolib.KeyProviderPublic], error) {
	k, err := c.getKeyring(ctx, keyring)
	if err != nil {
		return cryptolib.Key[cryptolib.KeyProviderPublic]{}, logger.Error(ctx, err)
	}

	pk := value
	if len(strings.Split(pk, ":")) == 1 {
		pk = k.Values[value].Comment

		if pk == "" {
			f, err := os.ReadFile(value)
			if err == nil {
				pk = string(f)
			}
		}
	}

	key, er := cryptolib.ParseKey[cryptolib.KeyProviderPublic](pk)
	if er != nil {
		return key, logger.Error(ctx, errs.ErrReceiver.Wrap(errors.New("error parsing public key"), er))
	}

	return key, logger.Error(ctx, nil)
}

// GetUnmask returns the unmask value.
func (c *Config) GetUnmask() types.SliceString {
	return c.cfg.Unmask
}

func (c *Config) GetValue(ctx context.Context, keyring KeyringName, name string) (Value20231210, error) {
	c.rlock(ctx)
	defer c.runlock(ctx)

	k, err := c.getKeyring(ctx, keyring)
	if err != nil {
		return Value20231210{}, logger.Error(ctx, err)
	}

	v, ok := k.Values[name]
	if !ok {
		return Value20231210{}, logger.Error(ctx, ErrValueNotFound)
	}

	return v, logger.Error(ctx, nil)
}

// GetValueDecrypted will retrieve a decrypted Value from a Keyring.
func (c *Config) GetValueDecrypted(ctx context.Context, keyring KeyringName, name string) ([]byte, error) {
	k, err := c.decryptKeyring(ctx, keyring)
	if err != nil {
		return nil, logger.Error(ctx, err)
	}

	o, err := k.getValueDecrypted(ctx, name)

	return o, logger.Error(ctx, err)
}
