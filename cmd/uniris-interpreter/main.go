package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/uniris/uniris-core/pkg/interpreter"

	uniris "github.com/uniris/uniris-interpreter/pkg"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "uniris-interpreter"
	app.Usage = "Interpreter for UNIRIS smart contract"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Execute code from a `FILE` source code",
		},
		cli.BoolFlag{
			Name:  "console, c",
			Usage: "Open console to execute code instantly",
		},
	}

	app.Action = func(c *cli.Context) error {

		if c.String("file") != "" {
			code, err := ioutil.ReadFile(c.String("file"))
			if err != nil {
				return err
			}
			res, err := uniris.Interpret(string(code), nil)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			fmt.Print(res)
			return nil
		} else if c.Bool("console") {
			fmt.Println("Type Ctrl-C to exit the console")
			scope := interpreter.NewScope(nil)
			for {
				text := read()
				res, err := interpreter.Execute(text, scope)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
				} else {
					fmt.Print(res)
				}
			}
		}

		return cli.ShowAppHelp(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func read() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	text, _ := reader.ReadString('\n')
	return strings.Trim(text, "")
}

func handleSystemCommand(text string) string {
	switch text {
	case "help":
		//Show help
		return ""
	default:
		return text
	}
}
