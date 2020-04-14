package cpu

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

// Requirement values are taken from the dcos-2.0 documentation,
// https://docs.d2iq.com/mesosphere/dcos/2.0/installing/production/system-requirements/.
var (
	cpuRequirements = map[bundle.DirType]int{
		bundle.DTMaster:      4,
		bundle.DTAgent:       2,
		bundle.DTPublicAgent: 2,
	}
)

func init() {
	builder := checks.CheckFuncBuilder{
		CheckMasters:      checkCpus,
		CheckAgents:       checkCpus,
		CheckPublicAgents: checkCpus,
	}
	check := checks.Check{
		Name:        "cpu",
		Description: "Checks that DC/OS nodes meet CPU requirements",
		Cure: "Add CPUs to the problematic nodes to meet system requirements: " +
			"https://docs.d2iq.com/mesosphere/dcos/2.0/installing/production/system-requirements/",
		OKSummary:      "All nodes meet CPU requirements",
		ProblemSummary: "Some nodes do not meet CPU requirements",
		Run:            builder.Build(),
	}
	checks.RegisterCheck(check)
}

func checkCpus(host bundle.Host) checks.Result {
	cpuinfo, err := host.OpenFile("cpuinfo")
	if err != nil {
		return checks.Result{
			Status: checks.SUndefined,
			Host:   host,
			Value:  "Couldn't check. Error: " + err.Error(),
		}
	}
	defer cpuinfo.Close()
	scanner := bufio.NewScanner(cpuinfo)
	numCpus := 0
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "processor\t:") {
			numCpus++
		}
	}
	if numCpus < cpuRequirements[host.Type] {
		return checks.Result{
			Status: checks.SProblem,
			Value: fmt.Sprintf(
				"node has less than required CPUs: %.f%% (%d vs. %d)",
				100*float64(numCpus)/float64(cpuRequirements[host.Type]),
				numCpus,
				cpuRequirements[host.Type])}
	}
	return checks.Result{
		Status: checks.SOK,
		Host:   host,
	}
}
