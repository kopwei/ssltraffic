package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/kopwei/ssltraffic/cmd"
	"github.com/urfave/cli"
)

// VERSION is the global value for software version
var VERSION = "v0.1.0-dev"

func main() {
	if err := mainErr(); err != nil {
		logrus.Fatal(err)
	}
}

func mainErr() error {
	app := cli.NewApp()
	app.Name = "ssltraffictest"
	app.Usage = "Establish SSL Traffic"
	app.Version = VERSION
	app.Before = func(ctx *cli.Context) error {
		if ctx.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}
	app.Author = "kopkop"
	app.Commands = []cli.Command{
		cmd.SSH(),
		cmd.SFTP(),
		cmd.HTTPS(),
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug,d",
			Usage: "Debug logging",
		},
	}
	return app.Run(os.Args)
}
