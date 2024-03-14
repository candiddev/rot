package main

import (
	"context"
	"crypto/x509"
	"os"
	"strconv"
	"strings"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdGenCrt() cli.Command[*cfg] {
	return cli.Command[*cfg]{
		ArgumentsRequired: []string{
			"private key value, encrypted value name, or - for stdin",
		},
		ArgumentsOptional: []string{
			"public key value, encrypted value name, or path",
			"ca certificate or path",
		},
		Flags: cli.Flags{
			"c": {
				Usage: "Create a CA certificate",
			},
			"d": {
				Placeholder: "hostname",
				Usage:       "DNS hostname (can be used multiple times)",
			},
			"e": {
				Default:     []string{"31536000"},
				Placeholder: "seconds",
				Usage:       "Expiration in seconds",
			},
			"eu": {
				Default:     []string{"clientAuth", "serverAuth"},
				Placeholder: "extended key usage",
				Usage:       "Extended key usage, valid values: " + strings.Join(cryptolib.ValidX509ExtKeyUsages(), ", "),
			},
			"i": {
				Placeholder: "address",
				Usage:       "IP address (can be used multiple times)",
			},
			"ku": {
				Default:     []string{"digitalSignature"},
				Placeholder: "key usage",
				Usage:       "Key usage, valid values: " + strings.Join(cryptolib.ValidX509KeyUsages(), ", "),
			},
			"n": {
				Placeholder: "name",
				Usage:       "Common Name (CN)",
			},
		},
		Usage: "Generate an X.509 certificate and output a PEM-formatted certificate to stdout.  Must specify the private key of the signer (for CA signed certificates) or the private key of the certificate (for self-signed certificates).  A public key can be specified, otherwise the public key of the private key will be used.",
		Run: func(ctx context.Context, args []string, f cli.Flags, c *cfg) errs.Err {
			pk := args[1]
			if pk == "-" {
				pk = string(cli.ReadStdin())
			}

			privateKey, errr := c.decryptValuePrivateKey(ctx, pk)
			if errr != nil {
				return logger.Error(ctx, errr)
			}

			var publicKey cryptolib.Key[cryptolib.KeyProviderPublic]

			if len(args) >= 3 {
				var err error

				publicKey, err = c.publicKey(args[2])
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}
			}

			var ca *x509.Certificate

			if len(args) == 4 {
				var key cryptolib.Key[cryptolib.X509Certificate]

				k := args[3]

				f, err := os.ReadFile(k)
				if err == nil {
					k = string(f)
				}

				key, err = cryptolib.ParseKey[cryptolib.X509Certificate](k)
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}

				ca, err = key.Key.Certificate()
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
				}
			}

			_, isCA := f.Value("c")

			dns, _ := f.Values("d")

			var expires int

			e, _ := f.Value("e")
			if i, err := strconv.Atoi(e); err == nil {
				expires = i
			}

			eu, _ := f.Values("eu")
			ku, _ := f.Values("ku")
			ips, _ := f.Values("i")

			cn, _ := f.Value("n")

			crt, err := cryptolib.NewX509Certificate(privateKey, publicKey, cn, cryptolib.NewX509CertificateOpts{
				CACertificate: ca,
				DNSNames:      dns,
				ExtKeyUsages:  eu,
				IPAddresses:   ips,
				IsCA:          isCA,
				KeyUsages:     ku,
				NotAfterSec:   expires,
			})

			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			logger.Raw(string(cryptolib.KeyToPEM(crt)))

			return nil
		},
	}
}
