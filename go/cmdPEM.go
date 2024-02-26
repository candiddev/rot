package main

import (
	"context"
	"strings"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/cryptolib"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdPEM(ctx context.Context, args []string, f cli.Flags, _ *cfg) errs.Err {
	s := []byte(args[1])
	if args[1] == "-" {
		s = cli.ReadStdin()
	}

	var out string

	if strings.HasPrefix(string(s), "--") {
		s, err := cryptolib.PEMToKey[cryptolib.KeyProvider](s)
		if err != nil {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}

		s.ID, _ = f.Value("i")

		out = s.String() + "\n"
	} else {
		k, err := cryptolib.ParseKey[cryptolib.KeyProvider](string(s))
		if err != nil {
			return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
		}

		out = string(cryptolib.KeyToPEM(k))
	}

	logger.Raw(out)

	return nil
}
