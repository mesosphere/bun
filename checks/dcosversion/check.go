package dcosversion

import (
	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

func init() {
	builder := checks.CheckFuncBuilder{
		CollectFromMasters:      collect,
		CollectFromAgents:       collect,
		CollectFromPublicAgents: collect,
		Aggregate:               aggregate,
	}
	check := checks.Check{
		Name: "dcos-version",
		Description: "Verify that all hosts in the cluster have the " +
			"same DC/OS version installed",
		Cure:           "Upgrade the nodes which have older DC/OS versions.",
		OKSummary:      "All the nodes have the same DC/OS version.",
		ProblemSummary: "The nodes have different DC/OS versions installed.",
		CheckFunc:      builder.Build(),
	}
	checks.RegisterCheck(check)
}

// Version represents the dcos-version JSON file
type Version struct {
	Version string
}

func collect(host bundle.Host) checks.Detail {
	v := Version{}
	if err := host.ReadJSON("dcos-version", &v); err != nil {
		return checks.Detail{
			Status: checks.SUndefined,
			Err:    err,
		}
	}
	return checks.Detail{
		Status: checks.SOK,
		Value:  v.Version,
	}
}

func aggregate(details checks.Details) checks.Details {
	ok := true
	var version string
	for _, detail := range details {
		if version == "" {
			version = detail.Value.(string)
		}
		if version != detail.Value.(string) {
			ok = false
			break
		}
	}
	if ok {
		return details
	}
	for _, d := range details {
		d.Status = checks.SProblem
	}
	return details
}
