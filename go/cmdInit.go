package main

import (
	"context"
	"os"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdInit() cli.Command[*cfg] { //nolint:gocognit
	return cli.Command[*cfg]{
		ArgumentsOptional: []string{
			"key name or id of an existing key",
			"initial public key, default: generate a PBKDF symmetric key",
		},
		Usage: "Initialize a new Rot configuration.  Will look for a .rot-keys file and use the first available key if none specified as the initial user key.",
		Run: func(ctx context.Context, args []string, f cli.Flags, c *cfg) errs.Err {
			c.DecryptKeys = map[string]cfgDecryptKey{}
			c.Values = map[string]cfgValue{}

			if _, err := os.ReadFile(c.CLI.ConfigPath); err == nil {
				b, err := cli.Prompt(c.CLI.ConfigPath+"%s aleady exists, overwite (yes/no)?", "", false)
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}

				if string(b) != "yes" {
					logger.Raw("Canceling...\n")

					return nil
				}
			}

			c.decryptKeysEncrypted(ctx)

			var id string

			var key cryptolib.KeyProviderPublic

			var err error

			if len(c.keys) > 0 {
				if len(args) > 1 {
					for i := range c.keys {
						if c.keys[i].ID == args[1] {
							id = c.keys[i].ID
							key, err = c.keys[i].Key.Public()

							break
						}
					}
				} else {
					id = c.keys[0].ID
					key, err = c.keys[0].Key.Public()
				}

				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}
			}

			if key != nil {
				args = []string{
					args[0],
					id,
					cryptolib.Key[cryptolib.KeyProviderPublic]{
						ID:  id,
						Key: key,
					}.String(),
				}
			} else if len(args) == 1 {
				args = append(args, "rot")
			}

			prv, pub, err := cryptolib.NewKeysAsymmetric(c.Algorithms.Asymmetric)
			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			c.privateKey = prv
			c.PublicKey = pub

			return logger.Error(ctx, cmdAddKey().Run(ctx, args, f, c))
		},
	}
}
