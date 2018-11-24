package cmd

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"
	"github.com/pkg/sftp"
	"io"
	"os"
	"path"
)

func SFTP() cli.Command {
	sftpFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "target, t",
			Usage: "The IP/Hostname of the target node",
			Value: "localhost",
		},
		cli.StringFlag{
			Name:   "user, u",
			Usage:  "The user of SFTP connection",
			Value:  "",
			EnvVar: "USER",
		},
		cli.StringFlag{
			Name:  "password, p",
			Usage: "The password of SFTP connection",
			Value: "",
		},
		cli.StringFlag{
			Name:  "file, f",
			Usage: "The path of the remote file",
			Value: "",
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
	conn, err := establishSSHTunnelConnection(ctx)
	if err != nil {
		logrus.Error("Unable to establish ssh connection towards target")
		return err
	}
	defer conn.Close()
	client, err := sftp.NewClient(conn)
	if err != nil {
		panic("Failed to create client: " + err.Error())
	}
	// Close connection
	defer client.Close()
	if ctx.String("file") != "" {
		f, err := client.Open(ctx.String("file"))
		if err != nil {
			logrus.Errorf("Unable to open remote file %s", ctx.String("file"))
			return err
		}
		defer f.Close()
		localFile, err := os.Create(path.Join(os.TempDir(), path.Base(ctx.String("file"))))
		if err != nil {
			logrus.Errorf("Unable to create local file with name %s", path.Base(ctx.String("file")))
			return err
		}
		defer localFile.Close()
		if _, err = io.Copy(localFile, f); err != nil {
			logrus.Error("Unable to copy file to local")
			return err
		}
		logrus.Infof("Remote file downloaded to local successfully")
	}
	return nil
}


func establishSSHTunnelConnection(ctx *cli.Context) (*ssh.Client, error) {
	logrus.Debugf("Establishing SSH connection towards switch %s", ctx.String("target"))
	config := &ssh.ClientConfig{
		User: ctx.String("user"),
		Auth: []ssh.AuthMethod{
			ssh.Password(ctx.String("password")),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ctx.String("target"), 22), config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	return conn, nil
}