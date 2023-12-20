package main

import (
	"context"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdShowValue(ctx context.Context, args []string, c *cfg) errs.Err {
	if err := c.decryptPrivateKey(ctx); err != nil {
		return logger.Error(ctx, err)
	}

	v, err := c.decryptValue(ctx, args[1])
	if err != nil {
		return logger.Error(ctx, err)
	}

	m := map[string]any{
		"comment":  c.Values[args[1]].Comment,
		"modified": c.Values[args[1]].Modified,
		"value":    string(v),
	}

	return logger.Error(ctx, cli.Print(m))
}
