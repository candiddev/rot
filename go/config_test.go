package main

import (
	"context"
	"testing"

	"github.com/candiddev/shared/go/assert"
	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/logger"
)

func TestCfgParse(t *testing.T) {
	logger.UseTestLogger(t)

	tests := map[string]struct {
		c       cfg
		wantErr error
	}{
		"not exist": {
			c: cfg{
				CLI: cli.Config{
					ConfigPath: "not/exist.jsonnet",
				},
			},
			wantErr: errUnknownAlgorithmsAsymmetric,
		},
		"invalid pbkdf": {
			c: cfg{
				Algorithms: cfgAlgorithms{
					Asymmetric: cryptolib.BestEncryptionAsymmetric,
				},
			},
			wantErr: errUnknownAlgorithmsPBKDF,
		},
		"invalid symmetric": {
			c: cfg{
				Algorithms: cfgAlgorithms{
					Asymmetric: cryptolib.BestEncryptionAsymmetric,
					PBKDF:      cryptolib.BestKDF,
				},
			},
			wantErr: errUnknownAlgorithmsSymmetric,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.HasErr(t, tc.c.Parse(context.Background(), nil), tc.wantErr)
		})
	}
}
