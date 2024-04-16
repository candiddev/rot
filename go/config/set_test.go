package config

import (
	"os"
	"strings"
	"testing"

	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
	"github.com/candiddev/shared/go/types"
)

func TestConfigInit(t *testing.T) {
	logger.UseTestLogger(t)

	c := Default()
	c.cfg.CLI.ConfigPath = "testdata/rot.jsonnet"

	defer os.Remove("testdata/rot.jsonnet")

	c.cfg.KeyPath = "testdata/.rot-keys"

	defer os.Remove("testdata/.rot-keys")

	tests := []struct {
		keyname string
		name    string
		path    string
		stdin   string
		wantErr error
	}{
		{
			keyname: "test",
			name:    "good",
			path:    "testdata/rot.jsonnet",
			stdin:   "\n\n",
		},
		{
			keyname: "test1",
			name:    "no overwrite",
			path:    "testdata/rot.jsonnet",
			stdin:   "no\n",
			wantErr: ErrInitCancelled,
		},
		{
			keyname: "test1",
			name:    "overwrite1",
			path:    "testdata/rot.jsonnet",
			stdin:   "yes\n",
		},
		{
			keyname: "test1",
			name:    "overwrite2",
			path:    "testdata/rot.jsonnet",
			stdin:   "yes\n",
		},
		{
			keyname: "",
			name:    "overwrite3",
			path:    "testdata/rot.jsonnet",
			stdin:   "yes\n",
		},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			logger.SetStdin(tests[i].stdin)

			c.cfg.CLI.ConfigPath = tests[i].path
			err := c.Init(ctx, "keyring", tests[i].keyname)

			assert.HasErr(t, err, tests[i].wantErr)

			if err == nil {
				f, _ := os.ReadFile("testdata/.rot-keys")

				for _, k := range strings.Fields(string(f)) {
					key, _ := cryptolib.ParseKey[cryptolib.KeyProviderPrivate](k)
					c.cfg.keys = append(c.cfg.keys, key)
				}
			}
		})
	}

	f, _ := os.ReadFile("testdata/.rot-keys")

	assert.Equal(t, len(strings.Fields(string(f))), 2)
}

func TestConfigNewKeyring(t *testing.T) {
	logger.UseTestLogger(t)

	c := Default()
	ctx = c.lock(ctx)

	_, pub, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)

	c.SetDecryptKey(ctx, "test1", pub)
	c.SetDecryptKey(ctx, "test2", pub)

	// New keyring, no decryptKeys
	assert.HasErr(t, c.NewKeyring(ctx, "test", nil), ErrNoDecryptKeys)

	// New keyring, with unknown decryptKey
	assert.HasErr(t, c.NewKeyring(ctx, "test", []string{"na"}), ErrDecryptKeyNotFound)

	// New keyring, valid decryptKey
	assert.HasErr(t, c.NewKeyring(ctx, "test", []string{"test1"}), nil)
	assert.Equal(t, len(c.cfg.Keyrings) == 1, true)
	assert.Equal(t, len(c.cfg.DecryptKeys["test1"].PrivateKeys) == 1, true)
	pk1 := c.cfg.Keyrings["test"].PrivateKey
	pk2 := c.cfg.DecryptKeys["test1"].PrivateKeys["test"].PrivateKey
	assert.Equal(t, len(c.cfg.DecryptKeys["test2"].PrivateKeys) == 0, true)

	// Existing keyring, valid decryptKey
	assert.HasErr(t, c.NewKeyring(ctx, "test", []string{"test2"}), nil)
	assert.Equal(t, len(c.cfg.Keyrings) == 1, true)
	assert.Equal(t, len(c.cfg.DecryptKeys["test1"].PrivateKeys) == 1, true)
	assert.Equal(t, len(c.cfg.DecryptKeys["test2"].PrivateKeys) == 1, true)
	assert.Equal(t, c.cfg.Keyrings["test"].PrivateKey, pk1)
	assert.Equal(t, c.cfg.DecryptKeys["test1"].PrivateKeys["test"].PrivateKey, pk2)

	// Too many keyrings
	c.cfg.License.Keyrings = 1
	assert.HasErr(t, c.NewKeyring(ctx, "test1", []string{"test2"}), errs.ErrReceiver)
}

func TestConfigNewKeyPathKey(t *testing.T) {
	logger.UseTestLogger(t)

	c := Default()
	ctx = c.lock(ctx)

	c.cfg.KeyPath = "testdata/.rot-keys"

	defer os.Remove(c.cfg.KeyPath)

	// New key, no encryption
	logger.SetStdin("\n\n")

	pub1, err := c.NewKeyPathKey(ctx, "test1")

	assert.HasErr(t, err, nil)
	assert.Equal(t, pub1.IsNil(), false)

	// New key, with encryption
	logger.SetStdin("123\n123\n")

	pub2, err := c.NewKeyPathKey(ctx, "test2")

	assert.HasErr(t, err, nil)
	assert.Equal(t, pub2.IsNil(), false)
	assert.Equal(t, pub1 != pub2, true)

	f, _ := os.ReadFile(c.cfg.KeyPath)
	keys := strings.Split(string(f), "\n")

	key, _ := cryptolib.ParseKey[cryptolib.KeyProviderPrivate](keys[0])
	kpub, _ := key.Key.Public()
	assert.Equal(t, kpub, pub1.Key)

	logger.SetStdin("123\n123\n")

	v, _ := cryptolib.ParseEncryptedValue(keys[1])
	k, _ := v.Decrypt(nil)
	key, _ = cryptolib.ParseKey[cryptolib.KeyProviderPrivate](string(k))
	kpub, _ = key.Key.Public()
	assert.Equal(t, kpub, pub2.Key)
}

func TestConfigRekey(t *testing.T) {
	logger.UseTestLogger(t)

	c := Default()
	ctx = c.lock(ctx)
	ctx = logger.SetFormat(ctx, logger.FormatKV)

	prv1, pub1, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)
	_, pub2, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)
	_, pub3, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)

	c.SetDecryptKey(ctx, "1", pub1)
	c.SetDecryptKey(ctx, "2", pub2)
	c.SetDecryptKey(ctx, "3", pub3)
	c.NewKeyring(ctx, "k", []string{"1", "2"})
	c.SetValue(ctx, "k", []byte("v1"), "v1", "v1")
	c.SetValue(ctx, "k", []byte("v2"), "v2", "v2")
	c.clear(ctx)
	c.cfg.keys = cryptolib.Keys[cryptolib.KeyProviderPrivate]{
		prv1,
	}

	pk1 := c.cfg.DecryptKeys["1"].PrivateKeys["k"].PrivateKey
	pk2 := c.cfg.DecryptKeys["2"].PrivateKeys["k"].PrivateKey
	pkk := c.cfg.Keyrings["k"].PublicKey

	// No keyring
	assert.HasErr(t, c.Rekey(ctx, "k1"), ErrKeyringNotFound)

	// Good
	assert.HasErr(t, c.Rekey(ctx, "k"), nil)
	assert.Equal(t, c.cfg.DecryptKeys["1"].PrivateKeys["k"].PrivateKey != pk1, true)
	assert.Equal(t, c.cfg.DecryptKeys["2"].PrivateKeys["k"].PrivateKey != pk2, true)
	assert.Equal(t, c.cfg.Keyrings["k"].PublicKey != pkk, true)
	assert.Equal(t, len(c.cfg.DecryptKeys["3"].PrivateKeys), 0)
	assert.HasErr(t, c.cfg.DecryptKeys["1"].PrivateKeys["k"].Signature.Verify([]byte(c.cfg.DecryptKeys["1"].PublicKey.String()), cryptolib.Keys[cryptolib.KeyProviderPublic]{c.cfg.Keyrings["k"].PublicKey}), nil)
	assert.HasErr(t, c.cfg.DecryptKeys["2"].PrivateKeys["k"].Signature.Verify([]byte(c.cfg.DecryptKeys["2"].PublicKey.String()), cryptolib.Keys[cryptolib.KeyProviderPublic]{c.cfg.Keyrings["k"].PublicKey}), nil)

	c.clear(ctx)

	v, _ := c.GetValueDecrypted(ctx, "k", "v1")
	assert.Equal(t, string(v), "v1")

	v, _ = c.GetValueDecrypted(ctx, "k", "v2")
	assert.Equal(t, string(v), "v2")
}

func TestConfigSetDecryptKey(t *testing.T) {
	logger.UseTestLogger(t)

	c := Default()
	ctx = c.lock(ctx)

	_, pub1, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)
	_, pub2, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)

	assert.HasErr(t, c.SetDecryptKey(ctx, "1", pub1), nil)
	assert.Equal(t, c.cfg.DecryptKeys["1"].PublicKey, pub1)

	tm := c.cfg.DecryptKeys["1"].Modified

	assert.HasErr(t, c.SetDecryptKey(ctx, "1", pub2), nil)
	assert.Equal(t, c.cfg.DecryptKeys["1"].PublicKey, pub2)
	assert.Equal(t, c.cfg.DecryptKeys["1"].Modified.Equal(tm), false)
}

func TestConfigSetDecryptKeyPrivateKey(t *testing.T) {
	logger.UseTestLogger(t)

	c := Default()
	ctx = c.lock(ctx)

	prv1, pub1, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)
	prv2, pub2, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)

	c.SetDecryptKey(ctx, "k1", pub1)
	c.SetDecryptKey(ctx, "k2", pub2)
	c.NewKeyring(ctx, "k", []string{"k1"})
	c.clear(ctx)
	c.cfg.keys = cryptolib.Keys[cryptolib.KeyProviderPrivate]{
		prv1,
	}

	// No keyring
	assert.HasErr(t, c.SetDecryptKeyPrivateKey(ctx, "k2", "k2"), ErrKeyringNotFound)

	// No decrypt key
	assert.HasErr(t, c.SetDecryptKeyPrivateKey(ctx, "k", "k3"), ErrDecryptKeyNotFound)

	// Good
	assert.HasErr(t, c.SetDecryptKeyPrivateKey(ctx, "k", "k2"), nil)

	c.clear(ctx)
	c.cfg.keys = cryptolib.Keys[cryptolib.KeyProviderPrivate]{
		prv2,
	}

	_, err := c.decryptKeyring(ctx, "k")

	assert.HasErr(t, err, nil)
}

func TestConfigSetKeyringValue(t *testing.T) {
	logger.UseTestLogger(t)

	c := Default()
	ctx = c.lock(ctx)

	prv, pub, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)

	c.SetDecryptKey(ctx, "k1", pub)
	c.NewKeyring(ctx, "k", []string{"k1"})
	c.clear(ctx)
	c.cfg.keys = cryptolib.Keys[cryptolib.KeyProviderPrivate]{
		prv,
	}

	// Invalid name
	assert.HasErr(t, c.SetValue(ctx, "k", []byte("b"), "!invalid", "comment"), types.ErrEnvAllowedCharacters)

	// Missing keyring
	assert.HasErr(t, c.SetValue(ctx, "k1", []byte("b"), "n", "comment"), ErrKeyringNotFound)

	// Good
	assert.HasErr(t, c.SetValue(ctx, "k", []byte("b"), "n", "comment"), nil)
	assert.Equal(t, c.cfg.Keyrings["k"].Values["n"].Comment, "comment")
	assert.Equal(t, c.cfg.Keyrings["k"].Values["n"].Modified.IsZero(), false)

	o, _ := c.GetValueDecrypted(ctx, "k", "n")

	assert.Equal(t, string(o), "b")
}
