package unregisteredagents

import (
	"fmt"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

func init() {
	builder := checks.CheckFuncBuilder{CheckMasters: collect}

	check := checks.Check{
		Name:           "mesos-unregistered-agents",
		Description:    "Checks for unregistered Mesos agents",
		Cure:           "Determine why Mesos agents cannot register by examining the Mesos agent log.",
		OKSummary:      "All Mesos agents appear to be registered.",
		ProblemSummary: "Some Mesos agents appear to be unregistered.",
		Run:            builder.Build(),
	}

	checks.RegisterCheck(check)
}

type slave struct {
	PID string `json:"pid"`
}

type state struct {
	RecoveredSlaves []slave `json:"recovered_slaves"`
}

func collect(host bundle.Host) checks.Result {
	var state state

	if err := host.ReadJSON("mesos-master-state", &state); err != nil {
		return checks.Result{
			Status: checks.SUndefined,
			Value:  err,
		}
	}

	var unregistered []string
	for _, slave := range state.RecoveredSlaves {
		unregistered = append(unregistered, fmt.Sprintf("(Mesos) %v appears unregistered", slave.PID))
	}

	if len(unregistered) > 0 {
		return checks.Result{
			Host:   host,
			Status: checks.SProblem,
			Value:  unregistered,
		}
	}

	return checks.Result{
		Status: checks.SOK,
	}
}
