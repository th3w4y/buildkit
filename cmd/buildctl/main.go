package main

import (
	"fmt"
	"os"

	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/util/appdefaults"
	"github.com/moby/buildkit/util/profiler"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "buildctl"
	app.Usage = "build utility"

	defaultAddress := os.Getenv("BUILDKIT_HOST")
	if defaultAddress == "" {
		defaultAddress = appdefaults.Address
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output in logs",
		},
		cli.StringFlag{
			Name:  "addr",
			Usage: "buildkitd address",
			Value: defaultAddress,
		},
		cli.StringFlag{
			Name:  "server-name",
			Usage: "buildkitd server name for certificate validation",
			Value: "",
		},
		cli.StringFlag{
			Name:  "ca-cert",
			Usage: "CA certificate for validation",
			Value: "",
		},
		cli.StringFlag{
			Name:  "cert",
			Usage: "client certificate",
			Value: "",
		},
		cli.StringFlag{
			Name:  "key",
			Usage: "client key",
			Value: "",
		},
	}

	app.Commands = []cli.Command{
		diskUsageCommand,
		buildCommand,
		debugCommand,
	}

	var debugEnabled bool

	app.Before = func(context *cli.Context) error {
		debugEnabled = context.GlobalBool("debug")
		if debugEnabled {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}

	profiler.Attach(app)

	if err := app.Run(os.Args); err != nil {
		if debugEnabled {
			fmt.Fprintf(os.Stderr, "error: %+v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
		os.Exit(1)
	}
}

func resolveClient(c *cli.Context) (*client.Client, error) {
	serverName := c.GlobalString("server-name")
	caCert := c.GlobalString("ca-cert")
	cert := c.GlobalString("cert")
	key := c.GlobalString("key")
	opts := []client.ClientOpt{client.WithBlock()}
	if serverName != "" {
		opts = append(opts, client.WithCredentials(serverName, caCert, cert, key))
	}
	return client.New(c.GlobalString("addr"), opts...)
}
