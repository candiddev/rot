package main

import (
	"os"
	"strings"
	"testing"

	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/logger"
)

func TestCmdRun(t *testing.T) {
	logger.UseTestLogger(t)

	cli.RunMain(m, "\n\n", "init")
	cli.RunMain(m, "hello world!\nhello world!\n", "add-value", "test")
	cli.RunMain(m, "avalue\navalue\n", "add-value", "value")

	defer os.Remove("rot.jsonnet")
	defer os.Remove(".rot-keys")

	out, err := cli.RunMain(m, "123\n123\n", "run", "env")
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, "test=***"), true)
	assert.Equal(t, strings.Contains(out, "value=***"), true)

	out, err = cli.RunMain(m, "123\n123\n", "-x", `unmask=["test"]`, "run", "env")
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, "test=hello world!"), true)
	assert.Equal(t, strings.Contains(out, "value=***"), true)
}
