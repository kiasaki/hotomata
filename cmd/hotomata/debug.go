package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/kiasaki/hotomata"
)

func setupDebug(c *cli.Context) *hotomata.Run {
	var err error
	var cwd string

	if cwd, err = os.Getwd(); err != nil {
		panic(err)
	}

	// Discover plans
	run := hotomata.NewRun()
	if err = run.DiscoverPlans("/etc/hotomata/plans"); err != nil {
		writeError("Error: Unable to load core plans", err)
	}
	if err = run.DiscoverPlans(path.Join(cwd, "plans")); err != nil {
		writeError("Error: Unable to load plans", err)
	}

	return run
}

func debugPlan(c *cli.Context) {
	run := setupDebug(c)

	// Parse plan args
	var planName = c.Args().First()
	if planName == "" {
		writeError("Error: A plan is required. e.g. `hotomata debug plan redis`", nil)
	}

	// Fetch concerned plan
	plan, ok := run.Plan(planName)
	if !ok {
		writeError("Error: Unable to find plan '"+planName+"'", nil)
	}

	writePlan("", run, plan, true)
}

func debugPlans(c *cli.Context) {
	run := setupDebug(c)

	fmt.Println("All plans")

	for _, p := range run.Plans() {
		writePlan("", run, p, false)
		fmt.Println("")
	}
}

func writePlan(in string, run *hotomata.Run, p *hotomata.Plan, detailed bool) {
	// Bump indentation each level
	in = in + "  "

	if detailed {
		var vars string
		for k, v := range p.Vars {
			vars = color(hotomata.ColorCyan, fmt.Sprintf("%s %s=%v", vars, k, v))
		}
		writef(hotomata.ColorMagenta, "%s%s%s", in, p.Name, vars)
	} else {
		writef(hotomata.ColorMagenta, "%s%s", in, p.Name)
	}

	for _, planCall := range p.PlanCalls {
		if detailed {
			fmt.Printf("%s  %s\n", in, planCall.Name)
		}
		if planCall.Run != "" {
			var extra string
			if planCall.Local {
				extra = extra + color(hotomata.ColorCyan, " local=true")
			}
			if planCall.Sudo {
				extra = extra + color(hotomata.ColorCyan, " sudo=true")
			}
			if planCall.IgnoreErrors {
				extra = extra + color(hotomata.ColorCyan, " ignore_errors=true")
			}

			writef(hotomata.ColorGreen, "%s  %s%s", in, strings.Split(planCall.Run, "\n")[0], extra)
		} else {
			plan, found := run.Plan(planCall.Plan)
			if found {
				writePlan(in, run, plan, detailed)
			} else {
				if detailed {
					writef(hotomata.ColorRed, "%s  %s (plan missing)", in, planCall.Plan)
				} else {
					writef(hotomata.ColorRed, "%s  %s (missing)", in, planCall.Plan)
				}
			}
		}
	}
}
