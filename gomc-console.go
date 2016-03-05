package main

import (
	"fmt"

	"github.com/initsuj/gomc-console/console"
	"github.com/initsuj/gomc/mcchat"
)

func main() {
	//fmt.Print("test = " + mcchat.Black)
	console.Print(mcchat.Yellow, "hello ", mcchat.Red, "Red, ", mcchat.DarkRed, " DarkRed")
	console.Println(mcchat.Yellow, "hello ", mcchat.Red, "Red, ", mcchat.DarkRed, " DarkRed")
	console.Println("test")
	console.Println("test §1Blue text §fand §cRed text")
	fmt.Println("hello")
}
