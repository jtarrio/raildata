package cmd

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/jtarrio/raildata"
	"github.com/rogpeppe/go-internal/lockedfile"
	"github.com/urfave/cli/v2"
)

func App() *cli.App {
	return &cli.App{
		Name:  "raildata-cli",
		Usage: "An application to query the RailData API",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "tokenfile",
				Usage:    "the pathname of a file containing the RailData API token. If the token is updated, the new value will be written to this file",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "username",
				Usage:   "the RailData API user name",
				EnvVars: []string{"RAILDATA_USERNAME"},
			},
			&cli.StringFlag{
				Name:    "password",
				Usage:   "the RailData API password",
				EnvVars: []string{"RAILDATA_PASSWORD"},
			},
			&cli.BoolFlag{
				Name:  "use-test-endpoint",
				Usage: "use the RailData test endpoint",
			},
		},
		Before: createClient,
		Commands: []*cli.Command{
			cmdGetStationMsg,
			cmdGetStationSchedule,
			cmdGetTrainSchedule,
			cmdGetTrainSchedule19Rec,
			cmdGetTrainStopList,
			cmdGetVehicleData,
		},
	}
}

func createClient(ctx *cli.Context) error {
	var options []raildata.Option

	if ctx.Bool("use-test-endpoint") {
		options = append(options, raildata.WithTestEndpoint(true))
	}
	tokenfile := ctx.String("tokenfile")
	token, err := readTokenFile(tokenfile)
	if err != nil {
		return err
	}
	options = append(options, raildata.WithToken(token))
	options = append(options, raildata.WithTokenUpdateListener(tokenFileUpdater(tokenfile)))

	username := ctx.String("username")
	password := ctx.String("password")
	if (len(username) == 0) != (len(password) == 0) {
		return errors.New("you must specify both --username and --password or none of them")
	}
	if len(username) > 0 {
		options = append(options, raildata.WithCredentials(username, password))
	}

	client, err := raildata.NewClient(options...)
	if err != nil {
		return err
	}
	ctx.Context = context.WithValue(ctx.Context, clientKey, client)
	return nil
}

func readTokenFile(name string) (string, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	token, _, _ := strings.Cut(string(b), "\n")
	return token, nil
}

func tokenFileUpdater(name string) raildata.TokenUpdateListener {
	return func(newToken string, oldToken string) {
		_ = lockedfile.Transform(name, func(old []byte) ([]byte, error) {
			token, _, _ := strings.Cut(string(old), "\n")
			if token != oldToken {
				return old, errors.New("")
			}
			return []byte(newToken + "\n"), nil
		})
	}
}

type clientKeyType struct{}

var clientKey = clientKeyType{}

func GetClientFromContext(ctx context.Context) raildata.Client {
	value := ctx.Value(clientKey)
	if o, ok := value.(raildata.Client); ok {
		return o
	}
	panic("No client found in context")
}
