package nodecount

import (
	"fmt"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

func init() {
	check := checks.Check{
		Name:           "node-count",
		Description:    "Checks if the cluster has 3 or 5 masters and more than 0 nodes.",
		Cure:           "Check Mesos logs of the disconnected masters and agents.",
		OKSummary:      "Cluster has correct amount of masters and more than 0 agents.",
		ProblemSummary: "Cluster doesn't have correct amount of masters or agents.",
		Run:            checkFunc,
	}
	checks.RegisterCheck(check)
}

func checkFunc(b bundle.Bundle) checks.Results {
	nMasters := len(b.Masters())
	nAgents := len(b.Agents()) + len(b.PublicAgents())
	if (nMasters == 3 || nMasters == 5) && nAgents != 0 {
		return []checks.Result{{Status: checks.SOK}}
	}
	return []checks.Result{
		{
			Status: checks.SProblem,
			Value: fmt.Sprintf("Expected 3 or 5 masters and more than 0 agents, observe %v masters and %v agents.",
				nMasters, nAgents),
		},
	}
}
