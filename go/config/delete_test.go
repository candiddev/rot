package config

import (
	"testing"

	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func TestConfigDeleteDecryptKey(t *testing.T) {
	logger.UseTestLogger(t)

	c := Default()
	ctx = addLock(ctx)

	c.cfg.DecryptKeys["test"] = DecryptKey20240410{}
	assert.HasErr(t, c.DeleteDecryptKey(ctx, "test1"), errs.ErrReceiver)
	assert.HasErr(t, c.DeleteDecryptKey(ctx, "test"), nil)
	assert.Equal(t, c.cfg.DecryptKeys, map[string]DecryptKey20240410{})
}

func TestConfigDeleteDecryptKeyPrivateKey(t *testing.T) {
	logger.UseTestLogger(t)

	c := Default()
	ctx = addLock(ctx)

	c.cfg.DecryptKeys["test"] = DecryptKey20240410{
		PrivateKeys: map[KeyringName]DecryptKeyPrivateKey20240410{
			"test": {},
		},
	}
	assert.HasErr(t, c.DeleteDecryptKeyPrivateKey(ctx, "test", "test1"), errs.ErrReceiver)
	assert.HasErr(t, c.DeleteDecryptKeyPrivateKey(ctx, "test", "test"), nil)
	assert.Equal(t, c.cfg.DecryptKeys["test"].PrivateKeys, map[KeyringName]DecryptKeyPrivateKey20240410{})
	assert.Equal(t, c.cfg.DecryptKeys["test"].Modified.IsZero(), false)
}

func TestConfigDeleteKeyring(t *testing.T) {
	logger.UseTestLogger(t)

	c := Default()
	ctx = addLock(ctx)

	c.cfg.Keyrings["test"] = Keyring20240410{}

	assert.HasErr(t, c.DeleteKeyring(ctx, "test1"), errs.ErrReceiver)
	assert.HasErr(t, c.DeleteKeyring(ctx, "test"), nil)
	assert.Equal(t, c.cfg.Keyrings, map[KeyringName]Keyring20240410{})
}

func TestConfigDeleteKeyringValue(t *testing.T) {
	logger.UseTestLogger(t)

	c := Default()
	ctx = addLock(ctx)

	c.cfg.Keyrings["test"] = Keyring20240410{
		Values: map[string]Value20231210{
			"test": {},
		},
	}
	assert.HasErr(t, c.DeleteKeyringValue(ctx, "test", "test1"), errs.ErrReceiver)
	assert.HasErr(t, c.DeleteKeyringValue(ctx, "test", "test"), nil)
	assert.Equal(t, c.cfg.Keyrings["test"].Values, map[string]Value20231210{})
}
