package main

import (
	"context"
	"fmt"

	"github.com/candiddev/rot/go/config"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/logger"
)

type cmdShowKeysValuesType string

const (
	cmdShowVluesTypeKeys     cmdShowKeysValuesType = "Keys"
	cmdShowVluesTypeKeyrings cmdShowKeysValuesType = "Keyrings"
	cmdShowVluesTypeValues   cmdShowKeysValuesType = "Values"
)

func cmdShowKeysValues(t cmdShowKeysValuesType) cli.Command[*config.Config] {
	f := cli.Flags{}
	if t == cmdShowVluesTypeKeys {
		f["a"] = &cli.Flag{
			Usage: fmt.Sprintf("Show all %s, not just in the current keyring", t),
		}
	}

	return cli.Command[*config.Config]{
		Flags: f,
		Usage: fmt.Sprintf("Show %s in a configuration.", t),
		Run: func(ctx context.Context, _ []string, f cli.Flags, c *config.Config) error {
			var err error

			var s []string

			_, all := f.Value("a")

			switch t {
			case cmdShowVluesTypeKeys:
				if all {
					s = c.GetDecryptKeys(ctx)
				} else {
					s, err = c.GetDecryptKeysKeyring(ctx, c.GetKeyringName(ctx))
				}
			case cmdShowVluesTypeKeyrings:
				s = c.GetKeyrings(ctx)
			case cmdShowVluesTypeValues:
				s, err = c.GetKeyringValues(ctx, c.GetKeyringName(ctx))
			}

			if err != nil {
				return logger.Error(ctx, err)
			}

			return logger.Error(ctx, cli.Print(s))
		},
	}
}
