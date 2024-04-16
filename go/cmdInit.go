package main

import (
	"context"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/cli"
)

func cmdInit() cli.Command[*config.Config] {
	return cli.Command[*config.Config]{
		ArgumentsOptional: []string{
			"keyring name (default: rot)",
			"initial public key, key name or id of an existing key (default: generate a new PBKDF symmetric key)",
		},
		Usage: "Initialize a new Rot configuration.  Will look for a .rot-keys file and use the first available key if none specified as the initial user key.",
		Run: func(ctx context.Context, args []string, _ cli.Flags, c *config.Config) error {
			kr := config.KeyringName("rot")
			if len(args) > 1 {
				kr = config.KeyringName(args[1])
			}

			n := ""
			if len(args) > 2 {
				n = args[2]
			}

			return c.Init(ctx, kr, n)
		},
	}
}
