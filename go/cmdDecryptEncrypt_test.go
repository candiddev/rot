package main

import (
	"os"
	"testing"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/logger"
)

func TestCmdDecryptEncrypt(t *testing.T) {
	defer os.Remove(".rot-keys")

	c := config.Default()

	logger.SetStdin("\n\n")

	pub, _ := c.NewKeyPathKey(ctx, "test")

	// KDF
	out, err := cli.RunMain(m, "a\nb!a\na\n", []string{"encrypt", "-d", "!"}...)
	assert.HasErr(t, err, nil)

	out, err = cli.RunMain(m, "a\na\n", []string{"decrypt", out}...)
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "a\nb")

	// Key
	out, err = cli.RunMain(m, "a\nb!", []string{"encrypt", "-d", "!", pub.String()}...)
	assert.HasErr(t, err, nil)

	out, err = cli.RunMain(m, "", []string{"decrypt", out}...)
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "a\nb")
}
