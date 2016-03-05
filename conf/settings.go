package conf

import (
	"io/ioutil"
	"github.com/zackbloom/go-ini"
	"os"
)

var (
	CurrentSettings Settings
)

type Settings struct {
	Connection struct {
			   Server   string
			   Login    string
			   Password string
			   Cache    bool
		   } `ini:"[Connection]"`
}

func Load(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		c, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		if err := ini.Unmarshal(c, &CurrentSettings); err != nil {
			return err
		}
	}

	return nil
}