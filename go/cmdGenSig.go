package main

import (
	"context"
	"encoding/base64"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdGenSig() cli.Command[*cfg] {
	return cli.Command[*cfg]{
		ArgumentsRequired: []string{
			"private key value or encrypted value name",
			"data to sign or - for stdin",
		},
		Flags: cli.Flags{
			"s": {
				Usage: "Output just the signature",
			},
		},
		Usage: "Generate a signature and output a standard encoding base64 string.  Must specify the private key of the signer and the data to be signed.  For ECP keys, the hash will be SHA256.  For Ed25519, the hash is unused.",
		Run: func(ctx context.Context, args []string, f cli.Flags, c *cfg) errs.Err {
			data := []byte(args[2])
			if string(data) == "-" {
				data = cli.ReadStdin()
			}

			privateKey, errr := c.decryptValuePrivateKey(ctx, args[1])
			if errr != nil {
				return logger.Error(ctx, errr)
			}

			s, err := cryptolib.NewSignature(privateKey, data)
			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			if _, sig := f.Value("s"); sig {
				logger.Raw(base64.StdEncoding.EncodeToString(s.Signature))
			} else {
				logger.Raw(s.String() + "\n")
			}

			return nil
		},
	}
}
