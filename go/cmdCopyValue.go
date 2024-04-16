package main

import (
	"context"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/logger"
)

func cmdCopyValue() cli.Command[*config.Config] {
	return cli.Command[*config.Config]{
		ArgumentsRequired: []string{
			"source keyring",
			"source value",
		},
		ArgumentsOptional: []string{
			"destination name",
			"destination comment",
		},
		Usage: "Copy a Value from a Keyring to an existing Keyring, optionally providing a new name and comment for it.",
		Run: func(ctx context.Context, args []string, _ cli.Flags, c *config.Config) error {
			srck := config.KeyringName(args[1])
			srcv := args[2]

			dstv := srcv
			if len(args) > 3 {
				dstv = args[3]
			}

			v, err := c.GetValue(ctx, srck, srcv)
			if err != nil {
				return logger.Error(ctx, err)
			}

			k, err := c.GetValueDecrypted(ctx, srck, srcv)
			if err != nil {
				return logger.Error(ctx, err)
			}

			dstc := v.Comment
			if len(args) > 4 {
				dstc = args[4]
			}

			return logger.Error(ctx, c.SetValue(ctx, c.GetKeyringName(ctx), k, dstv, dstc))
		},
	}
}
