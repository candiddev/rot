package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
	"github.com/candiddev/shared/go/types"
)

func cmdAddValue() cli.Command[*config.Config] {
	return cli.Command[*config.Config]{
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
		Usage: "Add a Value to a Keyring.",
		Run: func(ctx context.Context, args []string, f cli.Flags, c *config.Config) error {
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
				b, err = logger.Prompt("Value:", delimiter, true)
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}
			}

			return logger.Error(ctx, c.SetValue(ctx, c.GetKeyringName(ctx), b, name, comment))
		},
	}
}
