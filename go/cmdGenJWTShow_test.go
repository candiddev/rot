package main

import (
	"os"
	"strings"
	"testing"

	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/jwt"
)

func TestCmdGenJWTSHow(t *testing.T) {
	cli.RunMain(m, "\n\n", "init")

	defer os.Remove("rot.jsonnet")
	defer os.Remove(".rot-keys")

	cli.RunMain(m, "", "add-pk", "hello")
	pub, _ := cli.RunMain(m, "", "show-value", "-c", "hello")

	// generate-jwt
	j, err := cli.RunMain(m, "hello", "gen-jwt", "-a", "audience", "-e", "60", "-f", "bool=true", "-f", `string="1"`, "-f", "int=1", "-id", "id", "-is", "issuer", "-s", "subject", "-")
	assert.HasErr(t, err, nil)

	out, err := cli.RunMain(m, "", "show-jwt", j, pub)
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, `"audience"`), true)
	assert.Equal(t, strings.Contains(out, `: true`), true)
	assert.Equal(t, strings.Contains(out, `: 1`), true)
	assert.Equal(t, strings.Contains(out, `: "1"`), true)
	assert.Equal(t, strings.Contains(out, `"id"`), true)
	assert.Equal(t, strings.Contains(out, `"issuer"`), true)
	assert.Equal(t, strings.Contains(out, `"subject"`), true)

	out, err = cli.RunMain(m, "", "show-jwt", j)
	assert.HasErr(t, err, jwt.ErrParseNoPublicKeys)
	assert.Equal(t, strings.Contains(out, `"audience"`), true)

	_, jp, _ := cryptolib.NewKeysAsymmetric(cryptolib.AlgorithmBest)

	out, err = cli.RunMain(m, "", "show-jwt", j, jp.String())
	assert.HasErr(t, err, cryptolib.ErrVerify)
	assert.Equal(t, strings.Contains(out, `"audience"`), true)

	j, err = cli.RunMain(m, "", "gen-jwt", "hello")
	assert.HasErr(t, err, nil)

	_, err = cli.RunMain(m, "", "show-jwt", j, pub)
	assert.HasErr(t, err, nil)
}
