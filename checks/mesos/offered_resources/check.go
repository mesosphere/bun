package offeredresources

import (
	"strings"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

func init() {
	check := checks.Check{
		Name: "offered-resources",
		Description: "Checks if any framework is hoarding resources by considering them for a long time. " +
			"Example: https://jira.d2iq.com/browse/COPS-211",
		Cure:           "Update the framework to a new version or restart it to reset",
		OKSummary:      "No frameworks hoarding resources",
		ProblemSummary: "Some of the frameworks is hoarding resources",
		Run:            checkFunc,
	}
	checks.RegisterCheck(check)
}

func checkFunc(b bundle.Bundle) checks.Results {
	var frameworks Frameworks
	err := b.ReadFirstJSON("mesos-master-frameworks", &frameworks)
	if err != nil {
		return checks.Results{
			checks.Result{
				Status: checks.SUndefined,
				Value:  "Couldn't read any of the mesos-framework JSONs. The last error: " + err.Error(),
			},
		}
	}
	badFrameworks := make([]string, 0, len(frameworks.Frameworks))
	for _, framework := range frameworks.Frameworks {
		o := framework.OfferedResources
		if o.Mem != 0 || o.Cpus != 0 || o.Disk != 0 || o.Ports != "" {
			badFrameworks = append(badFrameworks, framework.Name+" ("+framework.ID+")")
		}
	}
	if len(badFrameworks) == 0 {
		return checks.Results{
			checks.Result{
				Status: checks.SOK,
			},
		}
	}
	return checks.Results{
		checks.Result{
			Status: checks.SProblem,
			Value:  "The following frameworks might be hoarding resources: " + strings.Join(badFrameworks, ", "),
		},
	}
}
