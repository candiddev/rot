package main

import (
	"context"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
	"github.com/candiddev/shared/go/types"
)

func cmdShowCrt() cli.Command[*cfg] {
	return cli.Command[*cfg]{
		ArgumentsRequired: []string{
			"Certificate value, path, or - for stdin",
		},
		ArgumentsOptional: []string{
			"CA certificate value or path, can be specified multiple times",
		},
		Usage: "Show a certificate, optionally validating the certificate with CA public keys.",
		Run: func(ctx context.Context, args []string, flags cli.Flags, config *cfg) errs.Err {
			cs := args[1]
			if cs == "-" {
				cs = string(cli.ReadStdin())
			}

			f, err := os.ReadFile(cs)
			if err == nil {
				cs = string(f)
			}

			key, err := cryptolib.ParseKey[cryptolib.X509Certificate](cs)
			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			crt, err := key.Key.Certificate()
			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			var e errs.Err

			if len(args) > 2 {
				roots := x509.NewCertPool()

				for i := range args[2:] {
					k := args[i+2]
					f, err := os.ReadFile(k)
					if err == nil {
						k = string(f)
					}

					key, err = cryptolib.ParseKey[cryptolib.X509Certificate](k)
					if err != nil {
						return logger.Error(ctx, errs.ErrReceiver.Wrap(fmt.Errorf("error parsing CA certificate: %w", err)))
					}

					x, err := key.Key.Certificate()
					if err != nil {
						return logger.Error(ctx, errs.ErrReceiver.Wrap(fmt.Errorf("error parsing CA certificate: %w", err)))
					}

					roots.AddCert(x)
				}

				_, err = crt.Verify(x509.VerifyOptions{
					KeyUsages: []x509.ExtKeyUsage{
						x509.ExtKeyUsageAny,
					},
					Roots: roots,
				})

				if err != nil {
					e = errs.ErrReceiver.Wrap(err)
				}
			}

			logger.Raw(types.JSONToString(crt))

			return logger.Error(ctx, e)
		},
	}
}
