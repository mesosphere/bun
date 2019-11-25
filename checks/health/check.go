package health

import (
	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

func init() {
	builder := checks.CheckFuncBuilder{
		CollectFromMasters:      collect,
		CollectFromAgents:       collect,
		CollectFromPublicAgents: collect,
	}
	check := checks.Check{
		Name:           "diagnostics-health",
		Description:    "Check if all DC/OS components are healthy",
		Cure:           "Check the logs of the unhealthy component.",
		OKSummary:      "All components are healthy.",
		ProblemSummary: "Found unhealthy components.",
		CheckFunc:      builder.Build(),
	}
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

func collect(host bundle.Host) checks.Detail {
	h := Host{}
	if err := host.ReadJSON("diagnostics-health", &h); err != nil {
		return checks.Detail{
			Status: checks.SUndefined,
			Err:    err,
		}
	}
	var unhealthy []string
	for _, u := range h.Units {
		if u.Health != 0 {
			unhealthy = append(unhealthy, u.ID)
		}
	}
	if len(unhealthy) > 0 {
		return checks.Detail{
			Status: checks.SProblem,
			Value:  unhealthy,
		}
	}
	return checks.Detail{
		Status: checks.SOK,
	}
}
