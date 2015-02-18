package pb

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"net/http"
  	"bytes"
  	"io/ioutil"
  	"encoding/xml"
  	"time"
  	"strconv"

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
	User     	   string
	Password	   string
	IPAddress      string
	VDCName 	   string
	StorageSize	   string
	Cores 		   string
	RamSize		   string
	MachineName    string
	CaCertPath     string
	PrivateKeyPath string
	DriverKeyPath  string
	storePath      string
}

type StorageCreateReturn struct {
    RequestId int `xml:"requestId"`
    DataCenterId string `xml:"dataCenterId"`
    DataCenterVersion int `xml:"dataCenterVersion"`
    StorageId string `xml:"storageId"`
}

type ServerCreateReturn struct {
	RequestId int `xml:"requestId"`
    DataCenterId string `xml:"dataCenterId"`
    DataCenterVersion int `xml:"dataCenterVersion"`
    ServerId string `xml:"serverId"`
}

type VDCGetReturn struct {
	DataCenterId string `xml:"dataCenterId"`
	DataCenterName string `xml:"dataCenterName"`
	DataCenterVersion int `xml:"dataCenterVersion"`
	ProvisioningState string `xml:"provisioningState"`
}

type ConnectedStrg struct {
	BootDevice bool `xml:"bootDevice"`
	BusType string `xml:"busType"`
	DeviceNumber int `xml:"deviceNumber"`
	Size int `xml:"size"`
	StorageId string `xml:"storageId"`
	StorageName string `xml:"storageName"`
}

type Frwl struct {
	Active bool `xml:"active"`
	FirewallId string `xml:"firewallId"`
	NicId string `xml:"nicId"`
	ProvisioningState string `xml:"provisioningState"`
}

type Nic struct {
	DataCenterId string `xml:"dataCenterId"`
	DataCenterVersion int `xml:"dataCenterVersion"`
	NicId string `xml:"nicId"`
	LanId int `xml:"lanId"`
	InternetAccess bool `xml:"internetAccess"`
	ServerId string `xml:"serverId"`
	Ips string `xml:"ips"`
	MacAddress string `xml:"macAddress"`
	Firewall Frwl `xml:"firewall"`
	DhcpActive bool `xml:"dhcpActive"`
	GatewayIp string `xml:"gatewayIp"`
	ProvisioningState string `xml:"provisioningState"`
}

type GetServerCallReturn struct {
	RequestId int `xml:"requestId"`
	DataCenterId string `xml:"dataCenterId"`
	DataCenterVersion int `xml:"dataCenterVersion"`
	ServerId string `xml:"serverId"`
	ServerName string `xml:"serverName"`
	Cores int `xml:"cores"`
	Ram int `xml:"ram"`
	InternetAccess bool `xml:"internetAccess"`
	Ips string `xml:"ips"`
	ConnectedStorages ConnectedStrg `xml:"connectedStorages"`
	Nics Nic `xml:"nics"`
	ProvisioningState string `xml:"provisioningState"`
	VirtualMachineState string `xml:"virtualMachineState"`
	CreationTime string `xml:"creationTime"`
	LastModificationTime string `xml:"lastModificationTime"`
	OsType string `xml:"osType"`
	AvailabilityZone string `xml:"availabilityZone"`
	CpuHotPlug bool `xml:"cpuHotPlug"`
	RamHotPlug bool `xml:"ramHotPlug"`
	NicHotPlug bool `xml:"nicHotPlug"`
	NicHotUnPlug bool `xml:"nicHotUnPlug"`
	DiscVirtioHotPlug bool `xml:"discVirtioHotPlug"`
	DiscVirtioHotUnPlug bool `xml:"discVirtioHotUnPlug"`
}

type StorageReturn struct{
    Ret  StorageCreateReturn `xml:"return"`
}

type ServerReturn struct {
	Ret ServerCreateReturn `xml:"return"`
}

type GetServerReturn struct {
	Ret GetServerCallReturn `xml:"return"`
}

type VDCGetAllReturn struct {
	Ret VDCGetReturn `xml:"return"`
}

type VDCBody struct {
	VDCResposne VDCGetAllReturn `xml:"getAllDataCentersResponse"`
}

type StorageBody struct {
    StrgRet StorageReturn `xml:"createStorageReturn"`
}

type ServerBody struct {
	ServerRet ServerReturn `xml:"createServerReturn"`
}

type GetServerBody struct {
	GetServerResponse GetServerReturn `xml:"getServerResponse"`
}

type StorageResponse struct{
    XMLName xml.Name `xml:"Envelope"`
    RespBody StorageBody `xml:"Body"`
}

type ServerResponse struct{
	XMLName xml.Name `xml:"Envelope"`
    RespBody ServerBody `xml:"Body"`
}

type GetServerResponse struct{
	XMLName xml.Name `xml:"Envelope"`
    RespBody GetServerBody `xml:"Body"`
}

type VDCResponse struct {
	XMLName xml.Name `xml:"Envelope"`
    RespBody VDCBody `xml:"Body"`
}

func makeReq(reqStr string, userid string, pass string) string{
	buf := []byte(reqStr)
	body := bytes.NewBuffer(buf)
	client := &http.Client{}
    req, err := http.NewRequest("POST", "https://api.profitbricks.com/1.3", body)
    if err != nil {
    	log.Errorf("Error in creating http client")
    	log.Errorf("%v", err)
    	return ""
    }
    req.SetBasicAuth(userid, pass)
    resp, err := client.Do(req)
    if err != nil{
        log.Errorf("Error in calling pb api")
        log.Errorf("%v", err)
    	return ""
    }
    bodyText, err := ioutil.ReadAll(resp.Body)
    if err != nil{
        log.Errorf("Error in response")
        log.Errorf("%v", err)
    	return ""
    }
    s := string(bodyText)
    //fmt.Printf("%v", s)
    return s
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
		cli.StringFlag{
			EnvVar: "PB_DCNAME",
			Name: "pb-vdc-name",
			Usage: "Profitbicks data centre name",
		},
		cli.StringFlag{
			EnvVar: "PB_STORAGE",
			Name:   "pb-storagesizeGB",
			Usage: "Profitbricks Virtual Server storage space size",
		},
		cli.StringFlag{
			EnvVar: "PB_CORES",
			Name:   "pb-cores",
			Usage: "Profitbricks Virtual Server compute cores",
		},
		cli.StringFlag{
			EnvVar: "PB_RAM",
			Name:   "pb-ramGB",
			Usage: "Profitbricks Virtual Server RAM size",
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
	d.User = flags.String("pb-user")
	d.Password = flags.String("pb-password")
	d.VDCName = flags.String("pb-vdc-name")
	d.StorageSize = flags.String("pb-storagesizeGB")
	d.Cores = flags.String("pb-cores")
	d.RamSize = flags.String("pb-ramGB")
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
	key, err := d.createSSHKey()
	if err != nil {
		return err
	}
	log.Infof("User ---- %v", d)
	//Get vdc ID from name
	log.Infof("%s", d.VDCName)
	log.Debugf(" ssssssssssssssss -------------- %+v", key)	
	return nil

	soapreq_str := `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ws="http://ws.api.profitbricks.com/">
					<soapenv:Header>
					</soapenv:Header>
					<soapenv:Body>					
					<ws:getAllDataCenters>
					<request>
					</request>
					</ws:getAllDataCenters>
					</soapenv:Body>
					</soapenv:Envelope>`
	s := makeReq(soapreq_str, d.User, d.Password)
	if s == ""{
		log.Debugf("Error Happened while getting VDC-------------------")
		log.Debugf(s)
		return nil
	}
	//fmt.Printf("%s", s)
	v3 := VDCResponse{}
    err = xml.Unmarshal([]byte(s), &v3)
	if err != nil {
		log.Infof("error: %v", err)
		log.Infof("Return XML  - %s", s)
		return err
	}
	log.Infof("%s", v3.RespBody.VDCResposne.Ret.DataCenterName)
	vdcId := ""
	if v3.RespBody.VDCResposne.Ret.DataCenterName == d.VDCName{
		vdcId = v3.RespBody.VDCResposne.Ret.DataCenterId
	}
	if vdcId == "" {
		log.Errorf("Could not find the data center named - %s", v3.RespBody.VDCResposne.Ret.DataCenterName)
		return nil
	}

	//create storage
	// Assumed the region is us/las
	//Need to recode that part for all the region
	soapreq_str = `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ws="http://ws.api.profitbricks.com/">
					<soapenv:Header>
					</soapenv:Header>
					<soapenv:Body>					
					<ws:createStorage>
					<request>
					<size>%s</size>
					<dataCenterId>%s</dataCenterId>
					<mountImageId>b60ca23c-a9c5-11e4-9f44-52540066fee9</mountImageId>
					<profitBricksImagePassword>mEs234Ppq</profitBricksImagePassword>
					</request>
					</ws:createStorage>
					</soapenv:Body>
					</soapenv:Envelope>`
	
	soapreq_str = fmt.Sprintf(soapreq_str, d.StorageSize, vdcId)
	s = makeReq(soapreq_str, d.User, d.Password)
	if s == ""{
		log.Errorf("Error Happened while creting the storage-----------------")
		log.Debugf(s)
		return nil
	}
	v := StorageResponse{}
    err = xml.Unmarshal([]byte(s), &v)
	if err != nil {
		log.Infof("error: %v", err)
		log.Infof("Return XML  - %s", s)
		return err
	}
	if v.RespBody.StrgRet.Ret.StorageId == ""{
		log.Errorf("Could not unmarshal the XML")
		log.Errorf(s)
		return nil	
	}

	//Create server
	i, err := strconv.Atoi(d.RamSize)
	if err != nil {
        // handle error
        log.Errorf("RAM size must be an integer")
        return err
    }
    i = i * 1024 //GB to MB as pb accepts this param in MB
    t := strconv.Itoa(i)

	soapreq_str = `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ws="http://ws.api.profitbricks.com/">
					<soapenv:Header>
					</soapenv:Header>
					<soapenv:Body>					
					<ws:createServer>
					<request>
					<cores>%s</cores>
					<ram>%s</ram>
					<dataCenterId>%s</dataCenterId>
					<bootFromStorageId>%s</bootFromStorageId>
					<serverName>%s</serverName>
					<internetAccess>true</internetAccess>
					<availabilityZone>AUTO</availabilityZone>
					<osType>OTHER</osType>
					<cpuHotPlug>true</cpuHotPlug>
					<ramHotPlug>true</ramHotPlug>
					<nicHotPlug>true</nicHotPlug>
					<nicHotUnPlug>true</nicHotUnPlug>
					<discVirtioHotPlug>true</discVirtioHotPlug>
					<discVirtioHotUnPlug>true</discVirtioHotUnPlug>
					</request>
					</ws:createServer>
					</soapenv:Body>
					</soapenv:Envelope>`
	soapreq_str = fmt.Sprintf(soapreq_str, d.Cores, t, vdcId, v.RespBody.StrgRet.Ret.StorageId, d.MachineName)
	s = makeReq(soapreq_str, d.User, d.Password)
	//fmt.Printf("%s", s)
	v1 := ServerResponse{}
    err = xml.Unmarshal([]byte(s), &v1)
	if err != nil {
		log.Infof("error: %v", err)
		log.Infof("Return XML  - %s", s)
		return err
	}
	if v1.RespBody.ServerRet.Ret.ServerId == ""{
		log.Errorf("Could not unmarshal the XML")
		log.Errorf(s)
		return nil	
	}

	//Ping the server to see it is ready
	i = 0
	for {
		log.Infof("Pinging the server ....")
		soapreq_str = `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ws="http://ws.api.profitbricks.com/">
						<soapenv:Header>
						</soapenv:Header>
						<soapenv:Body>					
						<ws:getServer>
						<serverId>%s</serverId>
						</ws:getServer>
						</soapenv:Body>
						</soapenv:Envelope>`
		soapreq_str = fmt.Sprintf(soapreq_str, v1.RespBody.ServerRet.Ret.ServerId)
		s = makeReq(soapreq_str, d.User, d.Password)
		v2 := GetServerResponse{}
		err = xml.Unmarshal([]byte(s), &v2)
		if err != nil {
			log.Infof("error: %v", err)
			log.Infof("Return XML  - %s", s)
			return err
		}
		log.Infof("%s", v2.RespBody.GetServerResponse.Ret.VirtualMachineState)
		if v2.RespBody.GetServerResponse.Ret.VirtualMachineState == "RUNNING"{
			d.IPAddress = v2.RespBody.GetServerResponse.Ret.Ips
			break
		}
		i = i + 1
		time.Sleep(10 * time.Second)
		if i == 100 { //Not more than 100 try
			break
		}
	}

	//Run docker command
	IPCommand := "-e IP="+d.IPAddress
	Passcommand := "-e PASSWORD=mEs234Ppq"
	DockerImage := "freiit/dockerize"
	cmd := exec.Command("docker", "run", IPCommand, Passcommand, DockerImage)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Infof("Error running docker command")
		log.Errorf("%v", err.Error())
		return err
	}
	log.Infof("******************* OUTPUT *******************")
	log.Infof("%q\n", out.String())
	log.Infof("******************* OUTPUT *******************")
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

func (d *Driver) createSSHKey() (string, error) {

	if err := ssh.GenerateSSHKey(d.sshKeyPath()); err != nil {
		return "", err
	}

	publicKey, err := ioutil.ReadFile(d.publicSSHKeyPath())
	if err != nil {
		return "", err
	}

	return string(publicKey), nil
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

func (d *Driver) publicSSHKeyPath() string {
	return d.sshKeyPath() + ".pub"
}
