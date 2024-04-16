package config

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/logger"
)

func TestConfigDecryptKeysEncrypt(t *testing.T) {
	c := Default()

	prv1, pub1, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)
	prv2, _, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)
	prv3, _, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)

	v1, _ := pub1.Key.EncryptAsymmetric([]byte(prv2.String()), pub1.ID, c.cfg.Algorithms.Asymmetric)
	v2, _ := pub1.Key.EncryptAsymmetric([]byte(prv3.String()), pub1.ID, c.cfg.Algorithms.Asymmetric)

	c.cfg.keys = cryptolib.Keys[cryptolib.KeyProviderPrivate]{
		prv1,
	}
	c.cfg.keysEncrypted = cryptolib.EncryptedValues{
		v1,
		v2,
	}

	c.decryptKeysEncrypted(ctx)

	assert.Equal(t, c.cfg.keys, cryptolib.Keys[cryptolib.KeyProviderPrivate]{
		prv1,
		prv2,
		prv3,
	})
	assert.Equal(t, c.cfg.keysEncrypted, nil)
}

func TestConfigDecryptKeyring(t *testing.T) {
	c := Default()
	ctx = addLock(ctx)

	prv, pub, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)

	k, err := c.decryptKeyring(ctx, "test")

	assert.HasErr(t, err, ErrNotInitialized)
	assert.Equal(t, k, nil)

	c.SetDecryptKey(ctx, "key", pub)
	c.NewKeyring(ctx, "test", []string{"key"})
	c.cfg.keys = cryptolib.Keys[cryptolib.KeyProviderPrivate]{
		prv,
	}

	k, err = c.decryptKeyring(ctx, "test1")

	assert.HasErr(t, err, ErrKeyringNotFound)
	assert.Equal(t, k, nil)

	k, err = c.decryptKeyring(ctx, "test")

	assert.HasErr(t, err, nil)
	assert.Equal(t, k.privateKey.IsNil(), false)
}

func TestConfigGetDecryptKey(t *testing.T) {
	logger.UseTestLogger(t)

	n := time.Now()

	c := Default()
	c.cfg.DecryptKeys["test"] = DecryptKey20240410{
		Modified: n,
	}

	d, err := c.getDecryptKey(ctx, "test1")
	assert.HasErr(t, err, ErrDecryptKeyNotFound)
	assert.Equal(t, d.Modified.IsZero(), true)

	d, err = c.getDecryptKey(ctx, "test")
	assert.HasErr(t, err, nil)
	assert.Equal(t, d.Modified, n)
}

func TestConfigGetKeyring(t *testing.T) {
	logger.UseTestLogger(t)

	prv, _, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)

	c := Default()
	c.cfg.Keyrings["test"] = Keyring20240410{
		PrivateKey: prv,
	}

	k, err := c.getKeyring(ctx, "test1")
	assert.HasErr(t, err, ErrKeyringNotFound)
	assert.Equal(t, k.PrivateKey, cryptolib.Key[cryptolib.KeyProviderPrivate]{})

	k, err = c.getKeyring(ctx, "test")
	assert.HasErr(t, err, nil)
	assert.Equal(t, k.PrivateKey, prv)
}

func TestLocks(t *testing.T) {
	c := Default()
	c.cfg.CLI.ConfigPath = "testdata/test.jsonnet"
	c.cfg.DecryptKeys = map[string]DecryptKey20240410{
		"a": {},
	}

	defer os.Remove("testdata/test.jsonnet")

	ctx = context.Background()
	ctxL := c.lock(ctx)
	ctxL = c.lock(ctxL)
	ctxL = c.lock(ctxL)

	assert.Equal(t, c.mu.TryLock(), false)
	assert.Equal(t, getLocks(ctxL), 3)

	c.unlock(ctxL)
	c.unlock(ctx)

	assert.HasErr(t, c.save(ctx, ErrDecryptKeyNotFound), ErrDecryptKeyNotFound)
	assert.HasErr(t, c.save(ctx, nil), nil)

	f, _ := os.ReadFile("testdata/test.jsonnet")

	assert.Contains(t, string(f), "a: {")
	assert.Equal(t, c.mu.TryLock(), true)

	c.mu.Unlock()
	c.rlock(ctx)

	assert.Equal(t, c.mu.TryRLock(), true)
	assert.Equal(t, c.mu.TryLock(), false)

	c.runlock(ctx)
	c.mu.RUnlock()

	assert.Equal(t, c.mu.TryLock(), true)
}
