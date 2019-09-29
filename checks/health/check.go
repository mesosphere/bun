package health

import (
	"fmt"
	"strings"

	"github.com/mesosphere/bun/bundle"
	"github.com/mesosphere/bun/checks"
)

func init() {
	builder := checks.CheckBuilder{
		Name:                    "diagnostics-health",
		Description:             "Check if all DC/OS components are healthy",
		CollectFromMasters:      collect,
		CollectFromAgents:       collect,
		CollectFromPublicAgents: collect,
		Aggregate:               checks.DefaultAggregate,
	}
	check := builder.Build()
	checks.RegisterCheck(check)
}

// Host represents the "host" object in the health JSON file
type Host struct {
	Units []Unit
}

// Unit represents the "unit" object in the health JSON file
type Unit struct {
	ID     string `json:"id"`
	Health int
}

func collect(host bundle.Host) (ok bool, details interface{}, err error) {
	h := Host{}
	if err = host.ReadJSON("diagnostics-health", &h); err != nil {
		return
	}
	var unhealthy []string
	for _, u := range h.Units {
		if u.Health != 0 {
			unhealthy = append(unhealthy,
				fmt.Sprintf("%v: health = %v", u.ID, u.Health))
		}
	}
	if len(unhealthy) > 0 {
		details = fmt.Sprintf("The following components are not healthy:\n%v",
			strings.Join(unhealthy, "\n"))
		ok = false
	} else {
		ok = true
	}
	return
}
