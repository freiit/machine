package pb

import (
	_ "fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/machine/drivers"
)

const (
	dockerConfigDir = "/etc/docker"
)

type Driver struct {
	Userid		   string
	Password	   int
	MachineName    string
	CaCertPath     string
	PrivateKeyPath string
	DriverKeyPath  string
	storePath      string
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

func (d *Driver) Create() error {
	log.Infof("Creating SSH key...")
	return nil
}

