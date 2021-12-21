package hotomata

import (
	"encoding/json"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// MasterPlan represents a set of machine filters and plans to be applied
// to matched machines at runtime
type MasterPlan struct {
	MachineFilters []*MachineFilter
	Vars           PlanVars
	Plans          []string
}

// FilterMachines takes whole inventory of machines and returns only the machines
// that were matched by the MasterPlan's filters
func (mp *MasterPlan) FilterMachines(inventory []InventoryMachine) []InventoryMachine {
	var machines = []InventoryMachine{}
	for _, machine := range inventory {
		allMatch := true
		// Pass through all filters, if one doesn't match it is enough to
		// not retain the current machine
		for _, filter := range mp.MachineFilters {
			match := filter.MatchesMachine(machine)
			if !match {
				allMatch = false
			}
		}
		if allMatch {
			machines = append(machines, machine)
		}
	}
	return machines
}

// MachineFilter represents a machine param to filter on and the patter it must
// match
type MachineFilter struct {
	Param   string
	Pattern string
}

func (mf *MachineFilter) Regexp() *regexp.Regexp {
	if mf.Pattern[0] == '/' && mf.Pattern[len(mf.Pattern)] == '/' {
		return regexp.MustCompile(mf.Pattern[1 : len(mf.Pattern)-1])
	} else {
		// Make sure we escape those meta characters
		fragments := strings.Split(mf.Pattern, "*")
		for i, f := range fragments {
			fragments[i] = regexp.QuoteMeta(f)
		}
		return regexp.MustCompile(strings.Join(fragments, ".+"))
	}
}

// MatchesMachine checks if the current MachineFilter matches the given machine's
// vars
func (mf *MachineFilter) MatchesMachine(machine InventoryMachine) bool {
	machineVars := machine.Vars()
	matcher := mf.Regexp()

	if rawValue, ok := machineVars[mf.Param]; ok {
		var stringValue string
		if err := json.Unmarshal(rawValue, &stringValue); err == nil {
			return matcher.MatchString(stringValue)
		}
		// TODO(kiasaki) Support arrays of string in inventory params
	}

	return false
}

// ParseMasterPlan takes a raw yaml file and parses it into a set of MasterPlans
func ParseMasterPlan(yamlSource []byte) ([]*MasterPlan, error) {
	var plans = []*MasterPlan{}

	// Unmarshal raw yaml
	var rawPlans []struct {
		Machines map[string]string
		Vars     map[string]interface{}
		Plans    []string
	}
	err := yaml.Unmarshal(yamlSource, &rawPlans)
	if err != nil {
		return plans, err
	}

	// Fill structs that are nicer to work with
	for _, rawPlan := range rawPlans {
		plan := &MasterPlan{MachineFilters: []*MachineFilter{}}

		for k, v := range rawPlan.Machines {
			plan.MachineFilters = append(plan.MachineFilters, &MachineFilter{
				Param:   k,
				Pattern: v,
			})
		}

		plan.Vars = rawPlan.Vars

		plan.Plans = rawPlan.Plans

		plans = append(plans, plan)
	}

	return plans, nil
}
