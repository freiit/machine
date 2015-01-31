package pb

import (
	"fmt"
	"os/exec"
	"path/filepath"

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
	IPAddress      string
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
	log.Infof("Creating SSH key...")
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

