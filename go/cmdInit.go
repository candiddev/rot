package main

import (
	"context"
	"fmt"
	"os"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdInit(ctx context.Context, args []string, c *cfg) errs.Err {
	c.DecryptKeys = map[string]cfgDecryptKey{}
	c.Values = map[string]cfgValue{}

	if _, err := os.ReadFile(c.CLI.ConfigPath); err == nil {
		b, err := cli.Prompt(fmt.Sprintf("%s aleady exists, overwite (yes/no)?", c.CLI.ConfigPath), "", false)
		if err != nil {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}

		if string(b[0]) != "yes" {
			logger.Raw("Canceling...\n")

			return nil
		}
	}

	prv, pub, err := cryptolib.NewKeysAsymmetric(c.Algorithms.Asymmetric)
	if err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
	}

	c.privateKey = prv
	c.PublicKey = pub

	return logger.Error(ctx, cmdAddKey(ctx, args, c))
}
