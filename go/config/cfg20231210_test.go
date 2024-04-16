package config

import (
	"testing"

	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/jsonnet"
	"github.com/candiddev/shared/go/logger"
	"github.com/candiddev/shared/go/types"
)

func get20231210() *cfg20231210 {
	c := default20231210()
	r := jsonnet.NewRender(ctx, nil)
	i, _ := r.GetPath(ctx, "testdata/cfg20231210.jsonnet")
	r.Import(i)
	r.Render(ctx, c)

	return c
}

func TestCfg20231210Upgrade(t *testing.T) {
	logger.UseTestLogger(t)

	c := get20231210()
	prv, _, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)

	c.CLI.ConfigPath = "configpath"
	c.Keys = types.SliceString{prv.String()}
	c.KeyPath = "keypath"

	assert.EqualJSON(t, c.upgrade(ctx), &cfg20240410{
		Algorithms: c.Algorithms,
		CLI:        c.CLI,
		DecryptKeys: map[string]DecryptKey20240410{
			"test1": {
				Modified: c.DecryptKeys["test1"].Modified,
				PrivateKeys: map[KeyringName]DecryptKeyPrivateKey20240410{
					"rot": {
						PrivateKey: c.DecryptKeys["test1"].PrivateKey,
						Signature:  c.DecryptKeys["test1"].Signature,
					},
				},
				PublicKey: c.DecryptKeys["test1"].PublicKey,
			},
			"test2": {
				Modified: c.DecryptKeys["test2"].Modified,
				PrivateKeys: map[KeyringName]DecryptKeyPrivateKey20240410{
					"rot": {
						PrivateKey: c.DecryptKeys["test2"].PrivateKey,
						Signature:  c.DecryptKeys["test2"].Signature,
					},
				},
				PublicKey: c.DecryptKeys["test2"].PublicKey,
			},
		},
		KeyPath: c.KeyPath,
		Keys:    c.Keys,
		Keyring: "rot",
		Keyrings: map[KeyringName]Keyring20240410{
			"rot": {
				PrivateKey: c.PrivateKey,
				PublicKey:  c.PublicKey,
				Values:     c.Values,
			},
		},
		License: License20240410{
			Keyrings: 2,
		},
		Version: cfgVersion20240410,
	})
}
