package unregisteredagents

import (
	"fmt"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"

	"github.com/lithammer/dedent"
)

func init() {
	builder := checks.CheckFuncBuilder{CheckMasters: collect}

	check := checks.Check{
		Name:        "mesos-unregistered-agents",
		Description: "Checks for unregistered Mesos agents",
		Cure: dedent.Dedent(`
			Determine why Mesos agents cannot register by examining the Mesos agent log.
			In order to register an agent needs to have finished its revovery, have detected the master, and be able to connect to it.`),
		OKSummary:      "All Mesos agents appear to be registered.",
		ProblemSummary: "Some Mesos agents appear to be unregistered.",
		Run:            builder.Build(),
	}

	checks.RegisterCheck(check)
}

type slave struct {
	PID string `json:"pid"`
}

type leaderinfo struct {
	ID string `json:"id"`
}

type state struct {
	RecoveredSlaves []slave    `json:"recovered_slaves"`
	ID              string     `json:"id"`
	LeaderInfo      leaderinfo `json:"leader_info"`
}

func collect(host bundle.Host) checks.Result {
	var state state

	if err := host.ReadJSON("mesos-master-state", &state); err != nil {
		return checks.Result{
			Status: checks.SUndefined,
			Value:  err,
		}
	}

	// While in general Mesos followers should redirect any requests for
	// `mesos-master-state` to the leader so that all output files should be
	// similar, defensively skip checking non-leaders.
	if state.ID != state.LeaderInfo.ID {
		return checks.Result{
			Status: checks.SOK,
			Value:  "Node is not the current leader"}
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
