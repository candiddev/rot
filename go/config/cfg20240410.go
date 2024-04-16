package config

import (
	"time"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/types"
)

const cfgVersion20240410 cfgVersion = "2024.04.10"

type cfg20240410 struct {
	Algorithms  Algorithms20231210              `json:"algorithms"`
	CLI         cli.Config                      `json:"cli"`
	DecryptKeys map[string]DecryptKey20240410   `json:"decryptKeys"`
	Keyring     KeyringName                     `json:"keyring"`
	Keyrings    map[KeyringName]Keyring20240410 `json:"keyrings"`
	Keys        types.SliceString               `json:"keys"`
	KeyPath     string                          `json:"keyPath"`
	License     License20240410                 `json:"license"`
	LicenseKey  string                          `json:"licenseKey"`
	Unmask      types.SliceString               `json:"unmask"`
	Version     cfgVersion                      `json:"version"`

	keys          cryptolib.Keys[cryptolib.KeyProviderPrivate] //nolint:revive
	keysEncrypted cryptolib.EncryptedValues
}

// DecryptKey20240410 contains config values.
type DecryptKey20240410 struct {
	Modified    time.Time                                    `json:"modified"`
	PrivateKeys map[KeyringName]DecryptKeyPrivateKey20240410 `json:"privateKeys"`
	PublicKey   cryptolib.Key[cryptolib.KeyProviderPublic]   `json:"publicKey"`
}

// DecryptKeyPrivateKey20240410 contains config values.
type DecryptKeyPrivateKey20240410 struct {
	PrivateKey cryptolib.EncryptedValue `json:"privateKey"`
	Signature  cryptolib.Signature      `json:"signature"`
}

// Keyring20240410 contains config values.
type Keyring20240410 struct {
	PrivateKey cryptolib.Key[cryptolib.KeyProviderPrivate] `json:"privateKey,omitempty"`
	PublicKey  cryptolib.Key[cryptolib.KeyProviderPublic]  `json:"publicKey"`
	Values     map[string]Value20231210                    `json:"values,omitempty"`

	privateKey cryptolib.Key[cryptolib.KeyProviderPrivate] //nolint:revive
}

// License20240410 configures licensing.
type License20240410 struct {
	Keyrings     int    `json:"keyrings"`
	Organization string `json:"sub"`
}

// Default generates a Config with known good defaults.
func default20240410() *cfg20240410 {
	return &cfg20240410{
		Algorithms: Algorithms20231210{
			Asymmetric: "best",
			PBKDF:      "best",
			Symmetric:  "best",
		},
		CLI:         cli.Config{},
		DecryptKeys: map[string]DecryptKey20240410{},
		Keyrings:    map[KeyringName]Keyring20240410{},
		KeyPath:     ".rot-keys",
		License: License20240410{
			Keyrings: 2,
		},
		Version: cfgVersion20240410,
	}
}
