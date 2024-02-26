package main

import (
	"context"
	"errors"
	"strconv"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
	"github.com/candiddev/shared/go/types"
)

func cmdGenerateValue(ctx context.Context, args []string, f cli.Flags, c *cfg) errs.Err {
	name := args[1]
	length := 20
	comment := ""

	if len(args) >= 3 {
		l, err := strconv.Atoi(args[2])
		if err != nil {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(errors.New("error parsing length"), err))
		}

		length = l
	}

	if v, ok := f.Value("l"); ok {
		if i, err := strconv.Atoi(v); err == nil {
			length = i
		}
	}

	if len(args) >= 4 {
		comment = args[3]
	}

	v := types.RandString(length)

	if err := c.encryptvalue(ctx, []byte(v), name, comment); err != nil {
		return logger.Error(ctx, err)
	}

	return logger.Error(ctx, c.save(ctx))
}
