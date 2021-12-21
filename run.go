package hotomata

import (
	"errors"
	"io/ioutil"
	"path"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

const planFileExt = ".yaml"

// Task represents a command to be ran on a remote machine and all the variables
// that represent it's context. Those variables come from the inventory file,
// the var files, the masterplans, all the plans util a $run was found
type Task struct {
	TopLevelName string
	Name         string
	Run          string
	SpecialFlags PlanSpecialFlags
	VarsChain    []PlanVars
}

// Run represents the context needed to run commands against machines. It hold
// all the plans discovered and the inventory of machines and has methods to
// either do remote execution of single commands or execution of a complete
// masterplan
type Run struct {
	plans     map[string]*Plan
	inventory []InventoryMachine
}

// NewRun creates an empty Run
func NewRun() *Run {
	return &Run{
		plans:     map[string]*Plan{},
		inventory: []InventoryMachine{},
	}
}

// DiscoverPlans searches recursively a directory for plan files and parses them
func (r *Run) DiscoverPlans(directory string) error {
	var loadFolder func(string) error
	loadFolder = func(folder string) error {
		folders, err := ioutil.ReadDir(folder)
		if err != nil {
			return err
		}
		for _, f := range folders {
			if f.IsDir() {
				err = loadFolder(path.Join(folder, f.Name()))
				if err != nil {
					return err
				}
				continue
			} else if !strings.HasSuffix(f.Name(), planFileExt) {
				continue
			}

			// Ok, at this point we got a .yaml file to load
			contents, err := ioutil.ReadFile(path.Join(folder, f.Name()))
			if err != nil {
				return err
			}

			planName := strings.TrimSuffix(f.Name(), planFileExt)
			plan, err := ParsePlan(planName, contents)
			if err != nil {
				return err
			}

			r.plans[planName] = plan
		}

		return nil
	}

	return loadFolder(directory)
}

// Plan fetches a plan from a Run's context
func (r *Run) Plan(name string) (*Plan, bool) {
	plan, ok := r.plans[name]
	return plan, ok
}

// Plans return all plans discovered to date
func (r *Run) Plans() map[string]*Plan {
	return r.plans
}

// LoadInventory appends a list of inventory machines to the current list of machines
func (r *Run) LoadInventory(machines []InventoryMachine) {
	r.inventory = append(r.inventory, machines...)
}

// FilterGroups removes machines not in the given group
func (r *Run) FilterGroups(filter string) {
	if filter == "*" {
		return
	}
	var newInventory = []InventoryMachine{}

	for _, m := range r.inventory {
		for _, g := range m.Groups {
			if g.GroupName == filter {
				newInventory = append(newInventory, m)
			}
		}
	}

	r.inventory = newInventory
}

// RunMasterPlans runs a set of masterplans
func (r *Run) RunMasterPlans(logger *Logger, masterplans []*MasterPlan) (*RunReport, error) {
	report := &RunReport{}

	for _, masterplan := range masterplans {
		err := r.RunMasterPlan(logger, report, masterplan)
		if err != nil {
			return report, err
		}
	}

	return report, nil
}

// RunMasterPlan runs a specific part of the masterplan
func (r *Run) RunMasterPlan(logger *Logger, report *RunReport, masterplan *MasterPlan) error {
	machines := masterplan.FilterMachines(r.inventory)
	runnerMachines := MachinesFromInventoryMachines(machines)

	logRunStart(logger, machines)

	// Convert plain plan names to PlanCalls
	var topPlanCalls []*PlanCall
	for _, planName := range masterplan.Plans {
		topPlanCalls = append(topPlanCalls, &PlanCall{
			Name: planName,
			Plan: planName,
			Vars: PlanVars{},
		})
	}

	// Build plan tree, dereferencing all sub plans
	tasks, err := r.dereferenceTasksFromPlanCalls(
		PlanSpecialFlags{},
		[]PlanVars{masterplan.Vars},
		topPlanCalls,
		"",
	)
	if err != nil {
		logger.Write(ColorRed, "abort: "+err.Error()+"\n")
		return err
	}

	for _, task := range tasks {
		logger.WriteLine(ColorCyan, "TASK: [ %s ] %s ", task.TopLevelName, task.Name)
		for i, m := range machines {
			spew.Dump(task)
			cmd, err := ExecuteTemplate(task.Run, append([]PlanVars{m.PlanVars()}, task.VarsChain...))
			if err != nil {
				logger.WriteLine(ColorRed, "abort: [ %s ] %s", m.Name, "Failed to compile template")
				logger.Write(ColorRed, err.Error()+"\n")
				logger.Writenc(task.Run + "\n")
				return err
			} else {
				// run the action!
				runner := SSHRunner{}

				// handle local and sudo special flags
				machine := *runnerMachines[i]
				if task.SpecialFlags.Local {
					machine.Hostname = "127.0.0.1"
					machine.Port = 0
				}
				if task.SpecialFlags.Sudo {
					cmd = "sudo " + cmd
				}

				// run and mark as skipped if ignore_errors=true and we errored
				response := runner.Run(machine, cmd)
				if task.SpecialFlags.IgnoreErrors &&
					response.Status == TaskStatusError &&
					response.Action == TaskActionAbort {
					response.Action = TaskActionContinue
					response.Status = TaskStatusSkip
				}

				// log result, handle verbose flag
				logTaskResponse(logger, task, m, response)
				logger.Writenc(cmd + "\n")
				if logger.Verbose {
					logger.Writenc(response.Log.String() + "\n")
				}

				// abort if we errored
				if response.Action == TaskActionAbort {
					return errors.New("aborted while processing task " + task.Name + " on host " + m.Name)
				}
			}
		}
	}

	return nil
}

// dereferenceTasksFromPlanCalls is a recursive function that extracts run commands
// and transforms them to tasks based on the context
func (r *Run) dereferenceTasksFromPlanCalls(
	specialFlags PlanSpecialFlags,
	varsChain []PlanVars,
	planCalls []*PlanCall,
	topLevelName string,
) ([]*Task, error) {
	var tasks = []*Task{}

	for _, pc := range planCalls {
		tln := topLevelName
		if tln == "" {
			tln = pc.Name
		}

		if pc.Run != "" {
			// Gather vars and create task
			tasks = append(tasks, &Task{
				TopLevelName: tln,
				Name:         pc.Name,
				Run:          pc.Run,
				SpecialFlags: specialFlags.Join(pc),
				VarsChain:    append(varsChain, pc.Vars),
			})
			spew.Dump(tasks[len(tasks)-1])
		} else {
			// Go deeper
			if plan, ok := r.Plan(pc.Plan); ok {
				newTasks, err := r.dereferenceTasksFromPlanCalls(
					specialFlags.Join(pc),
					append(varsChain, plan.Vars, pc.Vars),
					plan.PlanCalls,
					tln,
				)
				if err != nil {
					return tasks, err
				}
				tasks = append(tasks, newTasks...)
			} else {
				return tasks, errors.New("Plan " + pc.Plan + " is missing")
			}
		}
	}

	return tasks, nil
}

func logRunStart(logger *Logger, machines []InventoryMachine) {
	var runMachineNames string
	for _, m := range machines {
		runMachineNames = runMachineNames + m.Name + " "
	}
	logger.Writenc("\n")
	logger.WriteLine(ColorMagenta, "RUN: [ %s]", runMachineNames)
	logger.Writenc("\n")
}

func logTaskResponse(logger *Logger, task *Task, m InventoryMachine, r *TaskResponse) {
	logger.WriteLine(r.Color(), "%s(%s) [ %s ] local=%t sudo=%t ignore_errors=%t", r.Action, r.Status, m.Name, task.SpecialFlags.Local, task.SpecialFlags.Sudo, task.SpecialFlags.IgnoreErrors)
}
