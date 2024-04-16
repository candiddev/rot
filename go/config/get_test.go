package config

import (
	"os"
	"testing"

	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func TestConfigGetAlgorithms(t *testing.T) {
	assert.Equal(t, Default().GetAlgorithms(ctx), Default().cfg.Algorithms)
}

func TestGetDecryptKeys(t *testing.T) {
	c := Default()
	c.cfg.DecryptKeys = map[string]DecryptKey20240410{
		"c": {},
		"b": {},
		"a": {},
	}

	assert.Equal(t, c.GetDecryptKeys(ctx), []string{"a", "b", "c"})
}

func TestConfigGetDecryptKeysKeyring(t *testing.T) {
	logger.UseTestLogger(t)

	c := Default()
	c.cfg.DecryptKeys = map[string]DecryptKey20240410{
		"a": {
			PrivateKeys: map[KeyringName]DecryptKeyPrivateKey20240410{
				"1": {},
				"2": {},
				"3": {},
			},
		},
		"b": {
			PrivateKeys: map[KeyringName]DecryptKeyPrivateKey20240410{
				"2": {},
			},
		},
		"c": {
			PrivateKeys: map[KeyringName]DecryptKeyPrivateKey20240410{
				"1": {},
			},
		},
	}
	c.cfg.Keyrings = map[KeyringName]Keyring20240410{
		"1": {},
		"2": {},
		"3": {},
	}

	out, err := c.GetDecryptKeysKeyring(ctx, "0")
	assert.HasErr(t, err, ErrKeyringNotFound)
	assert.Equal(t, out, nil)

	out, err = c.GetDecryptKeysKeyring(ctx, "1")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, []string{"a", "c"})

	out, err = c.GetDecryptKeysKeyring(ctx, "3")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, []string{"a"})
}

func TestConfigGetKeys(t *testing.T) {
	c := Default()

	prv1, pub1, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)
	prv2, _, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)

	ev, _ := pub1.Key.EncryptAsymmetric([]byte(prv2.String()), pub1.ID, cryptolib.BestEncryptionSymmetric)

	c.cfg.keysEncrypted = cryptolib.EncryptedValues{
		ev,
	}
	c.cfg.keys = cryptolib.Keys[cryptolib.KeyProviderPrivate]{
		prv1,
	}

	assert.Equal(t, c.GetKeys(ctx), cryptolib.Keys[cryptolib.KeyProviderPrivate]{
		prv1,
		prv2,
	})
}

func TestConfigGetKeyrings(t *testing.T) {
	c := Default()
	c.cfg.Keyring = "a"
	ctx = addLock(ctx)

	c.cfg.Keyrings = map[KeyringName]Keyring20240410{
		"a":   {},
		"rot": {},
	}

	assert.Equal(t, c.GetKeyrings(ctx), []string{"a (current)", "rot"})
}

func TestGetKeyringValues(t *testing.T) {
	logger.UseTestLogger(t)

	ctx = addLock(ctx)
	c := Default()
	c.cfg.Keyrings["test"] = Keyring20240410{
		Values: map[string]Value20231210{
			"b": {},
			"a": {},
			"c": {},
		},
	}

	o, err := c.GetKeyringValues(ctx, "test1")
	assert.HasErr(t, err, ErrKeyringNotFound)
	assert.Equal(t, o, nil)

	o, err = c.GetKeyringValues(ctx, "test")
	assert.HasErr(t, err, nil)
	assert.Equal(t, o, []string{"a", "b", "c"})
}

func TestConfigGetPrivatePublicKey(t *testing.T) {
	logger.UseTestLogger(t)

	ctx = addLock(ctx)

	prv1, pub1, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)

	c := Default()
	c.SetDecryptKey(ctx, "key1", pub1)
	c.NewKeyring(ctx, "test", []string{"key1"})
	c.SetValue(ctx, "test", []byte(prv1.String()), "key", pub1.String())

	os.WriteFile("key.prv", []byte(prv1.String()), 0600)
	defer os.Remove("key.prv")

	os.WriteFile("key.pub", []byte(pub1.String()), 0600)
	defer os.Remove("key.pub")

	// Private Key
	tests := map[string]errs.Err{
		prv1.String(): nil,
		"key":         nil,
		"key.prv":     nil,
		pub1.String(): errs.ErrReceiver,
	}

	for in, out := range tests {
		t.Run("private_"+in, func(t *testing.T) {
			k, err := c.GetPrivateKey(ctx, "test", in)
			assert.HasErr(t, err, out)

			if out == nil {
				assert.Equal(t, k, prv1)
			} else {
				assert.Equal(t, k, cryptolib.Key[cryptolib.KeyProviderPrivate]{})
			}
		})
	}

	// Public Key
	tests = map[string]errs.Err{
		prv1.String(): errs.ErrReceiver,
		"key":         nil,
		"key.pub":     nil,
		pub1.String(): nil,
	}

	for in, out := range tests {
		t.Run("public_"+in, func(t *testing.T) {
			k, err := c.GetPublicKey(ctx, "test", in)
			assert.HasErr(t, err, out)

			if out == nil {
				assert.Equal(t, k, pub1)
			} else {
				assert.Equal(t, k, cryptolib.Key[cryptolib.KeyProviderPublic]{})
			}
		})
	}
}

func TestConfigGetValue(t *testing.T) {
	c := Default()
	c.cfg.Keyrings = map[KeyringName]Keyring20240410{
		"test": {
			Values: map[string]Value20231210{
				"test": {
					Comment: "hello",
				},
			},
		},
	}

	v, err := c.GetValue(ctx, "test1", "test")
	assert.HasErr(t, err, ErrKeyringNotFound)
	assert.Equal(t, v, Value20231210{})

	v, err = c.GetValue(ctx, "test", "test1")
	assert.HasErr(t, err, ErrValueNotFound)
	assert.Equal(t, v, Value20231210{})

	v, err = c.GetValue(ctx, "test", "test")
	assert.HasErr(t, err, nil)
	assert.Equal(t, v, Value20231210{
		Comment: "hello",
	})
}

func TestConfigGetvalueDecrypted(t *testing.T) {
	logger.UseTestLogger(t)

	c := Default()
	c.cfg.License.Keyrings = 0
	ctx = addLock(ctx)

	prv1, pub1, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)
	_, pub2, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)

	c.SetDecryptKey(ctx, "key1", pub1)
	c.SetDecryptKey(ctx, "key2", pub2)
	c.NewKeyring(ctx, "test1", []string{"key1"})
	c.NewKeyring(ctx, "test2", []string{"key2"})
	c.SetValue(ctx, "test1", []byte("hello"), "value1", "comment2")
	c.SetValue(ctx, "test2", []byte("world"), "value1", "comment2")

	c.cfg.keys = cryptolib.Keys[cryptolib.KeyProviderPrivate]{
		prv1,
	}

	tests := map[string]struct {
		keyring KeyringName
		value   string
		wantErr error
	}{
		"missing keyring": {
			keyring: "test",
			wantErr: ErrKeyringNotFound,
		},
		"no access": {
			keyring: "test2",
			value:   "value1",
			wantErr: cryptolib.ErrDecryptingKey,
		},
		"missing value": {
			keyring: "test1",
			value:   "value2",
			wantErr: ErrValueNotFound,
		},
		"good": {
			keyring: "test1",
			value:   "value1",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c.clear(ctx)

			out, err := c.GetValueDecrypted(ctx, tc.keyring, tc.value)
			assert.HasErr(t, err, tc.wantErr)

			if tc.wantErr == nil {
				assert.Equal(t, out, []byte("hello"))
			}
		})
	}
}
