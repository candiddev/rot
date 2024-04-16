package main

import (
	"crypto/x509"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
)

func TestCmdGenCrtShow(t *testing.T) {
	cli.RunMain(m, "\n\n", "init")

	defer os.Remove("rot.jsonnet")
	defer os.Remove(".rot-keys")

	cli.RunMain(m, "", "add-pk", "hello")
	prv, _ := cli.RunMain(m, "", "show-value", "-v", "hello")

	// generate-certificate
	exp := 60 * 60
	crtPEM, err := cli.RunMain(m, "hello", "gen-crt", "-c", "-d", "localhost", "-e", strconv.Itoa(exp), "-eu", "ocspSigning", "-eu", "clientAuth", "-ku", "digitalSignature", "-ku", "keyAgreement", "-i", "127.0.0.1", "-n", "My CA", "-")

	assert.HasErr(t, err, nil)

	out, err := cli.RunMain(m, crtPEM, "show-crt", "-")
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, "localhost"), true)

	os.WriteFile("ca.pem", []byte(crtPEM), 0600)

	defer os.Remove("ca.pem")

	cs, err := cli.RunMain(m, crtPEM, "pem", "-")

	assert.HasErr(t, err, nil)

	crt, e := cryptolib.ParseKey[cryptolib.X509Certificate](cs)

	assert.HasErr(t, e, nil)

	x, e := crt.Key.Certificate()

	assert.HasErr(t, e, nil)
	assert.Equal(t, x.IsCA, true)
	assert.Equal(t, x.DNSNames, []string{"localhost"})

	assert.Equal(t, x.ExtKeyUsage, []x509.ExtKeyUsage{
		x509.ExtKeyUsageOCSPSigning,
		x509.ExtKeyUsageClientAuth,
	})

	assert.Equal(t, x.KeyUsage, x509.KeyUsageCRLSign|x509.KeyUsageCertSign|x509.KeyUsageDigitalSignature|x509.KeyUsageKeyAgreement)
	assert.Equal(t, x.NotBefore.Before(time.Now().Add(60*60*24*time.Second)), true)
	assert.Equal(t, x.IPAddresses, []net.IP{net.ParseIP("127.0.0.1")})
	assert.Equal(t, x.Subject.CommonName, "My CA")

	_, pub, _ := cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)

	crtPEM, err = cli.RunMain(m, "", "gen-crt", prv, pub.String(), "ca.pem")

	assert.HasErr(t, err, nil)

	cs, err = cli.RunMain(m, crtPEM, "pem", "-")

	assert.HasErr(t, err, nil)

	crt, e = cryptolib.ParseKey[cryptolib.X509Certificate](cs)

	assert.HasErr(t, e, nil)

	x, e = crt.Key.Certificate()

	assert.HasErr(t, e, nil)

	assert.Equal(t, x.ExtKeyUsage, []x509.ExtKeyUsage{
		x509.ExtKeyUsageClientAuth,
		x509.ExtKeyUsageServerAuth,
	})

	assert.Equal(t, x.IsCA, false)
	assert.Equal(t, x.Issuer.CommonName, "My CA")
	assert.Equal(t, x.KeyUsage, x509.KeyUsageDigitalSignature)

	os.WriteFile("crt.pem", []byte(crtPEM), 0600)

	defer os.Remove("crt.pem")

	out, err = cli.RunMain(m, crtPEM, "show-crt", "crt.pem", "ca.pem")
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, "My CA"), true)

	crtPEM, _ = cli.RunMain(m, "", "gen-crt", prv, pub.String(), "ca.pem")
	os.WriteFile("ca.pem", []byte(crtPEM), 0600)

	out, err = cli.RunMain(m, crtPEM, "show-crt", "crt.pem", "ca.pem")
	assert.HasErr(t, err, errs.ErrReceiver)
	assert.Equal(t, strings.Contains(out, "My CA"), true)
}
