package main

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdGenerateCertificate(ctx context.Context, args []string, f cli.Flags, c *cfg) errs.Err {
	pk := args[1]
	if pk == "-" {
		pk = string(cli.ReadStdin())
	}

	if len(strings.Split(pk, ":")) == 1 {
		if err := c.decryptPrivateKey(ctx); err != nil {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}

		v, err := c.decryptValue(ctx, pk)
		if err != nil {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}

		pk = string(v)
	}

	privateKey, err := cryptolib.ParseKey[cryptolib.KeyProviderPrivate](pk)
	if err != nil {
		return logger.Error(ctx, errs.ErrReceiver.Wrap(errors.New("error parsing private key"), err))
	}

	var publicKey cryptolib.Key[cryptolib.KeyProviderPublic]

	if len(args) >= 3 {
		publicKey, _ = cryptolib.ParseKey[cryptolib.KeyProviderPublic](args[2])
	}

	var ca *x509.Certificate

	if len(args) == 4 {
		var key cryptolib.Key[cryptolib.X509Certificate]

		k := args[3]

		f, err := os.ReadFile(k)
		if err == nil {
			k = string(f)
		}

		switch {
		case strings.HasPrefix(k, "----"):
			key, err = cryptolib.PEMToKey[cryptolib.X509Certificate]([]byte(k))
		case len(strings.Split(k, ":")) >= 3:
			key, err = cryptolib.ParseKey[cryptolib.X509Certificate](k)
		default:
			return logger.Error(ctx, errs.ErrReceiver.Wrap(fmt.Errorf("error parsing CA certificate: %w", err)))
		}

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
}
