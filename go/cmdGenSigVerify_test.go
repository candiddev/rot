package main

import (
	"os"
	"testing"

	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/errs"
)

func TestCmdGenSigVerify(t *testing.T) {
	cli.RunMain(m, "\n\n", "init")

	defer os.Remove("rot.jsonnet")
	defer os.Remove(".rot-keys")

	cli.RunMain(m, "", "add-pk", "hello")
	cli.RunMain(m, "", "add-pk", "goodbye")

	sig, err := cli.RunMain(m, "helloworld", "gen-sig", "hello", "-")
	assert.HasErr(t, err, nil)
	_, err = cli.RunMain(m, "helloworld", "verify-sig", "hello", "-", sig)
	assert.HasErr(t, err, nil)
	_, err = cli.RunMain(m, "helloworld", "verify-sig", "goodbye", "-", sig)
	assert.HasErr(t, err, errs.ErrReceiver)
}
