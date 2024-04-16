package main

import (
	"context"
	"fmt"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/logger"
)

type cmdRemoveType string

const (
	cmdRemoveKeyPub  = "Decrypt Key"
	cmdRemoveKeyPrv  = "Decrypt Key Private Key for the current Keyring"
	cmdRemoveKeyring = "Keyring"
	cmdRemoveValue   = "Value"
)

func cmdRemove(t cmdRemoveType) cli.Command[*config.Config] {
	return cli.Command[*config.Config]{
		ArgumentsRequired: []string{
			"name",
		},
		Usage: fmt.Sprintf("Remove a %s from the configuration.", t),
		Run: func(ctx context.Context, args []string, _ cli.Flags, c *config.Config) error {
			var err error

			n := args[1]

			switch args[0] {
			case "remove-keyprv":
				err = c.DeleteDecryptKeyPrivateKey(ctx, c.GetKeyringName(ctx), n)
			case "remove-keypub":
				err = c.DeleteDecryptKey(ctx, n)
			case "remove-keyring":
				err = c.DeleteKeyring(ctx, config.KeyringName(n))
			case "remove-value":
				err = c.DeleteKeyringValue(ctx, c.GetKeyringName(ctx), n)
			}

			return logger.Error(ctx, err)
		},
	}
}
