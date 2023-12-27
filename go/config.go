package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/config"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/jsonnet"
	"github.com/candiddev/shared/go/logger"
)

var (
	errUnknownAlgorithmsAsymmetric = errors.New("unknown algorithm.asymmetric")
	errUnknownAlgorithmsPBKDF      = errors.New("unknown algorithm.pbkdf")
	errUnknownAlgorithmsSymmetric  = errors.New("unknown algorithm.symmetric")
)

type cfg struct {
	Algorithms  cfgAlgorithms                               `json:"algorithms"`
	CLI         cli.Config                                  `json:"cli"`
	DecryptKeys map[string]cfgDecryptKey                    `json:"decryptKeys"`
	Keys        []string                                    `json:"keys"`
	KeyPath     string                                      `json:"keyPath"`
	PrivateKey  cryptolib.Key[cryptolib.KeyProviderPrivate] `json:"privateKey,omitempty"`
	PublicKey   cryptolib.Key[cryptolib.KeyProviderPublic]  `json:"publicKey"`
	Unmask      []string                                    `json:"unmask"`
	Values      map[string]cfgValue                         `json:"values"`

	keys          cryptolib.Keys[cryptolib.KeyProviderPrivate] //nolint:revive
	keysEncrypted cryptolib.EncryptedValues
	privateKey    cryptolib.Key[cryptolib.KeyProviderPrivate] //nolint:revive
}

type cfgAlgorithms struct {
	Asymmetric cryptolib.Encryption `json:"asymmetric"`
	PBKDF      cryptolib.KDF        `json:"pbkdf"`
	Symmetric  cryptolib.Encryption `json:"symmetric"`
}

type cfgDecryptKey struct {
	Modified   time.Time                                  `json:"modified"`
	PrivateKey cryptolib.EncryptedValue                   `json:"privateKey"`
	PublicKey  cryptolib.Key[cryptolib.KeyProviderPublic] `json:"publicKey"`
	Signature  cryptolib.Signature                        `json:"signature"`
}

type cfgValue struct {
	Comment  string                   `json:"comment"`
	Key      cryptolib.EncryptedValue `json:"key"`
	Modified time.Time                `json:"modified"`
	Value    cryptolib.EncryptedValue `json:"value"`
}

// defaultCfg generates a cfg with known good defaults.
func defaultCfg() *cfg {
	return &cfg{
		Algorithms: cfgAlgorithms{
			Asymmetric: "best",
			PBKDF:      "best",
			Symmetric:  "best",
		},
		CLI:         cli.Config{},
		DecryptKeys: map[string]cfgDecryptKey{},
		KeyPath:     ".rot-keys",
		Values:      map[string]cfgValue{},
	}
}

// save generates a jsonnet string of cfg and writes it to path.
func (c *cfg) save(ctx context.Context) errs.Err {
	out, err := config.Mask(ctx, c, []string{"cli", "key", "keys", "keyPath", "unmask", "privateKey"})
	if err != nil {
		return logger.Error(ctx, err)
	}

	s, err := jsonnet.Convert(ctx, out)
	if err != nil {
		return logger.Error(ctx, err)
	}

	if err := os.WriteFile(c.CLI.ConfigPath, []byte(s), 0600); err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(errors.New("error writing file"), err))
	}

	return logger.Error(ctx, nil)
}

func (c *cfg) CLIConfig() *cli.Config {
	return &c.CLI
}

func (c *cfg) Parse(ctx context.Context, configArgs []string) errs.Err { //nolint:gocognit
	if err := config.Parse(ctx, c, configArgs, "ROT", c.CLI.ConfigPath); err != nil {
		err = logger.Error(ctx, err)

		if !err.Is(jsonnet.ErrImport) {
			return err
		}
	}

	for i := range cryptolib.EncryptionAsymmetric {
		if string(c.Algorithms.Asymmetric) == cryptolib.EncryptionAsymmetric[i] {
			break
		} else if i == len(cryptolib.EncryptionAsymmetric)-1 {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(errUnknownAlgorithmsAsymmetric, fmt.Errorf("%s, valid values are [%s]", c.Algorithms.Asymmetric, strings.Join(cryptolib.EncryptionAsymmetric, " "))))
		}
	}

	for i := range cryptolib.ValidPBKDF {
		if string(c.Algorithms.PBKDF) == cryptolib.ValidPBKDF[i] {
			break
		} else if i == len(cryptolib.ValidPBKDF)-1 {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(errUnknownAlgorithmsPBKDF, fmt.Errorf("%s, valid values are [%s]", c.Algorithms.PBKDF, strings.Join(cryptolib.ValidPBKDF, " "))))
		}
	}

	for i := range cryptolib.EncryptionSymmetric {
		if string(c.Algorithms.Symmetric) == cryptolib.EncryptionSymmetric[i] {
			break
		} else if i == len(cryptolib.EncryptionSymmetric)-1 {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(errUnknownAlgorithmsSymmetric, fmt.Errorf("%s, valid values are [%s]", c.Algorithms.Symmetric, strings.Join(cryptolib.EncryptionSymmetric, " "))))
		}
	}

	if c.KeyPath != "" {
		p := c.KeyPath

		if filepath.Base(c.KeyPath) == c.KeyPath {
			f := config.FindPathAscending(ctx, c.KeyPath)
			if f != "" {
				p = f
			}
		}

		out, err := os.ReadFile(p)
		if err == nil {
			for _, s := range strings.Split(string(out), "\n") {
				if s != "" {
					c.Keys = append(c.Keys, s)
				}
			}
		}
	}

	for i, key := range c.Keys {
		// Check if it's an encrypted value
		ev, err := cryptolib.ParseEncryptedValue(key)
		if err == nil {
			c.keysEncrypted = append(c.keysEncrypted, ev)
		} else {
			// Try parsing the key
			k, err := cryptolib.ParseKey[cryptolib.KeyProviderPrivate](key)
			if err == nil {
				c.keys = append(c.keys, k)
			} else {
				logger.Error(ctx, errs.ErrReceiver.Wrap(fmt.Errorf("error parsing key %d", i+1), err)) //nolint:errcheck
			}
		}

		c.Keys = []string{}
	}

	tamper := false

	for i := range c.DecryptKeys {
		if err := c.DecryptKeys[i].Signature.Verify([]byte(c.DecryptKeys[i].PublicKey.String()), cryptolib.Keys[cryptolib.KeyProviderPublic]{
			c.PublicKey,
		}); err != nil {
			tamper = true

			logger.Error(ctx, errs.ErrReceiver.Wrap(fmt.Errorf("decryptKey %s has an invalid signature, this could indicate tampering: skipping this decryptKey", i))) //nolint:errcheck

			delete(c.DecryptKeys, i)
		}
	}

	if tamper && len(c.DecryptKeys) == 0 {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(errors.New("no valid decryptKeys, fix tampering or reinitialize")))
	}

	if !c.PrivateKey.IsNil() {
		c.privateKey = c.PrivateKey
	}

	return logger.Error(ctx, nil)
}
