package main

import (
	"os"
	"testing"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/logger"
)

func TestCmdRemove(t *testing.T) {
	logger.UseTestLogger(t)

	cli.RunMain(m, "\n\n", "init")

	defer os.Remove("rot.jsonnet")
	defer os.Remove(".rot-keys")

	_, pub, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)

	cli.RunMain(m, "", "add-keypub", "remove", pub.String())
	cli.RunMain(m, "", "add-keyprv", "remove", "rot")

	// remove
	_, err := cli.RunMain(m, "", "remove-keyprv", "remove")
	assert.HasErr(t, err, nil)

	_, err = cli.RunMain(m, "", "remove-keypub", "remove")
	assert.HasErr(t, err, nil)

	_, err = cli.RunMain(m, "", "remove-value", "value")
	assert.HasErr(t, err, config.ErrValueNotFound)

	cli.RunMain(m, "1\n1\n", "add-value", "value1")
	cli.RunMain(m, "1\n1\n", "add-value", "value2")

	_, err = cli.RunMain(m, "", "remove-value", "value1")
	assert.HasErr(t, err, nil)

	c := config.Default()
	c.CLIConfig().ConfigPath = "./rot.jsonnet"
	c.Parse(ctx, nil)

	v, _ := c.GetKeyringValues(ctx, c.GetKeyringName(ctx))

	assert.Equal(t, len(c.GetDecryptKeys(ctx)), 1)

	out, _ := c.GetDecryptKeysKeyring(ctx, "rot")

	assert.Equal(t, len(out), 1)
	assert.Equal(t, len(v), 1)
}
