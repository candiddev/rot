package main

import (
	"strings"
	"testing"

	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/logger"
)

func TestCmdBase64(t *testing.T) {
	tests := map[string][]string{
		"raw std": {"-r"},
		"raw url": {"-r", "-u"},
		"std":     {},
		"url":     {"-u"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			logger.SetStd()

			out, err := cli.RunMain(m, "", append(append([]string{"base64"}, tc...), "input?v")...)

			assert.HasErr(t, err, nil)

			if !strings.Contains(name, "raw") {
				assert.Contains(t, out, "=")
			}

			if strings.Contains(name, "std") {
				assert.Contains(t, out, "/")
			} else {
				assert.Contains(t, out, "_")
			}

			out, err = cli.RunMain(m, out, append(append([]string{"base64", "-d"}, tc...), "-")...)

			assert.HasErr(t, err, nil)
			assert.Equal(t, out, "input?v")
		})
	}
}
