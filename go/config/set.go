package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
	"github.com/candiddev/shared/go/types"
)

var (
	ErrInitCancelled = errors.New("init cancelled, config already exists")
	ErrKeyNotFound   = errors.New("no key found with provided ID")
	ErrNoDecryptKeys = errors.New("cannot add keyring without decrypt keys")
)

func (c *Config) Init(ctx context.Context, keyring KeyringName, name string) error {
	ctx = c.lock(ctx)
	defer c.unlock(ctx)

	if _, err := os.ReadFile(c.cfg.CLI.ConfigPath); err == nil {
		b, err := logger.Prompt(c.cfg.CLI.ConfigPath+" aleady exists, overwite (yes/no)?", "", false)
		if err != nil {
			return logger.Error(ctx, err)
		}

		if string(b) != "yes" {
			return logger.Error(ctx, ErrInitCancelled)
		}
	}

	c.decryptKeysEncrypted(ctx)

	var err error

	var pub cryptolib.Key[cryptolib.KeyProviderPublic]

	var prv cryptolib.Key[cryptolib.KeyProviderPrivate]

	for i := range c.cfg.keys {
		if c.cfg.keys[i].ID == name {
			prv = c.cfg.keys[i]

			break
		}
	}

	if name == "" && len(c.cfg.keys) > 0 {
		prv = c.cfg.keys[0]
	} else if prv.IsNil() {
		pub, err = c.NewKeyPathKey(ctx, name)
	}

	if !prv.IsNil() {
		p, err := prv.Key.Public()
		if err != nil {
			return logger.Error(ctx, ErrKeyNotFound)
		}

		pub = cryptolib.Key[cryptolib.KeyProviderPublic]{
			ID:  prv.ID,
			Key: p,
		}
	}

	if err != nil {
		return logger.Error(ctx, err)
	}

	if err := c.SetDecryptKey(ctx, pub.ID, pub); err != nil {
		return logger.Error(ctx, err)
	}

	return logger.Error(ctx, c.save(ctx, c.NewKeyring(ctx, keyring, []string{pub.ID})))
}

// NewKeyring will add a new Keyring or modify an existing one and add DecryptKeys to it.
func (c *Config) NewKeyring(ctx context.Context, name KeyringName, decryptKeys []string) error {
	ctx = c.lock(ctx)
	defer c.unlock(ctx)

	if l := len(c.GetKeyrings(ctx)); c.cfg.License.Keyrings != 0 && l >= c.cfg.License.Keyrings {
		return logger.Error(ctx, fmt.Errorf("adding a new keyring would exceed the license limit (%d), please upgrade your license", l))
	}

	if _, err := c.decryptKeyring(ctx, name); err != nil {
		if len(decryptKeys) == 0 {
			return logger.Error(ctx, ErrNoDecryptKeys)
		}

		prv, pub, err := cryptolib.NewKeysAsymmetric(c.cfg.Algorithms.Asymmetric)
		if err != nil {
			return logger.Error(ctx, errs.ErrReceiver.Wrapf("error generating keyring keys: %w", err))
		}

		c.cfg.Keyrings[name] = Keyring20240410{
			PublicKey:  pub,
			Values:     map[string]Value20231210{},
			privateKey: prv,
		}
	}

	for i := range decryptKeys {
		if err := c.SetDecryptKeyPrivateKey(ctx, name, decryptKeys[i]); err != nil {
			return logger.Error(ctx, err)
		}
	}

	if c.cfg.Keyring == "" {
		c.cfg.Keyring = name
	}

	return logger.Error(ctx, c.save(ctx, nil))
}

// NewKeyPathKey adds a Key to KeyPath.
func (c *Config) NewKeyPathKey(ctx context.Context, name string) (cryptolib.Key[cryptolib.KeyProviderPublic], error) {
	var err error

	var pub cryptolib.Key[cryptolib.KeyProviderPublic]

	var prv cryptolib.Key[cryptolib.KeyProviderPrivate]

	prv, pub, err = cryptolib.NewKeysAsymmetric(c.GetAlgorithms(ctx).Asymmetric)
	if err != nil {
		return pub, logger.Error(ctx, err)
	}

	if name != "" {
		prv.ID = name
		pub.ID = name
	}

	v, err := cryptolib.KDFSet(cryptolib.Argon2ID, prv.ID, []byte(prv.String()), c.GetAlgorithms(ctx).Symmetric)
	if err != nil {
		return pub, logger.Error(ctx, err)
	}

	var p string

	if v.Ciphertext == "" {
		p = prv.String()
	} else {
		p = v.String()
	}

	f, err := os.OpenFile(c.cfg.KeyPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return pub, logger.Error(ctx, errs.ErrReceiver.Wrap(errors.New("error opening keys file"), err))
	}

	defer f.Close()

	if _, err := f.WriteString(p + "\n"); err != nil {
		return pub, logger.Error(ctx, errs.ErrReceiver.Wrap(errors.New("error writing keys file"), err))
	}

	return pub, logger.Error(ctx, nil)
}

// Rekey will rekey a Keyring.
func (c *Config) Rekey(ctx context.Context, keyring KeyringName) error {
	ctx = c.lock(ctx)
	defer c.unlock(ctx)

	k, err := c.decryptKeyring(ctx, keyring)
	if err != nil {
		return logger.Error(ctx, err)
	}

	d := []string{}

	for i := range c.cfg.DecryptKeys {
		for j := range c.cfg.DecryptKeys[i].PrivateKeys {
			if j == keyring {
				d = append(d, i)

				break
			}
		}
	}

	if err := c.DeleteKeyring(ctx, keyring); err != nil {
		return logger.Error(ctx, err)
	}

	if err := c.NewKeyring(ctx, keyring, d); err != nil {
		return logger.Error(ctx, err)
	}

	for i := range k.Values {
		o, err := k.getValueDecrypted(ctx, i)
		if err != nil {
			return logger.Error(ctx, err)
		}

		if err := c.SetValue(ctx, keyring, o, i, k.Values[i].Comment); err != nil {
			return logger.Error(ctx, err)
		}
	}

	return logger.Error(ctx, c.save(ctx, nil))
}

// SetDecryptKey sets a DecryptKey to the config.
func (c *Config) SetDecryptKey(ctx context.Context, name string, publicKey cryptolib.Key[cryptolib.KeyProviderPublic]) error {
	ctx = c.lock(ctx)
	defer c.unlock(ctx)

	c.cfg.DecryptKeys[name] = DecryptKey20240410{
		Modified:  time.Now(),
		PublicKey: publicKey,
	}

	return logger.Error(ctx, c.save(ctx, nil))
}

// SetDecryptKeyPrivateKey sets a DecryptKey Keyring Private Key to the config.
func (c *Config) SetDecryptKeyPrivateKey(ctx context.Context, keyring KeyringName, name string) error {
	ctx = c.lock(ctx)
	defer c.unlock(ctx)

	ky, err := c.getDecryptKey(ctx, name)
	if err != nil {
		return logger.Error(ctx, err)
	}

	kr, err := c.decryptKeyring(ctx, keyring)
	if err != nil {
		return logger.Error(ctx, err)
	}

	v, err := ky.PublicKey.Key.EncryptAsymmetric([]byte(kr.privateKey.String()), ky.PublicKey.ID, c.cfg.Algorithms.Symmetric)
	if err != nil {
		return logger.Error(ctx, err)
	}

	sig, err := cryptolib.NewSignature(kr.privateKey, []byte(ky.PublicKey.String()))
	if err != nil {
		return logger.Error(ctx, err)
	}

	ky.Modified = time.Now()
	ky.PrivateKeys[keyring] = DecryptKeyPrivateKey20240410{
		PrivateKey: v,
		Signature:  sig,
	}
	c.cfg.DecryptKeys[name] = ky

	return logger.Error(ctx, c.save(ctx, nil))
}

// SetValue will encrypt a value and add it to a keyring in the config.
func (c *Config) SetValue(ctx context.Context, keyring KeyringName, value []byte, name, comment string) error {
	ctx = c.lock(ctx)
	defer c.unlock(ctx)

	if err := types.EnvValidate(name); err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(errors.New("invalid name"), err))
	}

	key, err := cryptolib.NewKeySymmetric(c.cfg.Algorithms.Symmetric)
	if err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
	}

	ev, err := key.Key.EncryptSymmetric(value, key.ID)
	if err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
	}

	k, err := c.getKeyring(ctx, keyring)
	if err != nil {
		return logger.Error(ctx, err)
	}

	ek, err := k.PublicKey.Key.EncryptAsymmetric([]byte(key.String()), k.PublicKey.ID, c.cfg.Algorithms.Symmetric)
	if err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
	}

	k.Values[name] = Value20231210{
		Comment:  comment,
		Key:      ek,
		Modified: time.Now(),
		Value:    ev,
	}
	c.cfg.Keyrings[keyring] = k

	return logger.Error(ctx, c.save(ctx, nil))
}
