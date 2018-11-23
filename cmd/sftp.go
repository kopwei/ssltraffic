package cmd

import "github.com/urfave/cli"

func SFTP() cli.Command {
	sftpFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "target, t",
			Usage: "The IP/Hostname of the target node",
			Value: "localhost",
		},
	}
	return cli.Command{
		Name:   "sftp",
		Usage:  "Establish sftp connection to target",
		Action: sftpAction,
		Flags:  sftpFlags,
	}
}

func sftpAction(ctx *cli.Context) error {
	return nil
}
