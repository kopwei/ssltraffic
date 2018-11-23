package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"os/user"
	"time"
)

func SSH() cli.Command {
	sshFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "target, t",
			Usage: "The IP/Hostname of the target node",
			Value: "localhost",
		},
		cli.StringFlag{
			Name:   "user, u",
			Usage:  "The user of SSH command",
			Value:  "",
			EnvVar: "USER",
		},
		cli.StringFlag{
			Name:  "password, p",
			Usage: "The password of SSH command",
			Value: "",
		},
		cli.StringFlag{
			Name:  "keyfile, k",
			Usage: "The path to private key file",
			Value: "",
		},
	}
	return cli.Command{
		Name:   "ssh",
		Usage:  "Establish SSH connection to target",
		Action: sshAction,
		Flags:  sshFlags,
	}
}

func sshAction(ctx *cli.Context) error {
	err := genSSHTraffic(ctx)

	return err
}

type dialer struct {
	sshKeyString string
	sshAddress   string
	username     string
	password     string
}

// Conn wraps a net.Conn, and sets a deadline for every read
// and write operation.
type Conn struct {
	net.Conn
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (c *Conn) Read(b []byte) (int, error) {
	err := c.Conn.SetReadDeadline(time.Now().Add(c.ReadTimeout))
	if err != nil {
		return 0, err
	}
	return c.Conn.Read(b)
}

func (c *Conn) Write(b []byte) (int, error) {
	err := c.Conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
	if err != nil {
		return 0, err
	}
	return c.Conn.Write(b)
}

func getSSHConfig(username, password, sshPrivateKeyString string) (*ssh.ClientConfig, error) {
	signer, err := getPrivateKeySigner(sshPrivateKeyString)
	if err != nil {
		return nil, err
	}
	auth := []ssh.AuthMethod{ssh.PublicKeys(signer)}
	if password != "" {
		auth = append(auth, ssh.Password(password))
	}
	config := &ssh.ClientConfig{
		User:            username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            auth,
	}
	return config, nil
}

func getPrivateKeySigner(sshPrivateKeyString string) (ssh.Signer, error) {
	key, err := parsePrivateKey(sshPrivateKeyString)

	return key, err
}

func parsePrivateKey(keyBuff string) (ssh.Signer, error) {
	return ssh.ParsePrivateKey([]byte(keyBuff))
}

func getSSHClient(d *dialer) (*ssh.Client, error) {
	cfg, err := getSSHConfig(d.username, d.password, d.sshKeyString)
	if err != nil {
		logrus.Error("unable to get ssh configuration")
		return nil, err
	}
	conn, err := net.DialTimeout("tcp", d.sshAddress+":22", 30*time.Second)
	if err != nil {
		return nil, err
	}
	timeoutConn := &Conn{conn, 30 * time.Second, 30 * time.Second}
	c, chans, reqs, err := ssh.NewClientConn(timeoutConn, d.sshAddress+":22", cfg)
	if err != nil {
		return nil, err
	}
	client := ssh.NewClient(c, chans, reqs)
	// this sends keepalive packets every 2 seconds
	// there's no useful response from these, so we can just abort if there's an error
	go func() {
		t := time.NewTicker(2 * time.Second)
		defer t.Stop()
		for {
			<-t.C
			_, _, err := client.Conn.SendRequest("keepalive@golang.org", true, nil)
			if err != nil {
				return
			}
		}
	}()
	return client, nil
}

func getKeyStringByPath(path string) (string, error) {
	if path == "" {
		usr, err := user.Current()
		if err != nil {
			logrus.Error("unable to retrieve current user")
			return "", err
		}
		path = usr.HomeDir + "/.ssh/id_rsa"
	}
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Error("unable to read user's private key")
		return "", err
	}
	return string(dat), nil
}

func genSSHTraffic(ctx *cli.Context) error {
	keyString, err := getKeyStringByPath(ctx.String("keyfile"))
	if err != nil {
		logrus.Warn("Key file reading failed, won't use it")
	}
	d := &dialer{username: ctx.String("user"),
		password:   ctx.String("password"),
		sshAddress: ctx.String("target"),
		sshKeyString: keyString,
		}
	logrus.Infof("Starting to establish ssh conn to host %s", d.sshAddress)
	c, err := getSSHClient(d)
	if err != nil {
		logrus.Errorf("unable to dial addr %s", d.sshAddress)
		return err
	}
	defer disconnectSSH(c)
	session, err := c.NewSession()
	if err != nil {
		logrus.Error("failed to get new session")
		return err
	}
	defer session.Close()
	logrus.Info("Finished ssh connection towards target machine")
	return nil
}

func disconnectSSH(c *ssh.Client) error {
	// Disconnect
	var err error
	if c != nil {
		err = c.Close()
		if err != nil {
			logrus.Warnf("error occurred during closing ssh conn with host %s", c.RemoteAddr().String())
		}
		c = nil
	}
	return err
}
