package conf

import (
	_ "embed"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imssyang/gweb/configs"
	"gopkg.in/yaml.v3"
)

var App AppOpts

func init() {
	if err := App.Parse(configs.Default); err != nil {
		log.Fatalf("config: failed loading defaults config (%s)", err)
		return
	}
}

type AppOpts struct {
	Config  string
	Debug   bool `yaml:"debug"`
	Silent  bool `yaml:"silent"`
	Service struct {
		Address string `yaml:"address"`
		Host    string
		Port    int
		Timeout struct {
			Read  time.Duration `yaml:"read"`
			Write time.Duration `yaml:"write"`
		} `yaml:"timeout"`
	} `yaml:"service"`
	Template struct {
		Files []string `yaml:"files"`
	} `yaml:"template"`
}

func (o *AppOpts) Load(fileName string) error {
	if fileName == "" {
		return nil
	}

	info, err := os.Stat(fileName)
	if err != nil || info.IsDir() {
		return fmt.Errorf("%s not found", fileName)
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	return o.Parse(data)
}

func (o *AppOpts) Parse(data []byte) error {
	err := yaml.Unmarshal(data, o)
	if err != nil {
		return err
	}

	return o.Decap()
}

func (o *AppOpts) Encap() error {
	o.Service.Address = o.Service.Host + ":" + fmt.Sprintf("%d", o.Service.Port)
	return nil
}

func (o *AppOpts) Decap() error {
	host, port, err := net.SplitHostPort(o.Service.Address)
	if err != nil {
		log.Fatalf("config: failed split host and port from %s (%s)", o.Service.Address, err)
		return err
	}

	o.Service.Host = host
	o.Service.Port, err = strconv.Atoi(port)
	return err
}

func GinMode() string {
	if App.Debug {
		return gin.DebugMode
	} else {
		return gin.ReleaseMode
	}
}
