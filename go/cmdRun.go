package main

import (
	"context"
	"fmt"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdRun() cli.Command[*cfg] {
	return cli.Command[*cfg]{
		ArgumentsRequired: []string{
			"command",
		},
		Usage: "Run a command and inject configuration values as environment variables.  Values written to stderr/stdout will be masked with ***.",
		Run: func(ctx context.Context, args []string, _ cli.Flags, c *cfg) errs.Err {
			err := c.decryptPrivateKey(ctx)
			if err != nil {
				return logger.Error(ctx, err)
			}

			env := []string{}
			mask := []string{}

			for k := range c.Values {
				v, err := c.decryptValue(ctx, k)
				if err != nil {
					return logger.Error(ctx, err)
				}

				m := true

				for i := range c.Unmask {
					if k == c.Unmask[i] {
						m = false

						break
					}
				}

				if m {
					mask = append(mask, string(v))
				}

				env = append(env, fmt.Sprintf("%s=%s", k, v))
			}

			stderr := logger.NewMaskLogger(logger.Stderr, mask)
			stdout := logger.NewMaskLogger(logger.Stdout, mask)

			out, err := c.CLI.Run(ctx, cli.RunOpts{
				Args:               args[2:],
				Command:            args[1],
				Environment:        env,
				EnvironmentInherit: true,
				Stderr:             stderr,
				Stdout:             stdout,
			})

			o := out.String()
			if o != "" {
				logger.Raw(out.String() + "\n")
			}

			return logger.Error(ctx, err)
		},
	}
}
