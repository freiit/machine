package pb

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"net/http"
  	"bytes"
  	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/machine/drivers"
	"github.com/docker/machine/ssh"
	"github.com/docker/machine/state"
)

const (
	dockerConfigDir = "/etc/docker"
)

type Driver struct {
	Userid		   string
	Password	   string
	MachineName    string
	CaCertPath     string
	PrivateKeyPath string
	DriverKeyPath  string
	storePath      string
	//IPAddress      string
}

func init() {
	drivers.Register("pb", &drivers.RegisteredDriver{
		New:            NewDriver,
		GetCreateFlags: GetCreateFlags,
	})
}

func GetCreateFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			EnvVar: "PB_USER",
			Name:   "pb-user",
			Usage:  "Profiitbricks user id",
		},
		cli.StringFlag{
			EnvVar: "PB_PASSWD",
			Name:   "pb-password",
			Usage:  "Profitbricks password",
		},
	}
}

func NewDriver(machineName string, storePath string, caCert string, privateKey string) (drivers.Driver, error) {
	return &Driver{MachineName: machineName, storePath: storePath, CaCertPath: caCert, PrivateKeyPath: privateKey}, nil
}

func (d *Driver) DriverName() string {
	return "pb"
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.Userid = flags.String("pb-user")
	d.Password = flags.String("pb-password")
	return nil
}

func (d *Driver) PreCreateCheck() error {
	return nil
}

///////////////
// CREATE
//////////////

func (d *Driver) Create() error {
	//d.IPAddress = "127.0.0.1"
	log.Debugf("1")
	soapreq_str := `<soapenv:Envelope xmlns:soapenv=”http://schemas.xmlsoap.org/soap/envelope/” xmlns:ws=”http://ws.api.profitbricks.com/”>
					<soapenv:Header>
					</soapenv:Header>
					<soapenv:Body>
					</soapenv:Body>
					<ws:getAllDataCenters>
					<request>
					</request>
					</ws:getAllDataCenters>
					</soapenv:Envelope>`
	buf := []byte(soapreq_str)
	log.Debugf("2")
	body := bytes.NewBuffer(buf)
	client := &http.Client{}
    req, err := http.NewRequest("POST", "https://api.profitbricks.com/1.3", body)
    if err != nil {
    	log.Debugf("Error in creating http client")
    	log.Debugf("%v", err)
    	return err
    }
    log.Debugf("3")
    req.SetBasicAuth(d.Userid, d.Password)
    resp, err := client.Do(req)
    if err != nil{
        log.Debugf("Error in calling pb api")
        log.Debugf("%v", err)
    	return err
    }
    log.Debugf("4")
    bodyText, err := ioutil.ReadAll(resp.Body)
    if err != nil{
        log.Debugf("Error in response")
        log.Debugf("%v", err)
    	return err
    }
    log.Debugf("5")
    s := string(bodyText)
    log.Debugf("%v", s)
    log.Debugf("6")
	return nil
}

////////////////
// GET STATE
///////////////
func (d *Driver) GetState() (state.State, error) {
	return state.None, nil
}

////////////////
// Kill
///////////////

func (d *Driver) Kill() error {
	return nil
}

///////////////
// Remove
//////////////

func (d *Driver) Remove() error {
	return nil
}

//////////////
// Restart
/////////////

func (d *Driver) Restart() error {
	return nil
}

/////////////
// Start
/////////////
func (d *Driver) Start() error {
	return nil
}

//////////////
// Start docker
//////////////

func (d *Driver) StartDocker() error {
	log.Debug("Starting Docker...")

	cmd, err := d.GetSSHCommand("sudo service docker start")
	if err != nil {
		return err
	}
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

///////////////
// Stop
//////////////

func (d *Driver) Stop() error {
	return nil
}

///////////////
// Stop docker
//////////////

func (d *Driver) StopDocker() error {
	log.Debug("Stopping Docker...")

	cmd, err := d.GetSSHCommand("sudo service docker stop")
	if err != nil {
		return err
	}
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

//////////////
// Upgrade
/////////////

func (d *Driver) Upgrade() error {
	log.Debugf("Upgrading Docker")

	cmd, err := d.GetSSHCommand("sudo apt-get update && apt-get install --upgrade lxc-docker")
	if err != nil {
		return err

	}
	if err := cmd.Run(); err != nil {
		return err

	}

	return cmd.Run()
}

func (d *Driver) GetURL() (string, error) {
	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("tcp://%s:2376", ip), nil
}

func (d *Driver) GetIP() (string, error) {
	if d.IPAddress == "" {
		return "", fmt.Errorf("IP address is not set")
	}
	return d.IPAddress, nil
}


func (d *Driver) GetSSHCommand(args ...string) (*exec.Cmd, error) {
	return ssh.GetSSHCommand(d.IPAddress, 22, "root", d.sshKeyPath(), args...), nil
}

func (d *Driver) sshKeyPath() string {
	return filepath.Join(d.storePath, "id_rsa")
}

func (d *Driver) GetDockerConfigDir() string {
	return dockerConfigDir
}

