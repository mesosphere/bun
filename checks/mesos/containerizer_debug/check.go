package containerizerdebug

import (
	"fmt"
	"encoding/json"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

func init() {
	builder := checks.CheckFuncBuilder{
		CheckAgents:       collect,
		CheckPublicAgents: collect,
	}
	check := checks.Check{
		Name: "mesos-containerizer-debug",
		Description: "Checks if the Mesos containerizer (UCR) has become unresponsive due to hanging synchronous syscalls",
		Cure: "Check that no Mesos-related OS processes stuck in a D-state on Mesos agents. Otherwise, this is a false positive alert. " +
			"Mesos containerizer (UCR) may become unresponsive due to a kernel bug such as kmem leak described in https://support.d2iq.com/s/article/Known-Issue-KMEM-with-Kubernetes-MSPH-2019-0002. " +
			"If the problem persists, reboot the agent node.",
		OKSummary:      "All Mesos agents are fine.",
		ProblemSummary: "Some Mesos agents may be stuck due to hanging Mesos containerizer processes.",
		Run:            builder.Build(),
	}
	checks.RegisterCheck(check)
}

type PendingOperations struct {
	Operations []struct {
		Operation string `json:"operation"`
		Args      interface{} `json:"args"`
	} `json:"pending"`
}

func collect(host bundle.Host) checks.Result {
	var pendingOperations PendingOperations

	if err := host.ReadJSON("mesos-agent-containerizer-debug", &pendingOperations); err != nil {
		return checks.Result{
			Status: checks.SUndefined,
			Value: err,
		}
	}
	if len(pendingOperations.Operations) > 0 {
		ops, _ := json.Marshal(pendingOperations.Operations)
		return checks.Result{
			Status: checks.SProblem,
			Value:  fmt.Sprintf("Mesos containerizer (UCR) has pending operations: %s", ops),
		}
	}
	return checks.Result{
		Status: checks.SOK,
	}
}
