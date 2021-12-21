package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/codegangsta/cli"
	"github.com/kiasaki/hotomata"
)

func main() {
	app := cli.NewApp()
	app.Name = "hotomata-inventory"
	app.Usage = "tool to check validity of inventory files and introspect their contents as seen by hotomata"
	app.Version = "0.1.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "kiasaki",
			Email: "kiasaki0000@gmail.com",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "check",
			Aliases: []string{"c"},
			Usage:   "Verifies a given inventory file is valid",
			Action:  checkCmd,
		},
		{
			Name:    "print",
			Aliases: []string{"p"},
			Usage:   "Prints the contents of an inventory file",
			Action:  printCmd,
		},
	}

	app.Run(os.Args)
}

func checkCmd(c *cli.Context) {
	contents, err := ioutil.ReadFile(c.Args().First())
	if err != nil {
		panic(err)
	}

	result, err := hotomata.ValidateInventory(string(contents))
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf(hotomata.Colorize("The document is valid\n", hotomata.ColorGreen))
	} else {
		fmt.Printf(hotomata.Colorize("The document is not valid. see errors :\n", hotomata.ColorRed))
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
}

func printCmd(c *cli.Context) {
	contents, err := ioutil.ReadFile(c.Args().First())
	if err != nil {
		fmt.Printf(hotomata.Colorize("Can't read file: %s", hotomata.ColorRed), c.Args().First())
	}

	machines, err := hotomata.ParseInventory(contents)
	if err != nil {
		fmt.Printf(hotomata.Colorize("%s\n", hotomata.ColorRed), err.Error())
	} else {
		for _, machine := range machines {
			fmt.Printf(hotomata.Colorize("Machine: %s\n", hotomata.ColorMagenta), machine.Name)
			fmt.Printf(hotomata.Colorize("Groups: %v\n", hotomata.ColorBlue), machine.Groups.Names())
			fmt.Print("Variables:\n")

			// Here's the tricky part, lets sort them alphabeticlay
			var pairs = PropertyPairs{}
			for k, v := range machine.Vars() {
				pairs = append(pairs, PropertyPair{k, v})
			}

			sort.Sort(&pairs)
			for _, pair := range pairs {
				valString, err := pair.Value.MarshalJSON()
				if err != nil {
					panic(err) // should never happen it's unmarshaled json
				}
				fmt.Printf("    %s: %s\n", pair.Property, string(valString))
			}

			fmt.Println("")
		}
	}
}

type PropertyPair struct {
	Property string
	Value    json.RawMessage
}
type PropertyPairs []PropertyPair

func (p PropertyPairs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PropertyPairs) Len() int           { return len(p) }
func (p PropertyPairs) Less(i, j int) bool { return p[i].Property > p[j].Property }
