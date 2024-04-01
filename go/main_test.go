package main

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
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
	"github.com/candiddev/shared/go/jwt"
)

func TestM(t *testing.T) {
	c := defaultCfg()
	c.CLI.ConfigPath = "./rot.jsonnet"
	ctx := context.Background()

	t.Setenv("ROT_cli_logFormat", "kv")
	t.Setenv("ROT_cli_noColor", "true")
	t.Setenv("ROT_keyPath", "./.rot-keys")

	// init
	out, err := cli.RunMain(m, "", "init")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "")

	os.Remove("./.rot-keys")
	os.Remove("./rot.jsonnet")

	out, err = cli.RunMain(m, "\n\n", "init", "test1")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "")

	os.Remove("./rot.jsonnet")

	out, err = cli.RunMain(m, "\n\n", "init", "test2")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "")

	os.Remove("./rot.jsonnet")

	out, err = cli.RunMain(m, "", "init")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "")

	// check config
	assert.HasErr(t, c.Parse(ctx, nil), nil)
	assert.Equal(t, len(c.DecryptKeys), 1)

	// show-public-key
	out, err = cli.RunMain(m, "", "show-pk", "test1")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, c.DecryptKeys["test1"].PublicKey.String())

	_, pub, _ := cryptolib.NewKeysAsymmetric(cryptolib.AlgorithmBest)

	// add-key
	out, err = cli.RunMain(m, "", "add-key", "test2", pub.String())
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "")

	out, err = cli.RunMain(m, "123\n123\n", "add-key", "test3")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "")

	// check config
	c.Parse(ctx, nil)
	assert.Equal(t, len(c.DecryptKeys), 3)

	// check keys
	f, _ := os.ReadFile(".rot-keys")
	assert.Equal(t, len(strings.Split(string(f), "\n")), 4)

	// add-value
	out, err = cli.RunMain(m, "hello world!", "add-value", "test", "t", "t")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "")

	out, err = cli.RunMain(m, "", "add-value", "-l", "20", "value", "vc")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "")

	// algorithms
	out, err = cli.RunMain(m, "", "show-alg")
	assert.HasErr(t, err, nil)
	assert.Equal(t, len(strings.Split(out, "\n")), 20)

	// gen-keys
	out, err = cli.RunMain(m, "\n\n", "gen-keys", "encrypt-asymmetric")
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, string(cryptolib.AlgorithmEd25519Private)), true)

	keys := map[string]any{}
	json.Unmarshal([]byte(out), &keys)
	pk := keys["publicKey"].(string) //nolint:revive

	// check config
	c.Parse(ctx, nil)
	assert.Equal(t, len(c.Values), 2)
	assert.Equal(t, c.Values["test"].Comment, "t")
	assert.Equal(t, c.Values["value"].Comment, "vc")

	// encrypt - asymmetric
	out, err = cli.RunMain(m, "secret", "encrypt", c.DecryptKeys["test1"].PublicKey.String())
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, string(cryptolib.BestEncryptionAsymmetric)), true)

	// decrypt - asymmetric
	out, err = cli.RunMain(m, "123\n123\n", "decrypt", out)
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "secret")

	// encrypt - symmetric
	out, err = cli.RunMain(m, "secret\n123\n123\n", "encrypt")
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, string(cryptolib.BestEncryptionSymmetric)), true)

	// decrypt - symmetric
	out, err = cli.RunMain(m, "123\n123\n", "decrypt", out)
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "secret")

	out, err = cli.RunMain(m, "", "-x", "keyPath=test", "-x", fmt.Sprintf(`keys=["%s"]`, keys["privateKey"]), "decrypt", out)
	assert.HasErr(t, err, cryptolib.ErrUnknownEncryption)
	assert.Equal(t, out != "secret", true)

	// show-value
	out, err = cli.RunMain(m, "123\n123\n", "show-value", "test")
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, `"value": "hello world!"`), true)

	out, err = cli.RunMain(m, "", "-x", "keyPath=test", "-x", fmt.Sprintf(`keys=["%s"]`, keys["privateKey"]), "show-value", "test")
	assert.HasErr(t, err, errs.ErrReceiver)
	assert.Equal(t, strings.Contains(out, `"value": "hello world!"`), false)

	out, err = cli.RunMain(m, "123\n123\n", "show-value", "value")
	assert.HasErr(t, err, nil)

	json.Unmarshal([]byte(out), &keys)

	assert.Equal(t, len(keys["value"].(string)), 20)

	// show-keys
	out, err = cli.RunMain(m, "", "show-keys", "")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, `[
  "test1",
  "test2",
  "test3"
]`)

	// show-values
	out, err = cli.RunMain(m, "", "show-values", "")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, `[
  "test",
  "value"
]`)

	// rekey
	t.Setenv("ROT_algorithms_asymmetric", string(cryptolib.KDFECDHP256))
	t.Setenv("ROT_algorithms_symmetric", string(cryptolib.EncryptionAES128GCM))

	out, err = cli.RunMain(m, "123\n123\n", "rekey")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "")

	n := defaultCfg()
	n.CLI = c.CLI
	n.Parse(ctx, nil)

	assert.Equal(t, c.PublicKey != n.PublicKey, true)
	assert.Equal(t, n.PublicKey.Key.Algorithm(), cryptolib.AlgorithmECP256Public)
	assert.Equal(t, n.DecryptKeys["test1"].PrivateKey != c.DecryptKeys["test1"].PrivateKey, true)
	assert.Equal(t, n.DecryptKeys["test1"].PrivateKey.KDF, c.DecryptKeys["test1"].PrivateKey.KDF)
	assert.Equal(t, n.Values["value"].Key.Ciphertext != c.Values["value"].Key.Ciphertext, true)
	assert.Equal(t, n.Values["value"].Key.Encryption, cryptolib.EncryptionAES128GCM)
	assert.Equal(t, n.Values["value"].Key.KDF, cryptolib.KDFECDHP256)

	// run
	out, err = cli.RunMain(m, "123\n123\n", "run", "env")
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, "test=***"), true)
	assert.Equal(t, strings.Contains(out, "value=***"), true)

	out, err = cli.RunMain(m, "123\n123\n", "-x", `unmask=["test"]`, "run", "env")
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, "test=hello world!"), true)
	assert.Equal(t, strings.Contains(out, "value=***"), true)

	// add-private-keys
	out, err = cli.RunMain(m, "", "add-pk", "hello", "world")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "")

	prv1, err := cli.RunMain(m, "123\n123\n", "show-value", "-v", "hello")
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Split(prv1, ":")[2], "world")

	pub1, err := cli.RunMain(m, "123\n123\n", "show-value", "-c", "hello")
	assert.HasErr(t, err, nil)

	// show-public-key
	pub2, err := cli.RunMain(m, "", "show-pk", prv1)
	assert.HasErr(t, err, nil)
	assert.Equal(t, pub1, pub2)

	pub2, err = cli.RunMain(m, prv1, "show-pk", "-")
	assert.HasErr(t, err, nil)
	assert.Equal(t, pub1, pub2)

	// pem
	p, _ := cryptolib.ParseKey[cryptolib.KeyProviderPublic](pub1)

	pubPEM, err := cli.RunMain(m, "", "pem", pub1)
	assert.HasErr(t, err, nil)

	pemPub, err := cli.RunMain(m, pubPEM, "pem", "-i", p.ID, "-")
	assert.HasErr(t, err, nil)

	assert.Equal(t, pub1, pemPub)

	// generate-certificate
	exp := 60 * 60
	crtPEM, err := cli.RunMain(m, "hello", "gen-crt", "-c", "-d", "localhost", "-e", strconv.Itoa(exp), "-eu", "ocspSigning", "-eu", "clientAuth", "-ku", "digitalSignature", "-ku", "keyAgreement", "-i", "127.0.0.1", "-n", "My CA", "-")

	assert.HasErr(t, err, nil)

	out, err = cli.RunMain(m, crtPEM, "show-crt", "-")
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, "localhost"), true)

	os.WriteFile("ca.pem", []byte(crtPEM), 0600)

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

	_, pub, _ = cryptolib.NewKeysAsymmetric(cryptolib.BestEncryptionAsymmetric)

	crtPEM, err = cli.RunMain(m, "", "gen-crt", prv1, pub.String(), "ca.pem")

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

	out, err = cli.RunMain(m, crtPEM, "show-crt", "crt.pem", "ca.pem")
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, "My CA"), true)

	crtPEM, _ = cli.RunMain(m, "", "gen-crt", prv1, pub.String(), "ca.pem")
	os.WriteFile("ca.pem", []byte(crtPEM), 0600)

	out, err = cli.RunMain(m, crtPEM, "show-crt", "crt.pem", "ca.pem")
	assert.HasErr(t, err, errs.ErrReceiver)
	assert.Equal(t, strings.Contains(out, "My CA"), true)

	// generate-jwt
	j, err := cli.RunMain(m, prv1, "gen-jwt", "-a", "audience", "-e", "60", "-f", "bool=true", "-f", `string="1"`, "-f", "int=1", "-id", "id", "-is", "issuer", "-s", "subject", "-")
	assert.HasErr(t, err, nil)

	out, err = cli.RunMain(m, "", "show-jwt", j, pub1)
	assert.HasErr(t, err, nil)
	assert.Equal(t, strings.Contains(out, `"audience"`), true)
	assert.Equal(t, strings.Contains(out, `: true`), true)
	assert.Equal(t, strings.Contains(out, `: 1`), true)
	assert.Equal(t, strings.Contains(out, `: "1"`), true)
	assert.Equal(t, strings.Contains(out, `"id"`), true)
	assert.Equal(t, strings.Contains(out, `"issuer"`), true)
	assert.Equal(t, strings.Contains(out, `"subject"`), true)

	out, err = cli.RunMain(m, "", "show-jwt", j)
	assert.HasErr(t, err, jwt.ErrParseNoPublicKeys)
	assert.Equal(t, strings.Contains(out, `"audience"`), true)

	_, jp, _ := cryptolib.NewKeysAsymmetric(cryptolib.AlgorithmBest)

	out, err = cli.RunMain(m, "", "show-jwt", j, jp.String())
	assert.HasErr(t, err, cryptolib.ErrVerify)
	assert.Equal(t, strings.Contains(out, `"audience"`), true)

	j, err = cli.RunMain(m, "", "gen-jwt", "hello")
	assert.HasErr(t, err, nil)

	_, err = cli.RunMain(m, "", "show-jwt", j, pub1)
	assert.HasErr(t, err, nil)

	// ssh
	_, err = cli.RunMain(m, "", "gen-ssh", "-c", "source-address=localhost", "-e", "permit-pty", "-h", "-i", "123", "-p", "root", "-v", "360", "hello", pub1)
	assert.HasErr(t, err, nil)
	_, err = cli.RunMain(m, "", "ssh", pub1)
	assert.HasErr(t, err, nil)

	// remove
	out, err = cli.RunMain(m, "123\n123\n", "add-key", "remove", pk)
	assert.HasErr(t, err, nil)
	assert.Equal(t, out, "")

	_, err = cli.RunMain(m, "", "remove-key", "remove")
	assert.HasErr(t, err, nil)
	_, err = cli.RunMain(m, "", "remove-value", "value")
	assert.HasErr(t, err, nil)

	c.DecryptKeys = map[string]cfgDecryptKey{}
	c.Values = map[string]cfgValue{}
	c.Parse(ctx, nil)
	assert.Equal(t, len(c.DecryptKeys), 3)
	assert.Equal(t, len(c.Values), 2)

	// sig
	cli.RunMain(m, "", "add-pk", "goodbye")
	sig, err := cli.RunMain(m, "helloworld", "gen-sig", "hello", "-")
	assert.HasErr(t, err, nil)
	_, err = cli.RunMain(m, "helloworld", "verify-sig", "hello", "-", sig)
	assert.HasErr(t, err, nil)
	_, err = cli.RunMain(m, "helloworld", "verify-sig", "goodbye", "-", sig)
	assert.HasErr(t, err, errs.ErrReceiver)

	// base64
	out1, err := cli.RunMain(m, "helloworld", "base64", "-")
	assert.HasErr(t, err, nil)
	out2, err := cli.RunMain(m, out1, "base64", "-d", "-")
	assert.HasErr(t, err, nil)
	assert.Equal(t, out2, "helloworld")

	// tamper
	k := n.DecryptKeys["test1"]
	k.PublicKey.ID = "new"
	n.DecryptKeys["test1"] = k
	delete(n.DecryptKeys, "test3")
	n.save(ctx)

	out, err = cli.RunMain(m, "123\n123\n", "show-value", "test")
	assert.HasErr(t, err, errs.ErrReceiver)
	assert.Equal(t, strings.Contains(out, "tampering"), true)

	os.RemoveAll("rot.jsonnet")
	os.RemoveAll(".rot-keys")
	os.RemoveAll("ca.pem")
	os.RemoveAll("crt.pem")
}
