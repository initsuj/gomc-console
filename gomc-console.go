package main

import (
	"github.com/initsuj/gomc-console/console"
	"github.com/initsuj/gomc/mcchat"
	"github.com/initsuj/gomc/mcauth"
	"flag"
	"github.com/initsuj/gomc-console/cache"
	"github.com/initsuj/gomc-console/conf"
	"os"
)

var (
	exit = make(chan bool)
)

func main() {

	var ini, server, user, pwd string
	var offline bool

	flag.StringVar(&ini, "i", "gomc.ini", "location of ini file to use.")
	flag.StringVar(&server, "s", "", "server to connect to.")
	flag.StringVar(&user, "u", "", "mojang account username.")
	flag.StringVar(&pwd, "p", "", "mojang account password.")
	flag.BoolVar(&offline, "offline", false, "play in offline mode.")
	flag.Parse()

	if err := cache.Init(); err != nil {
		panic(err)
	}
	defer cache.Close()

	// always try to write file. fails if file exists.
	conf.WriteDefault("gomc.ini")

	if err := conf.Load("gomc.ini"); err != nil {
		console.Println(mcchat.Red, "Error opening ini file: ", err.Error())
	}

	// user by priority: flag -> settings -> user input
	if user == "" {
		user = conf.Current.Connection.Login
	}

	if user == "" {
		user = console.Prompt("Please enter mojang account login: ")
	}

	acct := cache.Find(user)

	// password by priority: flag -> settings -> acct caching -> user input
	if !offline {
		if pwd == "" {
			pwd = conf.Current.Connection.Password
		}

		if pwd == "" && acct.AccessToken == "" {
			pwd = console.ReadPassword("Please enter mojang account password: ")
		}
	}

	authd := offline
	if !authd && acct.AccessToken != "" && acct.Login != "" {
		authd, _ = mcauth.Validate(acct)

		if !authd {
			authd = (mcauth.Refresh(&acct) == nil)
		}

		if authd {
			console.Println("Successfully validated!")
		}else {
			console.Println("Acccount token validation failed!")
		}
	}

	if !authd {

		if acct.Login == "" {
			acct.Login = user
		}

		if acct.ClientToken == "" {
			id, err := mcauth.NewUUID()
			if err != nil {
				panic(err)
			}
			acct.ClientToken = id
		}
		console.Println("Contacting Minecraft.net!")

		// give three attempts to login
		for i := 0; i < 2; i++ {
			err := mcauth.Authenticate(&acct, pwd)
			if err == nil {
				authd = true
				cache.Store(acct)
				break
			}

			aerr, ok := err.(mcauth.AuthError); if !ok || aerr.Type != "ForbiddenOperationException" {
				console.Println(mcchat.Red, "Error while logging into Minecraft: ", err)
				break
			}else {
				console.Println(mcchat.Red, err)

				user = console.Prompt("Please enter mojang account login: ")
				pwd = console.ReadPassword("Please enter mojang account password: ")
			}
		}

		if authd {
			console.Println("Successfully authenticated!")
		}else {
			console.Println(mcchat.Red, "Cannot authenticate!")
			os.Exit(1)
		}
	}

	go console.Scan(func(input string) {
		if input == "/quit" {
			exit <- true
		}
	})

	<-exit
}





