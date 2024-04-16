package main

import (
	"context"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/logger"
)

func cmdRekey() cli.Command[*config.Config] {
	return cli.Command[*config.Config]{
		Usage: "Rekey all Keys and Values in a Keyring.",
		Run: func(ctx context.Context, _ []string, _ cli.Flags, c *config.Config) error {
			return logger.Error(ctx, c.Rekey(ctx, c.GetKeyringName(ctx)))
		},
	}
}
