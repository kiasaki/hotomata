package hotomata

import (
	"strings"

	"github.com/flosch/pongo2"
)

// ExecuteTemplate uses the pongo template engine to render a template file
// contents or a command string
func ExecuteTemplate(template string, varsChain []PlanVars) (string, error) {
	vars, err := flattenVarsChain(varsChain)
	if err != nil {
		return "", err
	}
	out, err := executeTemplate(template, vars)
	return strings.Trim(out, " \n\r"), err
}

func executeTemplate(template string, vars PlanVars) (string, error) {
	tmpl, err := pongo2.FromString(template)
	if err != nil {
		return "", err
	}

	return tmpl.Execute(pongo2.Context(vars))
}

func flattenVarsChain(varsChain []PlanVars) (PlanVars, error) {
	var vars = PlanVars{}

	for _, planVars := range varsChain {
		if len(planVars) == 0 {
			continue
		}

		for k, v := range planVars {
			if v == nil {
				continue
			}

			if value, ok := v.(string); ok && strings.Contains(value, "{") {
				out, err := executeTemplate(value, vars)
				if err != nil {
					return vars, err
				}
				vars[k] = out
			} else {
				vars[k] = v
			}
		}
	}

	return vars, nil
}
