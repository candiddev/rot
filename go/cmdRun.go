package main

import (
	"context"
	"fmt"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/logger"
)

func cmdRun() cli.Command[*config.Config] {
	return cli.Command[*config.Config]{
		ArgumentsRequired: []string{
			"command",
		},
		Usage: "Run a command and inject Values as environment variables.  Values written to stderr/stdout will be masked with ***.",
		Run: func(ctx context.Context, args []string, _ cli.Flags, c *config.Config) error {
			env := []string{}
			mask := []string{}
			unmask := c.GetUnmask()
			values, err := c.GetKeyringValues(ctx, c.GetKeyringName(ctx))
			if err != nil {
				return logger.Error(ctx, err)
			}

			for k := range values {
				v, err := c.GetValueDecrypted(ctx, c.GetKeyringName(ctx), values[k])
				if err != nil {
					return logger.Error(ctx, err)
				}

				m := true

				for i := range unmask {
					if values[k] == unmask[i] {
						m = false

						break
					}
				}

				if m {
					mask = append(mask, string(v))
				}

				env = append(env, fmt.Sprintf("%s=%s", values[k], v))
			}

			stderr := logger.NewMaskLogger(logger.Stderr, mask)
			stdout := logger.NewMaskLogger(logger.Stdout, mask)

			out, err := c.CLIConfig().Run(ctx, cli.RunOpts{
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
