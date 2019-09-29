package actormailboxes

import (
	"fmt"
	"strings"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

// number of events in an actor's mailbox after which the actor is
// considered backlogged
const maxEvents = 30

func init() {
	builder := checks.CheckBuilder{
		Name: "mesos-actor-mailboxes",
		Description: "Check if actor mailboxes in the Mesos process " +
			"have a reasonable amount of messages",
		OKSummary:               "All Mesos actors are fine.",
		ProblemSummary:          "Some Mesos actors are backlogged.",
		CollectFromMasters:      collect,
		CollectFromAgents:       collect,
		CollectFromPublicAgents: collect,
		Aggregate:               checks.DefaultAggregate,
	}
	check := builder.Build()
	checks.RegisterCheck(check)
}

// MesosActor represents the structure of the __processess__ file.
type MesosActor struct {
	ID     string `json:"id"`
	Events []struct{}
}

func collect(host bundle.Host) (ok bool, details interface{}, err error) {
	var actors []MesosActor
	if err = host.ReadJSON("mesos-processes", &actors); err != nil {
		return
	}
	var u []string
	for _, a := range actors {
		if len(a.Events) > maxEvents {
			u = append(u, fmt.Sprintf("(Mesos) %v@%v: mailbox size = %v (> %v)",
				a.ID, host.IP, len(a.Events), maxEvents))
		}
	}
	if len(u) > 0 {
		details = "The following Mesos actor mailboxes are too big:\n" +
			strings.Join(u, "\n")
		return
	}
	ok = true
	return
}
