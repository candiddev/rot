package main

import (
	"context"
	"fmt"

	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdRemove(ctx context.Context, args []string, c *cfg) errs.Err {
	n := args[1]

	switch args[0] {
	case "remove-key":
		if _, ok := c.DecryptKeys[n]; !ok {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(fmt.Errorf("decryptKey not found with name: %s", n)))
		}

		delete(c.DecryptKeys, n)
	case "remove-value":
		if _, ok := c.Values[n]; !ok {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(fmt.Errorf("value not found with name: %s", n)))
		}

		delete(c.Values, n)
	}

	return c.save(ctx)
}
