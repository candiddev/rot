package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdShowKeysValues(keys bool) cli.Command[*cfg] {
	t := "value"
	if keys {
		t = "key"
	}

	return cli.Command[*cfg]{
		Usage: fmt.Sprintf("Show %s names in a configuration.", t),
		Run: func(ctx context.Context, args []string, _ cli.Flags, c *cfg) errs.Err {
			v := []string{}

			if args[0] == "show-keys" {
				for k := range c.DecryptKeys {
					v = append(v, k)
				}
			} else {
				for k := range c.Values {
					v = append(v, k)
				}
			}

			sort.Strings(v)

			return logger.Error(ctx, cli.Print(v))
		},
	}
}
