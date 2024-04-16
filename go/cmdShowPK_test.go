package main

import (
	"os"
	"testing"

	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/logger"
)

func TestCmdShowPK(t *testing.T) {
	logger.UseTestLogger(t)

	cli.RunMain(m, "\n\n", "init", "rot", "key")

	defer os.Remove("rot.jsonnet")
	defer os.Remove(".rot-keys")

	// show-public-key
	out, err := cli.RunMain(m, "", "show-pk", "key")
	assert.HasErr(t, err, nil)
	_, err = cryptolib.ParseKey[cryptolib.KeyProviderPublic](out)
	assert.HasErr(t, err, nil)

	prv1, pub1, _ := cryptolib.NewKeysAsymmetric(cryptolib.AlgorithmBest)

	// show-public-key
	pub2, err := cli.RunMain(m, "", "show-pk", prv1.String())
	assert.HasErr(t, err, nil)
	assert.Equal(t, pub1.String(), pub2)

	pub2, err = cli.RunMain(m, prv1.String(), "show-pk", "-")
	assert.HasErr(t, err, nil)
	assert.Equal(t, pub1.String(), pub2)
}
