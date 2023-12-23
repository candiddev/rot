package main

import (
	"context"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdAddValue(ctx context.Context, args []string, c *cfg) errs.Err {
	name := args[1]
	delimiter := ""
	comment := ""

	if len(args) >= 3 {
		comment = args[2]
	}

	if len(args) >= 4 {
		delimiter = args[3]
	}

	v, err := cli.Prompt("Value:", delimiter, true)
	if err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
	}

	if err := c.encryptvalue(ctx, v, name, comment); err != nil {
		return logger.Error(ctx, err)
	}

	return logger.Error(ctx, c.save(ctx))
}
