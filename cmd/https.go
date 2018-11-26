package cmd

import (
	"bytes"
	"net/http"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func HTTPS() cli.Command {
	httpsFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "target, t",
			Usage: "The IP/Hostname of the target node",
			Value: "localhost",
		},
	}
	return cli.Command{
		Name:   "https",
		Usage:  "Establish https connection to target",
		Action: httpsAction,
		Flags:  httpsFlags,
	}
}

func httpsAction(ctx *cli.Context) error {
	client := &http.Client{}
	reqURL := &url.URL{
		Host:   ctx.String("target"),
		Path:   "/",
		Scheme: "https",
	}
	req, err := http.NewRequest("GET", reqURL.String(), bytes.NewBuffer([]byte("")))
	if err != nil {
		logrus.Error(err)
		return err
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("%v", resp)
	return nil
}
