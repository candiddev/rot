package main

import (
	"context"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdShowPrivateKey(ctx context.Context, _ []string, _ cli.Flags, c *cfg) errs.Err {
	if err := c.decryptPrivateKey(ctx); err != nil {
		return logger.Error(ctx, err)
	}

	logger.Raw(c.privateKey.String())

	return logger.Error(ctx, nil)
}
