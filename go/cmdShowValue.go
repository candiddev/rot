package main

import (
	"context"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/logger"
)

func cmdShowValue() cli.Command[*config.Config] {
	return cli.Command[*config.Config]{
		ArgumentsRequired: []string{
			"name",
		},
		Flags: cli.Flags{
			"c": {
				Usage: "Show the comment only",
			},
			"v": {
				Usage: "Show the value only",
			},
		},
		Usage: "Show a decrypted value.",
		Run: func(ctx context.Context, args []string, f cli.Flags, c *config.Config) error {
			v, err := c.GetValue(ctx, c.GetKeyringName(ctx), args[1])
			if err != nil {
				return logger.Error(ctx, err)
			}

			if _, ok := f.Value("c"); ok {
				logger.Raw(v.Comment + "\n")

				return logger.Error(ctx, nil)
			}

			e, err := c.GetValueDecrypted(ctx, c.GetKeyringName(ctx), args[1])
			if err != nil {
				return logger.Error(ctx, nil)
			}

			if _, ok := f.Value("v"); ok {
				logger.Raw(string(e) + "\n")

				return logger.Error(ctx, nil)
			}

			m := map[string]any{
				"comment":  v.Comment,
				"modified": v.Modified,
				"value":    string(e),
			}

			return logger.Error(ctx, cli.Print(m))
		},
	}
}
