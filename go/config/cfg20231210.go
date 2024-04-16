// Package config contains tools for managing a Rot config.
package config

import (
	"context"
	"fmt"
	"time"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/logger"
	"github.com/candiddev/shared/go/types"
)

const cfgVersion20231210 cfgVersion = "2023.12.10"

type cfg20231210 struct {
	Algorithms  Algorithms20231210                          `json:"algorithms"`
	CLI         cli.Config                                  `json:"cli"`
	DecryptKeys map[string]DecryptKey20231210               `json:"decryptKeys"`
	Keys        types.SliceString                           `json:"keys"`
	KeyPath     string                                      `json:"keyPath"`
	PrivateKey  cryptolib.Key[cryptolib.KeyProviderPrivate] `json:"privateKey,omitempty"`
	PublicKey   cryptolib.Key[cryptolib.KeyProviderPublic]  `json:"publicKey"`
	Unmask      types.SliceString                           `json:"unmask"`
	Values      map[string]Value20231210                    `json:"values"`
}

// Algorithms20231210 contains config values.
type Algorithms20231210 struct {
	Asymmetric cryptolib.Encryption `json:"asymmetric"`
	PBKDF      cryptolib.KDF        `json:"pbkdf"`
	Symmetric  cryptolib.Encryption `json:"symmetric"`
}

// DecryptKey20231210 contains config values.
type DecryptKey20231210 struct {
	Modified   time.Time                                  `json:"modified"`
	PrivateKey cryptolib.EncryptedValue                   `json:"privateKey"`
	PublicKey  cryptolib.Key[cryptolib.KeyProviderPublic] `json:"publicKey"`
	Signature  cryptolib.Signature                        `json:"signature"`
}

// Value20231210 contains config values.
type Value20231210 struct {
	Comment  string                   `json:"comment"`
	Key      cryptolib.EncryptedValue `json:"key"`
	Modified time.Time                `json:"modified"`
	Value    cryptolib.EncryptedValue `json:"value"`
}

func default20231210() *cfg20231210 {
	return &cfg20231210{
		Algorithms: Algorithms20231210{
			Asymmetric: "best",
			PBKDF:      "best",
			Symmetric:  "best",
		},
		CLI:         cli.Config{},
		DecryptKeys: map[string]DecryptKey20231210{},
		KeyPath:     ".rot-keys",
		Values:      map[string]Value20231210{},
	}
}

func (c *cfg20231210) upgrade(ctx context.Context) *cfg20240410 {
	logger.Info(ctx, fmt.Sprintf("Upgrading config version from %s to %s", cfgVersion20231210, cfgVersion20240410))

	cn := default20240410()
	cn.Algorithms = c.Algorithms
	cn.CLI = c.CLI

	cn.Keys = c.Keys
	cn.KeyPath = c.KeyPath

	if !c.PublicKey.IsNil() {
		cn.Keyring = "rot"

		cn.Keyrings[cn.Keyring] = Keyring20240410{
			PrivateKey: c.PrivateKey,
			PublicKey:  c.PublicKey,
			Values:     c.Values,
		}
	}

	for k, v := range c.DecryptKeys {
		d := DecryptKey20240410{
			Modified: v.Modified,
			PrivateKeys: map[KeyringName]DecryptKeyPrivateKey20240410{
				cn.Keyring: {
					PrivateKey: v.PrivateKey,
					Signature:  v.Signature,
				},
			},
			PublicKey: v.PublicKey,
		}

		cn.DecryptKeys[k] = d
	}

	cn.Unmask = c.Unmask

	return cn
}
