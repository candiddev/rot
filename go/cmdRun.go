package main

import (
	"context"
	"fmt"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdRun(ctx context.Context, args []string, c *cfg) errs.Err {
	err := c.decryptPrivateKey(ctx)
	if err != nil {
		return logger.Error(ctx, err)
	}

	env := []string{}

	for k := range c.Values {
		v, err := c.decryptValue(ctx, k)
		if err != nil {
			return logger.Error(ctx, err)
		}

		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	out, err := c.CLI.Run(ctx, cli.RunOpts{
		Args:               args[2:],
		Command:            args[1],
		Environment:        env,
		EnvironmentInherit: true,
		StreamLogs:         true,
	})

	logger.Raw(out.String() + "\n")

	return logger.Error(ctx, err)
}
