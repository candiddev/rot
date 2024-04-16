package main

import (
	"context"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/logger"
)

func cmdAddKeyring() cli.Command[*config.Config] {
	return cli.Command[*config.Config]{
		ArgumentsRequired: []string{
			"keyring name",
			"decrypt key names",
		},
		Usage: "Add a new or modify an existing Keyring and give the provided Decrypt Keys access to it.",
		Run: func(ctx context.Context, args []string, _ cli.Flags, c *config.Config) error {
			return logger.Error(ctx, c.NewKeyring(ctx, config.KeyringName(args[1]), args[2:]))
		},
	}
}
