package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/candiddev/shared/go/cli"
	"github.com/candiddev/shared/go/errs"
	"github.com/candiddev/shared/go/jwt"
	"github.com/candiddev/shared/go/logger"
	"github.com/google/uuid"
)

func cmdGenJWT() cli.Command[*cfg] {
	return cli.Command[*cfg]{
		ArgumentsRequired: []string{
			"private key value, encrypted value name, or - for stdin",
		},
		Flags: cli.Flags{
			"a": {
				Placeholder: "audience",
				Usage:       "Audience (aud) for JWT.  Can be provided multiple times",
			},
			"e": {
				Default:     []string{"3600"},
				Placeholder: "seconds",
				Usage:       "Expiration (exp) in seconds",
			},
			"f": {
				Placeholder: "key=value",
				Usage:       "Add a key and value to the JWT.  Will attempt to parse bools and ints unless they are quoted.  Can be provided multiple times.",
			},
			"id": {
				Placeholder: "id",
				Usage:       "ID (jti) of the JWT, will generate a UUID if not provided",
			},
			"is": {
				Default:     []string{"Rot"},
				Placeholder: "issuer",
				Usage:       "Issuer (iss) of the JWT",
			},
			"s": {
				Placeholder: "subject",
				Usage:       "Subject (sub) of the JWT",
			},
		},
		Usage: "Generate a JWT and output it it to stdout.  Must specify the private key to sign the JWT.",
		Run: func(ctx context.Context, args []string, f cli.Flags, c *cfg) errs.Err {
			m := map[string]any{}

			aud, _ := f.Values("a")

			var expires time.Time

			e, _ := f.Value("e")
			if e != "" {
				var err error

				s, err := strconv.Atoi(e)
				if err != nil {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(fmt.Errorf("error parsing -e: %w", err)))
				}

				expires = time.Now().Add(time.Second * time.Duration(s))
			}

			id, _ := f.Value("id")
			if id == "" {
				id = uuid.New().String()
			}

			issuer, _ := f.Value("is")
			subject, _ := f.Value("s")

			vals, _ := f.Values("f")

			for i := range vals {
				s := strings.Split(vals[i], "=")
				if len(s) <= 1 {
					return logger.Error(ctx, errs.ErrReceiver.Wrap(fmt.Errorf("error parsing field: %s: invalid format, must be key=value", vals[i])))
				}

				value := strings.Join(s[1:], "=")

				if strings.HasPrefix(value, `"`) {
					str, err := strconv.Unquote(value)
					if err != nil {
						return logger.Error(ctx, errs.ErrReceiver.Wrap(fmt.Errorf("error parsing field: %s: %w", vals[i], err)))
					}

					m[s[0]] = str
				} else if b, err := strconv.ParseBool(value); err == nil {
					m[s[0]] = b
				} else if i, err := strconv.Atoi(value); err == nil {
					m[s[0]] = i
				} else {
					m[s[0]] = value
				}
			}

			j, _, err := jwt.New(m, expires, aud, id, issuer, subject)
			if err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			pk := args[1]
			if pk == "-" {
				pk = string(cli.ReadStdin())
			}

			privateKey, errr := c.decryptValuePrivateKey(ctx, pk)
			if errr != nil {
				return logger.Error(ctx, errr)
			}

			if err := j.Sign(privateKey); err != nil {
				return logger.Error(ctx, errs.ErrReceiver.Wrap(err))
			}

			logger.Raw(j.String())

			return nil
		},
	}
}
