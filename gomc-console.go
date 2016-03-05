package main

import (
	"fmt"

	"github.com/initsuj/gomc-console/console"
	"github.com/initsuj/gomc/mcchat"
	"flag"
	"github.com/initsuj/gomc-console/cache"
)

func main() {

	var ini, server, user, pwd string

	flag.StringVar(&ini, "i", "gomc.ini", "location of ini file to use.")
	flag.StringVar(&server, "s", "", "server to connect to.")
	flag.StringVar(&user, "u", "", "mojang account username.")
	flag.StringVar(&pwd, "p", "", "mojang account password.")
	flag.Parse()

	if err := cache.Init(); err != nil{
		panic(err)
	}
	defer cache.Close()



	//fmt.Print("test = " + mcchat.Black)
	console.Print(mcchat.Yellow, "hello ", mcchat.Red, "Red, ", mcchat.DarkRed, " DarkRed")
	console.Println(mcchat.Yellow, "hello ", mcchat.Red, "Red, ", mcchat.DarkRed, " DarkRed")
	console.Println("test")
	console.Println("test §1Blue text §fand §cRed text")
	fmt.Println("hello")
}
