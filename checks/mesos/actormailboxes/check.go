package actormailboxes

import (
	"fmt"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

// number of events in an actor's mailbox after which the actor is
// considered backlogged
const maxEvents = 30

func init() {
	builder := checks.CheckFuncBuilder{
		CollectFromMasters:      collect,
		CollectFromAgents:       collect,
		CollectFromPublicAgents: collect,
	}
	check := checks.Check{
		Name: "mesos-actor-mailboxes",
		Description: "Check if actor mailboxes in the Mesos process " +
			"have a reasonable amount of messages",
		Cure: "Check I/O on the correspondent hosts and if something is overloading Mesos agents or masters" +
			" with API calls.",
		OKSummary:      "All Mesos actors are fine.",
		ProblemSummary: "Some Mesos actors are backlogged.",
		CheckFunc:      builder.Build(),
	}
	checks.RegisterCheck(check)
}

type MesosActor struct {
	ID     string `json:"id"`
	Events []struct{}
}

func collect(host bundle.Host) checks.Detail {
	var actors []MesosActor
	if err := host.ReadJSON("mesos-processes", &actors); err != nil {
		return checks.Detail{
			Status: checks.SUndefined,
			Err:    err,
		}
	}
	var mailboxes []string
	for _, a := range actors {
		if len(a.Events) > maxEvents {
			mailboxes = append(mailboxes, fmt.Sprintf("(Mesos) %v@%v: mailbox size = %v (> %v)",
				a.ID, host.IP, len(a.Events), maxEvents))
		}
	}
	if len(mailboxes) > 0 {
		return checks.Detail{
			Status: checks.SProblem,
			Value:  mailboxes,
		}
	}
	return checks.Detail{
		Status: checks.SOK,
	}
}
