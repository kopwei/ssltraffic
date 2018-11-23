package cmd

import "github.com/urfave/cli"

func HTTPS() cli.Command {
	httpsFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "target, t",
			Usage: "The IP/Hostname of the target node",
			Value: "localhost",
		},
	}
	return cli.Command{
		Name:   "sftp",
		Usage:  "Establish https connection to target",
		Action: httpsAction,
		Flags:  httpsFlags,
	}
}

func httpsAction(ctx *cli.Context) error {
	return nil
}
