package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/kiasaki/hotomata"
)

func main() {
	app := cli.NewApp()
	app.Name = "hotomata"
	app.Usage = "tool to execute masterplan scripts againt an inventory via ssh"
	app.Version = "0.1.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "kiasaki",
			Email: "kiasaki0000@gmail.com",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "inventory, i",
			Value:  "inventory.json",
			Usage:  "inventory file location",
			EnvVar: "HOTOMATA_INVENTORY_FILE",
		},
		cli.StringFlag{
			Name:   "masterplan, m",
			Value:  "masterplan.yaml",
			Usage:  "masterplan file location",
			EnvVar: "HOTOMATA_MASTERPLAN_FILE",
		},
		cli.StringFlag{
			Name:   "plans-folder, f",
			Value:  "plans",
			Usage:  "plans folder location",
			EnvVar: "HOTOMATA_PLANS_FOLDER",
		},
		cli.StringFlag{
			Name:   "core-plans-folder",
			Value:  "/etc/hotomata/plans",
			Usage:  "core plans folder location",
			EnvVar: "HOTOMATA_CORE_PLANS_FOLDER",
		},
		cli.StringFlag{
			Name:   "color",
			Value:  "true",
			Usage:  "enable colored output (true or false)",
			EnvVar: "HOTOMATA_COLOR",
		},
		cli.BoolFlag{
			Name:   "verbose, V",
			Usage:  "Log more info",
			EnvVar: "HOTOMATA_VERBOSE",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Runs a given masterplan against an inventory of machines",
			Action:  run,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "group, g",
					Value: "*",
					Usage: "Limit to certain inventory groups '*' for all",
				},
			},
		},
		{
			Name:  "debug",
			Usage: "Few sub commands for debuging your plans",
			Subcommands: []cli.Command{
				{
					Name:   "plan",
					Usage:  "Visualise a specific plan",
					Action: debugPlan,
				},
				{
					Name:   "plans",
					Usage:  "Visualise all plans discovered",
					Action: debugPlans,
				},
			},
		},
	}

	app.Run(os.Args)
}

func writef(c hotomata.Color, message string, params ...interface{}) {
	fmt.Printf(color(c, message)+"\n", params...)
}

func writeError(message string, err error) {
	var completeMessage = message
	if err != nil {
		completeMessage = fmt.Sprintf("%s (%s)", completeMessage, err.Error())
	}
	fmt.Print(color(hotomata.ColorRed, completeMessage+"\n"))
	os.Exit(1)
}

func color(color hotomata.Color, message string) string {
	return hotomata.Colorize(message, color)
}
