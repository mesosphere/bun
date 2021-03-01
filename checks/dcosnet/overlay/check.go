package overlay

import (
	"fmt"
	"strings"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

func init() {
	builder := checks.CheckFuncBuilder{
		CheckAgents:       collect,
		CheckPublicAgents: collect,
		Aggregate:         checks.DefaultAggregate,
	}

	check := checks.Check{
		Name:        "overlay",
		Description: "Checks if all the overlay network interfaces created successfully on all the agents",
		Cure: "Docker already has the listed networks? Try to delete them with the `docker network rm` " +
			"command and restart the dcos-net. See https://jira.d2iq.com/browse/COPS-4789",
		OKSummary:      "dcos-net created all the required overlay network interfaces on all the agents",
		ProblemSummary: "dcos-net could not create some overlay network interfaces on some agents",
		Run:            builder.Build(),
	}
	checks.RegisterCheck(check)
}

func collect(host bundle.Host) checks.Result {
	var o overlay
	err := host.ReadJSON("mesos-agent-overlay", &o)
	if err != nil {
		return checks.Result{
			Status: checks.SUndefined,
			Value:  "Couldn't read Mesos agent overlay state: " + err.Error(),
		}
	}
	problems := make([]string, 0, len(o.Overlays))
	for _, network := range o.Overlays {
		if network.State.Error != "" {
			problem := fmt.Sprintf("Network: %s, error message: %s", network.Info.Name, network.State.Error)
			problems = append(problems, problem)
		}
	}
	if len(problems) > 0 {
		return checks.Result{
			Status: checks.SProblem,
			Value:  strings.Join(problems, "\n"),
		}
	}
	return checks.Result{Status: checks.SOK}
}
