package main

import (
	"context"
	"encoding/base64"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/logger"
)

func cmdBase64() cli.Command[*cfg] {
	return cli.Command[*cfg]{
		ArgumentsRequired: []string{
			"input value or - for stdin",
		},
		Flags: cli.Flags{
			"d": {
				Usage: "Decode base64 (default: encode)",
			},
			"r": {
				Usage: "Raw/no padding (default: padding)",
			},
			"u": {
				Usage: "URL encoding (default: standard encoding)",
			},
		},
		Usage: "Encode/decode a base64 value or stdin and output to stdout.",
		Run: func(ctx context.Context, args []string, f cli.Flags, c *cfg) errs.Err {
			v := []byte(args[1])
			if string(v) == "-" {
				v = cli.ReadStdin()
			}

			_, d := f.Value("d")
			_, url := f.Value("u")
			_, r := f.Value("r")

			var err error

			var out []byte

			switch {
			case d && url && r:
				out, err = base64.RawURLEncoding.DecodeString(string(v))
			case d && r:
				out, err = base64.RawStdEncoding.DecodeString(string(v))
			case d && url:
				out, err = base64.URLEncoding.DecodeString(string(v))
			case d:
				out, err = base64.StdEncoding.DecodeString(string(v))
			case r && url:
				out = []byte(base64.RawURLEncoding.EncodeToString(v))
			case r:
				out = []byte(base64.RawStdEncoding.EncodeToString(v))
			case url:
				out = []byte(base64.URLEncoding.EncodeToString(v))
			default:
				out = []byte(base64.StdEncoding.EncodeToString(v))
			}

			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			logger.Stdout.Write(out) //nolint:errcheck

			return nil
		},
	}
}
