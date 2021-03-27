package instances

import (
	"fmt"
	"strings"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

func init() {
	check := checks.Check{
		Name:        "marathon-app-instances",
		Description: "Checks if Marathon apps have required amount of instances",
		Cure: "Applications couldn't start for multiple reasons. Check Marathon logs for the NOT SATISFIED " +
			"pattern to make sure that there is enough resources to start tasks." +
			"Check if tasks are failing tasks and read through their logs.",
		OKSummary:      "Some Marathon apps have less or more instances than required.",
		ProblemSummary: "All Marathon tasks have the required amount of instances.",
		Run:            checkFunc,
	}
	checks.RegisterCheck(check)
}

type apps struct {
	Apps []struct {
		Id           string
		Instances    int
		TasksRunning int
	}
}

func checkFunc(b bundle.Bundle) checks.Results {
	var apps apps
	err := b.ReadAnyJSON("marathon-apps", &apps)
	if err != nil {
		return checks.Results{checks.Result{
			Status: checks.SUndefined,
			Value:  "Couldn't read any marathon apps JSONs",
		}}
	}
	problems := make([]string, 0, len(apps.Apps))
	for _, app := range apps.Apps {
		if app.TasksRunning < app.Instances {
			problems = append(problems, fmt.Sprintf("App %s has less instances than required: %v < %v", app.Id, app.TasksRunning, app.Instances))
			continue
		}
		if app.TasksRunning > app.Instances {
			problems = append(problems, fmt.Sprintf("App %s has more instances than required: %v > %v", app.Id, app.TasksRunning, app.Instances))
		}
	}
	if len(problems) > 0 {
		return checks.Results{
			checks.Result{
				Status: checks.SProblem,
				Value:  strings.Join(problems, "\n"),
			},
		}
	}
	return checks.Results{
		checks.Result{
			Status: checks.SOK,
		},
	}
}
