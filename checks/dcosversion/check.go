package dcosversion

import (
	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

func init() {
	builder := checks.CheckFuncBuilder{
		CheckMasters:      collect,
		CheckAgents:       collect,
		CheckPublicAgents: collect,
		Aggregate:         aggregate,
	}
	check := checks.Check{
		Name: "dcos-version",
		Description: "Verifies that all hosts in the cluster have the " +
			"same DC/OS version installed",
		Cure:           "Upgrade the nodes which have older DC/OS versions.",
		OKSummary:      "All the nodes have the same DC/OS version.",
		ProblemSummary: "The nodes have different DC/OS versions installed.",
		Run:            builder.Build(),
	}
	checks.RegisterCheck(check)
}

// Version represents the dcos-version JSON file
type Version struct {
	Version string
}

func collect(host bundle.Host) checks.Result {
	v := Version{}
	if err := host.ReadJSON("dcos-version", &v); err != nil {
		return checks.Result{
			Status: checks.SUndefined,
			Value:  err,
		}
	}
	return checks.Result{
		Status: checks.SOK,
		Value:  v.Version,
	}
}

func aggregate(results checks.Results) checks.Results {
	ok := true
	var version string
	for _, result := range results.OKs() {
		if version == "" {
			version = result.Value.(string)
		}
		if version != result.Value.(string) {
			ok = false
			break
		}
	}
	if ok {
		return results
	}
	for _, d := range results {
		d.Status = checks.SProblem
	}
	results = append(results, results.Undefined()...)
	return results
}
