package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/config"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/jsonnet"
	"github.com/candiddev/shared/go/logger"
)

// KeyringName is the name of a Keyring.
type KeyringName string

type cfgVersion string

var (
	ErrParsingConfig               = errors.New("error prasing config")
	ErrUnknownAlgorithmsAsymmetric = errors.New("unknown algorithm.asymmetric")
	ErrUnknownAlgorithmsPBKDF      = errors.New("unknown algorithm.pbkdf")
	ErrUnknownAlgorithmsSymmetric  = errors.New("unknown algorithm.symmetric")
	ErrTamper                      = errors.New("no valid decryptKeys, fix tampering or reinitialize")
)

// Config is the Rot config.
type Config struct {
	cfg *cfg20240410
	mu  sync.RWMutex
}

func Default() *Config {
	return &Config{
		cfg: default20240410(),
	}
}

func (c *Config) CLIConfig() *cli.Config {
	return &c.cfg.CLI
}

// Show returns the config.
func (c *Config) Show() any {
	return c.cfg
}

func (c *Config) Parse(ctx context.Context, configArgs []string) error { //nolint:gocognit,gocyclo
	m := map[string]any{}
	if err := config.Parse(ctx, &m, configArgs, "ROT", c.cfg.CLI.ConfigPath); err != nil {
		err = logger.Error(ctx, err)

		if !errors.Is(err, jsonnet.ErrImport) {
			return err
		}
	}

	b, err := json.Marshal(m)
	if err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(ErrParsingConfig, err))
	}

	a, ok := m["publicKey"]

	switch {
	case ok && a != "":
		p := c.cfg.CLI.ConfigPath
		cfg := default20231210()
		err = json.Unmarshal(b, &cfg)
		c.cfg = cfg.upgrade(ctx)
		c.cfg.CLI.ConfigPath = p
	default:
		err = json.Unmarshal(b, c.cfg)
	}

	if err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(ErrParsingConfig, err))
	}

	for i := range cryptolib.EncryptionAsymmetric {
		if string(c.cfg.Algorithms.Asymmetric) == cryptolib.EncryptionAsymmetric[i] {
			break
		} else if i == len(cryptolib.EncryptionAsymmetric)-1 {
			return errs.ErrReceiver.Wrap(ErrUnknownAlgorithmsAsymmetric, fmt.Errorf("%s, valid values are [%s]", c.cfg.Algorithms.Asymmetric, strings.Join(cryptolib.EncryptionAsymmetric, " ")))
		}
	}

	for i := range cryptolib.ValidPBKDF {
		if string(c.cfg.Algorithms.PBKDF) == cryptolib.ValidPBKDF[i] {
			break
		} else if i == len(cryptolib.ValidPBKDF)-1 {
			return errs.ErrReceiver.Wrap(ErrUnknownAlgorithmsPBKDF, fmt.Errorf("%s, valid values are [%s]", c.cfg.Algorithms.PBKDF, strings.Join(cryptolib.ValidPBKDF, " ")))
		}
	}

	for i := range cryptolib.EncryptionSymmetric {
		if string(c.cfg.Algorithms.Symmetric) == cryptolib.EncryptionSymmetric[i] {
			break
		} else if i == len(cryptolib.EncryptionSymmetric)-1 {
			return errs.ErrReceiver.Wrap(ErrUnknownAlgorithmsSymmetric, fmt.Errorf("%s, valid values are [%s]", c.cfg.Algorithms.Symmetric, strings.Join(cryptolib.EncryptionSymmetric, " ")))
		}
	}

	if _, ok := c.cfg.Keyrings[c.cfg.Keyring]; c.cfg.Keyring != "" && !ok {
		return logger.Error(ctx, fmt.Errorf("%w: %s", ErrKeyringNotFound, c.cfg.Keyring))
	}

	if c.cfg.KeyPath != "" {
		p := c.cfg.KeyPath

		if filepath.Base(c.cfg.KeyPath) == c.cfg.KeyPath {
			f := config.FindPathAscending(ctx, c.cfg.KeyPath)
			if f != "" {
				p = f
			}
		}

		out, err := os.ReadFile(p)
		if err == nil {
			for _, s := range strings.Split(string(out), "\n") {
				if s != "" {
					c.cfg.Keys = append(c.cfg.Keys, s)
				}
			}
		}
	}

	for i, key := range c.cfg.Keys {
		// Check if it's an encrypted value
		ev, err := cryptolib.ParseEncryptedValue(key)
		if err == nil {
			c.cfg.keysEncrypted = append(c.cfg.keysEncrypted, ev)
		} else {
			// Try parsing the key
			k, err := cryptolib.ParseKey[cryptolib.KeyProviderPrivate](key)
			if err == nil {
				c.cfg.keys = append(c.cfg.keys, k)
			} else {
				logger.Error(ctx, errs.ErrReceiver.Wrap(fmt.Errorf("error parsing key %d", i+1), err)) //nolint:errcheck
			}
		}

		c.cfg.Keys = []string{}
	}

	tamper := false

	for i := range c.cfg.DecryptKeys {
		for j := range c.cfg.DecryptKeys[i].PrivateKeys {
			match := false

			for k := range c.cfg.Keyrings {
				if j == k {
					match = true
				}
			}

			if !match {
				logger.Error(ctx, errs.ErrReceiver.Wrap(fmt.Errorf("decryptKey %s has a privateKey for an unknown keyring %s", i, j))) //nolint:errcheck

				continue
			}

			if err := c.cfg.DecryptKeys[i].PrivateKeys[j].Signature.Verify([]byte(c.cfg.DecryptKeys[i].PublicKey.String()), cryptolib.Keys[cryptolib.KeyProviderPublic]{
				c.cfg.Keyrings[j].PublicKey,
			}); err != nil {
				tamper = true

				logger.Error(ctx, fmt.Errorf("decryptKey %s has an invalid signature for keyring %s, this could indicate tampering: skipping this decryptKey", i, j)) //nolint:errcheck

				delete(c.cfg.DecryptKeys, i)
			}
		}
	}

	if tamper && len(c.cfg.DecryptKeys) == 0 {
		return logger.Error(ctx, ErrTamper)
	}

	for k, v := range c.cfg.Keyrings {
		if !v.PrivateKey.IsNil() {
			v.privateKey = v.PrivateKey

			c.cfg.Keyrings[k] = v
		}
	}

	l := License20240410{}

	if _, err := cli.ParseLicense(ctx, c.cfg.LicenseKey, &l); err == nil && c.cfg.LicenseKey != "" {
		c.cfg.License = l
	} else {
		logger.Error(ctx, err) //nolint:errcheck
	}

	if l.Keyrings != 0 && len(c.cfg.Keyrings) > l.Keyrings {
		return logger.Error(ctx, fmt.Errorf("number of keyrings (%d) exceeds the license amount (%d), please upgrade your license or reduce the number of keyrings", len(c.cfg.Keyrings), c.cfg.License.Keyrings))
	}

	return logger.Error(ctx, nil)
}
