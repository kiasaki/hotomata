package main

import (
	"io/ioutil"
	"os"

	"github.com/codegangsta/cli"
	"github.com/kiasaki/hotomata"
)

func run(c *cli.Context) {
	var contents []byte
	var err error

	// Parse inventory arg
	var inventoryFile = c.GlobalString("inventory")
	if inventoryFile == "" {
		writeError("Error: An inventory file is required. e.g. `hotomata --inventory inventory.json`", nil)
	}

	// Parse actual inventory
	contents, err = ioutil.ReadFile(inventoryFile)
	if err != nil {
		writeError("Error: Unable to read inventory file at "+inventoryFile, err)
	}
	inventory, err := hotomata.ParseInventory(contents)
	if err != nil {
		writeError("Error: Unable to parse inventory file, verify your JSON syntax", err)
	}

	// Parse masterplan arg
	var masterPlanFile = c.GlobalString("masterplan")
	if c.Args().First() != "" {
		masterPlanFile = c.Args().First()
	}

	// Parse actual masterplan
	contents, err = ioutil.ReadFile(masterPlanFile)
	if err != nil {
		writeError("Error: Unable to read masterplan file at "+masterPlanFile, err)
	}
	masterplans, err := hotomata.ParseMasterPlan(contents)
	if err != nil {
		writeError("Error: Unable to parse masterplan file, verify your YAML syntax", err)
	}

	// Create a run and parse plans
	run := hotomata.NewRun()
	err = run.DiscoverPlans(c.GlobalString("core-plans-folder"))
	if err != nil {
		writeError("Error: could not load core plans folder at "+c.GlobalString("core-plans-folder"), err)
	}
	err = run.DiscoverPlans(c.GlobalString("plans-folder"))
	if err != nil {
		writeError("Error: could not load plans folder at "+c.GlobalString("plans-folder"), err)
	}

	// load inventory and limit groups
	run.LoadInventory(inventory)
	run.FilterGroups(c.String("group"))

	logger := hotomata.NewLogger(os.Stderr, c.GlobalString("color") == "true", c.GlobalBool("verbose"))
	run.RunMasterPlans(logger, masterplans)
}
