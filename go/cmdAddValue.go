package main

import (
	"context"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdAddValue(ctx context.Context, args []string, f cli.Flags, c *cfg) errs.Err {
	if c.PublicKey.IsNil() {
		return logger.Error(ctx, errNotInitialized)
	}

	name := args[1]
	delimiter := ""
	comment := ""

	if len(args) >= 3 {
		comment = args[2]
	}

	if v, ok := f.Value("d"); ok {
		delimiter = v
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
