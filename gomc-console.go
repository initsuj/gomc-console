package main

import (
	"github.com/initsuj/gomc-console/console"
	"github.com/initsuj/gomc/mcchat"
	"github.com/initsuj/gomc/mcauth"
	"github.com/initsuj/gomc/mcauth/mcrequest"
	"flag"
	"github.com/initsuj/gomc-console/cache"
	"github.com/initsuj/gomc-console/conf"
)

var (
	exit = make(chan bool)
)

func main() {

	var ini, server, user, pwd string

	flag.StringVar(&ini, "i", "gomc.ini", "location of ini file to use.")
	flag.StringVar(&server, "s", "", "server to connect to.")
	flag.StringVar(&user, "u", "", "mojang account username.")
	flag.StringVar(&pwd, "p", "", "mojang account password.")
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
	if pwd == "" {
		pwd = conf.Current.Connection.Password
	}

	if pwd == "" && acct.AccessToken == "" {
		pwd = console.ReadPassword("Please enter mojang account password: ")
	}

	authd := false
	if acct.AccessToken != "" && acct.Login != "" {

	}

	if !authd {
		if pwd == "" {
			pwd = console.ReadPassword("Please enter mojang account password: ")
		}
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
		err := mcauth.Login(mcrequest.NewMinecraftLogin(acct.Login, pwd, acct.ClientToken), &acct)
		if err != nil {
			console.Println(mcchat.Red, "Error while logging into Minecraft: ", err)
		}else {
			authd = true
			cache.Store(acct)
		}

		if authd {
			console.Println("Successfully authenticated!")
		}
	}

	go console.Scan(func(input string) {
		if input == "/quit" {
			exit <- true
		}
	})

	<-exit
}



