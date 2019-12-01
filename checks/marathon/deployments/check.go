package deployments

import (
	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

// max number of Marathon deployments considered healthy
const maxDeployments = 10

func init() {
	check := checks.Check{
		Name:           "marathon-deployments",
		Description:    "Checks if Marathon has too many deployments",
		Cure:           "Too many deployments can mean that the cluster lost resources during some incident.",
		OKSummary:      "Marathon has less than 10 deployments.",
		ProblemSummary: "Marathon has more than 10 deployments.",
		Run:            checkFunc,
	}
	checks.RegisterCheck(check)
}

func checkFunc(b bundle.Bundle) checks.Results {
	deployments := make([]struct{}, 0)
	if len(b.Masters()) == 0 {
		return checks.Results{checks.Result{
			Status: checks.SUndefined,
			Value:  "Cannot find DC/OS masters to perform the check",
		}}
	}

	master := b.Masters()[0]
	if err := master.ReadJSON("marathon-deployments", &deployments); err != nil {
		return checks.Results{
			{
				Status: checks.SUndefined,
				Value:  err,
			},
		}
	}
	if len(deployments) > maxDeployments {
		return checks.Results{
			checks.Result{
				Status: checks.SProblem,
				Value:  len(deployments),
			},
		}
	}
	return checks.Results{
		checks.Result{
			Status: checks.SOK,
			Value:  len(deployments),
		},
	}
}
