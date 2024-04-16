package main

import (
	"os"
	"strings"
	"testing"

	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/logger"
)

func TestCmdShowValue(t *testing.T) {
	logger.UseTestLogger(t)

	cli.RunMain(m, "\n\n", "init", "rot", "key")

	defer os.Remove("rot.jsonnet")
	defer os.Remove(".rot-keys")

	cli.RunMain(m, "123\n123\n", "add-value", "test", "comment")
	cli.RunMain(m, "\n\n", "add-keypub", "test1")
	cli.RunMain(m, "\n\n", "add-keypub", "test2")
	cli.RunMain(m, "", "add-keyprv", "test1")

	// show-keys
	out, err := cli.RunMain(m, "", "show-keys")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, `[
  "key",
  "test1"
]`)

	out, err = cli.RunMain(m, "", "show-keys", "-a")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, `[
  "key",
  "test1",
  "test2"
]`)

	// show-value
	out, err = cli.RunMain(m, "", "show-value", "test")
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, `"comment": "comment"`), true)
	assert.Equal(t, strings.Contains(out, `"value": "123"`), true)

	out, _ = cli.RunMain(m, "123\n123\n", "show-value", "-c", "test")
	assert.Equal(t, out, "comment")

	out, _ = cli.RunMain(m, "123\n123\n", "show-value", "-v", "test")
	assert.Equal(t, out, "123")
}
