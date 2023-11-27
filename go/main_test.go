package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
)

func TestM(t *testing.T) {
	c := defaultCfg()
	c.CLI.ConfigPath = "rot.jsonnet"
	ctx := context.Background()

	t.Setenv("ROT_cli_logFormat", "kv")
	t.Setenv("ROT_cli_noColor", "true")

	// init
	out, err := cli.RunMain(m, "\n\n", "add-key", "test1")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "")

	// check config
	assert.HasErr(t, c.Parse(ctx, nil), nil)
	assert.Equal(t, len(c.DecryptKeys), 1)

	// add-key
	out, err = cli.RunMain(m, "123\n123\n", "add-key", "test2")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "")

	// check config
	c.Parse(ctx, nil)
	assert.Equal(t, len(c.DecryptKeys), 2)

	// check keys
	f, _ := os.ReadFile(".rot-keys")
	assert.Equal(t, len(strings.Split(string(f), "\n")), 3)

	// add-value
	out, err = cli.RunMain(m, "hello world!", "add-value", "test", "t", "t")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "")

	// algorithms
	out, err = cli.RunMain(m, "", "show-algorithms")
	assert.HasErr(t, err, nil)
	assert.Equal(t, len(strings.Split(out, "\n")), 19)

	// generate-key
	out, err = cli.RunMain(m, "\n\n", "generate-key", "encrypt-asymmetric")
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, string(cryptolib.AlgorithmEd25519Private)), true)

	keys := map[string]any{}
	json.Unmarshal([]byte(out), &keys)

	// generate-value
	out, err = cli.RunMain(m, "", "generate-value", "value", "20", "vc")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "")

	// check config
	c.Parse(ctx, nil)
	assert.Equal(t, len(c.Values), 2)
	assert.Equal(t, c.Values["test"].Comment, "t")
	assert.Equal(t, c.Values["value"].Comment, "vc")

	// encrypt
	out, err = cli.RunMain(m, "secret", "encrypt", c.DecryptKeys["test1"].PublicKey.String())
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, string(cryptolib.BestEncryptionAsymmetric)), true)

	// decrypt
	out, err = cli.RunMain(m, "123\n123\n", "decrypt", out)
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "secret")

	out, err = cli.RunMain(m, "", "-x", "keyPath=test", "-x", fmt.Sprintf(`keys=["%s"]`, keys["privateKey"]), "decrypt", out)
	assert.HasErr(t, err, cryptolib.ErrUnknownEncryption)
	assert.Equal(t, out != "secret", true)

	// show-value
	out, err = cli.RunMain(m, "123\n123\n", "show-value", "test")
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, `"value": "hello world!"`), true)

	out, err = cli.RunMain(m, "", "-x", "keyPath=test", "-x", fmt.Sprintf(`keys=["%s"]`, keys["privateKey"]), "show-value", "test")
	assert.HasErr(t, err, errs.ErrReceiver)
	assert.Equal(t, strings.Contains(out, `"value": "hello world!"`), false)

	out, err = cli.RunMain(m, "123\n123\n", "show-value", "value")
	assert.HasErr(t, err, nil)

	json.Unmarshal([]byte(out), &keys)

	assert.Equal(t, len(keys["value"].(string)), 20)

	// rekey
	t.Setenv("ROT_algorithms_asymmetric", string(cryptolib.KDFECDHP256))
	t.Setenv("ROT_algorithms_symmetric", string(cryptolib.EncryptionAES128GCM))

	out, err = cli.RunMain(m, "123\n123\n", "rekey")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "")

	n := defaultCfg()
	n.CLI = c.CLI
	n.Parse(ctx, nil)

	assert.Equal(t, c.PublicKey != n.PublicKey, true)
	assert.Equal(t, n.PublicKey.Key.Algorithm(), cryptolib.AlgorithmECP256Public)
	assert.Equal(t, n.DecryptKeys["test1"].PrivateKey != c.DecryptKeys["test1"].PrivateKey, true)
	assert.Equal(t, n.DecryptKeys["test1"].PrivateKey.KDF, c.DecryptKeys["test1"].PrivateKey.KDF)
	assert.Equal(t, n.Values["value"].Key.Ciphertext != c.Values["value"].Key.Ciphertext, true)
	assert.Equal(t, n.Values["value"].Key.Encryption, cryptolib.EncryptionAES128GCM)
	assert.Equal(t, n.Values["value"].Key.KDF, cryptolib.KDFECDHP256)

	// run
	out, err = cli.RunMain(m, "123\n123\n", "run", "env")
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, "test=hello world!"), true)

	// tamper
	k := n.DecryptKeys["test1"]
	k.PublicKey.ID = "new"
	n.DecryptKeys["test1"] = k
	delete(n.DecryptKeys, "test2")
	n.save(ctx)

	out, err = cli.RunMain(m, "123\n123\n", "show-value", "test")
	assert.HasErr(t, err, errs.ErrReceiver)
	assert.Equal(t, strings.Contains(out, "tampering"), true)

	os.RemoveAll("rot.jsonnet")
	os.RemoveAll(".rot-keys")
}
