package main

import (
	"testing"

	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
)

func TestCmdPEM(t *testing.T) {
	_, pub, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)

	// pem
	pubPEM, err := cli.RunMain(m, "", "pem", pub.String())
	assert.HasErr(t, err, nil)

	pemPub, err := cli.RunMain(m, pubPEM, "pem", "-i", pub.ID, "-")
	assert.HasErr(t, err, nil)

	assert.Equal(t, pub.String(), pemPub)
}
