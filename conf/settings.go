package conf

import (
	"io/ioutil"
	"github.com/zackbloom/go-ini"
	"os"
)

var (
	SettingsLoaded = false
	Current = Settings{
		Connection: struct {
			Server   string
			Login    string
			Password string
			Cache    bool
		}{
			Server: "",
			Login: "",
			Password: "",
			Cache: true,
		},
	}
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

		if err := ini.Unmarshal(c, &Current); err != nil {
			return err
		}

		SettingsLoaded = true
	}

	return nil
}

func WriteDefault(filename string) error {
	ini := `[Connection]
#Server = minecraft server that you would like to connect with. ex "us.mineplex.com", "192.168.1.100:25864"
#Login = Minecraft.net username
#Password = Minecraft.net password
#Cache = Whether account caching is enabled - [true/false]

Server=localhost
Login=
Password=
Cache=true
`
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err := ioutil.WriteFile(filename, []byte(ini), 7777)
		if err != nil {
			return err
		}
	}

	return nil
}