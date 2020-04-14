package mem

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

// Requirement values are taken from the dcos-2.0 documentation,
// https://docs.d2iq.com/mesosphere/dcos/2.0/installing/production/system-requirements/.
var (
	memRequirements = map[bundle.DirType]int{
		bundle.DTMaster:      32_000_000,
		bundle.DTAgent:       16_000_000,
		bundle.DTPublicAgent: 16_000_000,
	}
)

func init() {
	builder := checks.CheckFuncBuilder{
		CheckMasters:      checkMem,
		CheckAgents:       checkMem,
		CheckPublicAgents: checkMem,
	}
	check := checks.Check{
		Name:        "mem",
		Description: "Checks that DC/OS nodes meet memory requirements",
		Cure: "Add memory to the problematic nodes to meet system requirements: " +
			"https://docs.d2iq.com/mesosphere/dcos/2.0/installing/production/system-requirements/`",
		OKSummary:      "All nodes meet memory requirements",
		ProblemSummary: "Some nodes do not meet memory requirements",
		Run:            builder.Build(),
	}
	checks.RegisterCheck(check)
}

func checkMem(host bundle.Host) checks.Result {
	meminfo, err := host.OpenFile("meminfo")
	if err != nil {
		return checks.Result{
			Status: checks.SUndefined,
			Host:   host,
			Value:  "Couldn't check. Error: " + err.Error(),
		}
	}
	defer meminfo.Close()
	scanner := bufio.NewScanner(meminfo)
	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "MemTotal") {
			continue
		}
		tokens := strings.Split(scanner.Text(), " ")
		mem, err := strconv.Atoi(tokens[len(tokens)-2])
		if err != nil {
			return checks.Result{
				Status: checks.SUndefined,
				Host:   host,
				Value:  "Couldn't check. Error: " + err.Error(),
			}
		}
		if mem < memRequirements[host.Type] {
			return checks.Result{
				Status: checks.SProblem,
				Host:   host,
				Value: fmt.Sprintf(
					"node has less than required memory: %.f%% (%.2f GB vs. %.2f GB)",
					100*float64(mem)/float64(memRequirements[host.Type]),
					convertKBtoGB(mem),
					convertKBtoGB(memRequirements[host.Type])),
			}
		}
		return checks.Result{
			Status: checks.SOK,
			Host:   host,
		}
	}
	return checks.Result{
		Status: checks.SUndefined,
		Host:   host,
		Value: "Couldn't find MemTotal in one the following files: " +
			strings.Join(bundle.GetFileType("meminfo").Paths, ", "),
	}
}

func convertKBtoGB(kb int) float64 {
	return float64(kb) / 1e6
}
