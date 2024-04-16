package config

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/jsonnet"
	"github.com/candiddev/shared/go/logger"
)

var ctx = context.Background()

func TestConfigParse(t *testing.T) {
	logger.UseTestLogger(t)

	c20231210 := get20231210()
	c20240410 := c20231210.upgrade(ctx)

	o, _ := jsonnet.Convert(ctx, c20240410)

	os.WriteFile("testdata/cfg20240410.jsonnet", []byte(o), 0600)
	defer os.Remove("testdata/cfg20240410.jsonnet")

	badKeyring := Default()
	badKeyring.cfg.Keyring = "missing"

	d := c20240410.DecryptKeys["test2"]
	d.PublicKey = c20240410.DecryptKeys["test1"].PublicKey
	c20240410.DecryptKeys["test2"] = d
	c20240410.DecryptKeys["test3"] = DecryptKey20240410{
		PrivateKeys: map[KeyringName]DecryptKeyPrivateKey20240410{
			"rot": {},
		},
	}
	delete(c20240410.DecryptKeys, "test1")

	// Error testing
	tests := map[string]struct {
		c       *cfg20240410
		wantErr error
	}{
		"20231210": {
			c: &cfg20240410{
				CLI: cli.Config{
					ConfigPath: "testdata/cfg20231210.jsonnet",
				},
			},
		},
		"20240410": {
			c: &cfg20240410{
				CLI: cli.Config{
					ConfigPath: "testdata/cfg20240410.jsonnet",
				},
			},
		},
		"default": {
			c: Default().cfg,
		},
		"bad_asym": {
			c: &cfg20240410{
				Algorithms: Algorithms20231210{
					Asymmetric: "wrong",
				},
			},
			wantErr: ErrUnknownAlgorithmsAsymmetric,
		},
		"bad_pbk": {
			c: &cfg20240410{
				Algorithms: Algorithms20231210{
					Asymmetric: "best",
					PBKDF:      "wrong",
				},
			},
			wantErr: ErrUnknownAlgorithmsPBKDF,
		},
		"bad_sym": {
			c: &cfg20240410{
				Algorithms: Algorithms20231210{
					Asymmetric: "best",
					PBKDF:      "best",
					Symmetric:  "wrong",
				},
			},
			wantErr: ErrUnknownAlgorithmsSymmetric,
		},
		"bad_keyring": {
			c:       badKeyring.cfg,
			wantErr: ErrKeyringNotFound,
		},
		"tamper": {
			c:       c20240410,
			wantErr: ErrTamper,
		},
	}

	for name, tc := range tests { //nolint
		t.Run(name, func(t *testing.T) {
			assert.HasErr(t, (&Config{
				cfg: tc.c,
			}).Parse(ctx, nil), tc.wantErr)
		})
	}

	// Key resolution test
	c := Default()
	ctx = c.lock(ctx)

	c.cfg.KeyPath = "testdata/.rot-keys"

	defer os.Remove("testdata/.rot-keys")

	logger.SetStdin("\n\n")

	pub1, _ := c.NewKeyPathKey(ctx, "key1")

	logger.SetStdin("12345\n12345\n")

	pub2, _ := c.NewKeyPathKey(ctx, "key2")
	f, _ := os.ReadFile("testdata/.rot-keys")
	key, _ := cryptolib.ParseKey[cryptolib.KeyProviderPrivate](strings.Split(string(f), "\n")[0])
	ev, _ := cryptolib.ParseEncryptedValue(strings.Split(string(f), "\n")[1])

	c.cfg.CLI.ConfigPath = "testdata/rot.jsonnet"

	defer os.Remove("testdata/rot.jsonnet")

	c.SetDecryptKey(ctx, "key1", pub1)
	c.SetDecryptKey(ctx, "key2", pub2)
	c.NewKeyring(ctx, "keyring", []string{"key1", "key2"})
	c.SetDecryptKeyPrivateKey(ctx, "key1", "keyring")
	c.SetDecryptKeyPrivateKey(ctx, "key2", "keyring")
	c.clear(ctx)
	c.unlock(ctx)

	k := c.cfg.Keyrings["rot"]
	k.PrivateKey = key
	c.cfg.Keyrings["rot"] = k

	assert.HasErr(t, c.Parse(ctx, nil), nil)
	assert.Equal(t, c.cfg.keys, cryptolib.Keys[cryptolib.KeyProviderPrivate]{key})
	assert.Equal(t, c.cfg.keysEncrypted, cryptolib.EncryptedValues{ev})
	assert.Equal(t, c.cfg.Keyrings["rot"].privateKey, key)
}
