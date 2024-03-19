package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
	"github.com/candiddev/shared/go/types"
)

func cmdAddValue() cli.Command[*cfg] {
	return cli.Command[*cfg]{
		ArgumentsRequired: []string{
			"name",
		},
		ArgumentsOptional: []string{
			`comment, default: ""`,
		},
		Flags: cli.Flags{
			"d": {
				Default:     []string{`\n`},
				Placeholder: "delimiter",
				Usage:       "Delimiter",
			},
			"l": {
				Placeholder: "length",
				Usage:       "Generate a random string with this length instead of providing a value",
			},
		},
		Usage: "Add a value to the configuration values.",
		Run: func(ctx context.Context, args []string, f cli.Flags, c *cfg) errs.Err {
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

			var b []byte

			var err error

			if v, ok := f.Value("l"); ok {
				l, err := strconv.Atoi(v)
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(fmt.Errorf("error parsing length: %w", err)))
				}

				b = []byte(types.RandString(l))
			} else {
				b, err = cli.Prompt("Value:", delimiter, true)
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}
			}

			if err := c.encryptvalue(ctx, b, name, comment); err != nil {
				return logger.Error(ctx, err)
			}

			return logger.Error(ctx, c.save(ctx))
		},
	}
}
