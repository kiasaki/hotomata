package hotomata

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

type PlanVars map[string]interface{}

type PlanSpecialFlags struct {
	Local        bool
	Sudo         bool
	IgnoreErrors bool
}

func (psf PlanSpecialFlags) Join(pc *PlanCall) PlanSpecialFlags {
	if pc.Local {
		psf.Local = true
	}
	if pc.Sudo {
		psf.Sudo = true
	}
	if pc.Local {
		psf.IgnoreErrors = true
	}
	return psf
}

type PlanCall struct {
	Name         string
	Run          string
	Plan         string
	Local        bool
	Sudo         bool
	IgnoreErrors bool
	Vars         PlanVars
}

type Plan struct {
	Name      string
	Vars      PlanVars
	PlanCalls []*PlanCall
}

func ParsePlan(planName string, yamlSource []byte) (*Plan, error) {
	var plan = &Plan{Name: planName}

	// Unmarshal raw yaml
	var rawPlan struct {
		Vars  map[string]interface{}
		Plans []map[string]interface{}
	}
	err := yaml.Unmarshal(yamlSource, &rawPlan)
	if err != nil {
		return plan, errors.New(fmt.Sprintf("Error parsing yaml for plan: %s (%s)", planName, err.Error()))
	}

	// Fill structs that are nicer to work with
	plan.Vars = rawPlan.Vars
	for _, rawPlanCall := range rawPlan.Plans {
		planCall := &PlanCall{}

		// Parse PlanCall $name
		if rawName, ok := rawPlanCall["$name"]; ok {
			if name, ok := rawName.(string); ok {
				planCall.Name = name
			} else {
				return plan, newError("Error parsing plan: %s: $name is not a string", planName)
			}
		} else {
			return plan, newError("Error parsing plan: %s: $name is required", planName)
		}

		// Parse PlanCall $run
		if rawRun, ok := rawPlanCall["$run"]; ok {
			if run, ok := rawRun.(string); ok {
				planCall.Run = run
			} else {
				return plan, newError("Error parsing plan: %s: $run is not a string (%s)", planName, planCall.Name)
			}
		}

		// Parse PlanCall $plan
		if rawPlanCallPlan, ok := rawPlanCall["$plan"]; ok {
			if planCallPlan, ok := rawPlanCallPlan.(string); ok {
				planCall.Plan = planCallPlan
			} else {
				return plan, newError("Error parsing plan: %s: $plan is not a string (%s)", planName, planCall.Name)
			}
		}

		var boolValue bool
		boolValue, err = getRawPlanCallBool(rawPlanCall, "$local")
		planCall.Local = boolValue
		if err != nil {
			return plan, newError("Error parsing plan: %s: $local is not 'true' or 'false' (%s)", planName, planCall.Name)
		}
		boolValue, err = getRawPlanCallBool(rawPlanCall, "$sudo")
		planCall.Sudo = boolValue
		if err != nil {
			return plan, newError("Error parsing plan: %s: $sudo is not 'true' or 'false' (%s)", planName, planCall.Name)
		}
		boolValue, err = getRawPlanCallBool(rawPlanCall, "$ignore_errors")
		planCall.IgnoreErrors = boolValue
		if err != nil {
			return plan, newError("Error parsing plan: %s: $ignore_errors is not 'true' or 'false' (%s)", planName, planCall.Name)
		}

		// Verify we have an action to do
		if planCall.Run == "" && planCall.Plan == "" {
			return plan, newError("Error parsing plan: %s: $run or $plan is required (%s)", planName, planCall.Name)
		}

		delete(rawPlanCall, "$name")
		delete(rawPlanCall, "$run")
		delete(rawPlanCall, "$plan")
		delete(rawPlanCall, "$local")
		delete(rawPlanCall, "$sudo")
		delete(rawPlanCall, "$ignore_errors")
		planCall.Vars = PlanVars(rawPlanCall)

		plan.PlanCalls = append(plan.PlanCalls, planCall)
	}

	return plan, nil
}

func getRawPlanCallBool(rawPlanCall map[string]interface{}, key string) (bool, error) {
	if raw, ok := rawPlanCall[key]; ok {
		if value, ok := raw.(bool); ok {
			return value, nil
		} else {
			return false, errors.New("Not a bool")
		}
	}
	return false, nil
}
