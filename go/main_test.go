package main

import (
	"context"
	"os"
	"testing"
)

var ctx = context.Background()

func TestMain(m *testing.M) {
	os.Setenv("ROT_cli_configPath", "./rot.jsonnet")
	os.Setenv("ROT_keyPath", "./.rot-keys")

	os.Exit(m.Run())
}
