package main

import (
	"os"
	"testing"

	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/cli"
)

func TestCmdGenSSH(t *testing.T) {
	cli.RunMain(m, "\n\n", "init")

	defer os.Remove("rot.jsonnet")
	defer os.Remove(".rot-keys")

	cli.RunMain(m, "", "add-pk", "hello")
	pub, _ := cli.RunMain(m, "", "show-value", "-c", "hello")

	_, err := cli.RunMain(m, "", "gen-ssh", "-c", "source-address=localhost", "-e", "permit-pty", "-h", "-i", "123", "-p", "root", "-v", "360", "hello", pub)
	assert.HasErr(t, err, nil)
	_, err = cli.RunMain(m, "", "ssh", pub)
	assert.HasErr(t, err, nil)
}
