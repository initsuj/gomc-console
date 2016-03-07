package console

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/daviddengcn/go-colortext"
	"github.com/initsuj/gomc/mcchat"
	"os"
	"bufio"
	"github.com/howeyc/gopass"
)

var (
	colorStr  = regexp.MustCompile("(§{1}[0-9a-fA-F])")
	str2Color = make(map[string]mcchat.Colorer, 16)
	stopScan = make(chan bool, 1)
	scanning = false
)

type InputHandler(string)

func init() {
	str2Color[string(mcchat.Black)] = mcchat.Black
	str2Color[string(mcchat.DarkBlue)] = mcchat.DarkBlue
	str2Color[string(mcchat.DarkGreen)] = mcchat.DarkGreen
	str2Color[string(mcchat.DarkAqua)] = mcchat.DarkAqua
	str2Color[string(mcchat.DarkRed)] = mcchat.DarkRed
	str2Color[string(mcchat.DarkPurple)] = mcchat.DarkPurple
	str2Color[string(mcchat.Gold)] = mcchat.Gold
	str2Color[string(mcchat.Gray)] = mcchat.Gray
	str2Color[string(mcchat.DarkGray)] = mcchat.DarkGray
	str2Color[string(mcchat.Blue)] = mcchat.Blue
	str2Color[string(mcchat.Green)] = mcchat.Green
	str2Color[string(mcchat.Aqua)] = mcchat.Aqua
	str2Color[string(mcchat.Red)] = mcchat.Red
	str2Color[string(mcchat.Purple)] = mcchat.Purple
	str2Color[string(mcchat.Yellow)] = mcchat.Yellow
	str2Color[string(mcchat.White)] = mcchat.White

}

func Print(args ...interface{}) {
	for _, a := range args {
		doPrint(a)
	}

	ct.ResetColor()
}

func Println(args ...interface{}) {
	for _, a := range args {
		doPrint(a)
	}

	doPrint("\r\n")
	ct.ResetColor()
}

func doPrint(args ...interface{}) {
	for _, a := range args {

		if c, ok := a.(mcchat.Colorer); ok {
			SetForegroundColor(c)
			continue
		}

		if s, ok := a.(string); ok {
			cs := colorStr.FindAllString(s, -1)
			for _, str := range cs {
				//fmt.Println(str)
				i := strings.Index(s, str)
				if i != 0 {
					fmt.Print(s[:i])
				}
				if c, ok := str2Color[str]; ok {
					SetForegroundColor(c)
				}

				s = s[i+3:]
			}

			if len(s) > 0 {
				fmt.Print(s)
			}
			ct.ResetColor()
			continue
		}

		fmt.Print(a)
	}
}

func SetForegroundColor(c mcchat.Colorer) {
	switch c.Color() {
	case mcchat.Black:
		ct.Foreground(ct.Black, false)
	case mcchat.DarkBlue:
		ct.Foreground(ct.Blue, false)
	case mcchat.DarkGreen:
		ct.Foreground(ct.Green, false)
	case mcchat.DarkAqua:
		ct.Foreground(ct.Cyan, false)
	case mcchat.DarkRed:
		ct.Foreground(ct.Red, false)
	case mcchat.DarkPurple:
		ct.Foreground(ct.Magenta, false)
	case mcchat.Gold:
		ct.Foreground(ct.Yellow, false)
	case mcchat.DarkGray:
		ct.Foreground(ct.Black, true)
	case mcchat.Gray:
		ct.Foreground(ct.White, false)
	case mcchat.Blue:
		ct.Foreground(ct.Blue, true)
	case mcchat.Green:
		ct.Foreground(ct.Green, true)
	case mcchat.Aqua:
		ct.Foreground(ct.Cyan, true)
	case mcchat.Red:
		ct.Foreground(ct.Red, true)
	case mcchat.Purple:
		ct.Foreground(ct.Magenta, true)
	case mcchat.Yellow:
		ct.Foreground(ct.Yellow, true)
	case mcchat.White:
		ct.Foreground(ct.White, true)

	}
}

func Prompt(m string) (v string){
	for {
		Print(m)
		_, err := fmt.Scanln(&v)
		if err != nil{
			Print(mcchat.Red, "There was an error reading your input!")
			continue
		}
		if v != ""{
			break
		}
	}
	return
}

func ReadPassword(m string)(v string){
	for {
		Print(m)
		p, err := gopass.GetPasswdMasked()
		if err != nil{
			Print(mcchat.Red, "There was an error reading your password!")
		}
		if v = string(p); v != ""{
			break
		}
	}

	return
}

func Scan(handler func(string)){
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		go func(){handler(input.Text())}()
	}
}

func Close(){
	if scanning{
		stopScan <- true
	}
}