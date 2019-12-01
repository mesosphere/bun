package mesos9868

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

type networkInfo struct {
	IPAddresses []struct {
		IPAddress string `json:"ip_address"`
	} `json:"ip_addresses"`
}

type mesosContainer struct {
	ContainerID string `json:"container_id"`
	Status      struct {
		NetworkInfos []networkInfo `json:"network_infos"`
	} `json:"status"`
}

type mesosTaskStatus struct {
	Timestamp float64 `json:"timestamp"`

	ContainerStatus struct {
		ContainerID struct {
			Value  string
			Parent struct {
				Value string `json:"value"`
			} `json:"parent"`
		} `json:"container_id"`
		NetworkInfos []networkInfo `json:"network_infos"`
	} `json:"container_status"`
}

type mesosState struct {
	Frameworks []struct {
		Name  string `json:"name"`
		Tasks []struct {
			Name     string            `json:"name"`
			Statuses []mesosTaskStatus `json:"statuses"`
		} `json:"tasks"`
	} `json:"frameworks"`
}

type ipResults struct {
	TaskIPs      []string
	ContainerIPs []string
	HostIP       string
}

func init() {
	builder := checks.CheckFuncBuilder{
		CheckMasters:      collectMasters,
		CheckAgents:       collectAgents,
		CheckPublicAgents: collectAgents,
		Aggregate:         aggregate,
	}
	check := checks.Check{
		Name:           "mesos-9868",
		Description:    "Checks if the cluster is affected by the MESOS-9868 bug",
		Cure:           "Please, see https://issues.apache.org/jira/browse/MESOS-9868 for results.",
		OKSummary:      "The cluster is not affected by the MESOS-9868 bug.",
		ProblemSummary: "The cluster is affected by the MESOS-9868 bug.",
		Run:            builder.Build(),
	}
	checks.RegisterCheck(check)
}

func collectMasters(host bundle.Host) checks.Result {
	var mesosState mesosState

	if err := host.ReadJSON("mesos-master-state", &mesosState); err != nil {
		err = fmt.Errorf("unable to parse: %v", err)
		return checks.Result{
			Status: checks.SUndefined,
			Value:  err,
		}
	}
	summary := make(map[string]ipResults)
	for _, fw := range mesosState.Frameworks {
		for _, t := range fw.Tasks {
			var lastStatus *mesosTaskStatus
			for _, s := range t.Statuses {
				if lastStatus == nil || s.Timestamp > lastStatus.Timestamp {
					lastStatus = &s
				}
			}
			if lastStatus == nil {
				err := fmt.Errorf("no statuses found for task %s", t.Name)
				return checks.Result{
					Status: checks.SUndefined,
					Value:  err,
				}
			}
			var ips ipResults
			id := lastStatus.ContainerStatus.ContainerID.Value
			for _, n := range lastStatus.ContainerStatus.NetworkInfos {
				for _, i := range n.IPAddresses {
					ips.TaskIPs = append(ips.TaskIPs, i.IPAddress)
				}
			}
			summary[id] = ips
			// Do the same with the parent
			id = lastStatus.ContainerStatus.ContainerID.Parent.Value
			if id != "" {
				summary[id] = ips
			}
		}
	}
	return checks.Result{
		Value: summary,
	}
}

func collectAgents(host bundle.Host) checks.Result {
	var mesosContainers []mesosContainer

	if err := host.ReadJSON("mesos-agent-containers", &mesosContainers); err != nil {
		err = fmt.Errorf("unable to parse: %s", err.Error())
		return checks.Result{
			Status: checks.SUndefined,
			Value:  err,
		}
	}
	summary := make(map[string]ipResults)
	for _, c := range mesosContainers {
		var ips ipResults
		for _, n := range c.Status.NetworkInfos {
			for _, i := range n.IPAddresses {
				ips.ContainerIPs = append(ips.ContainerIPs, i.IPAddress)
			}
		}
		summary[c.ContainerID] = ips
	}
	return checks.Result{
		Status: checks.SOK,
		Value:  summary,
	}
}

func aggregate(r checks.Results) checks.Results {
	// Collect all the ipResults from tasks and containers
	summary := make(map[string]ipResults)
	for _, d := range r.OKs() {
		for container, ips := range d.Value.(map[string]ipResults) {
			v := summary[container]
			v.TaskIPs = append(v.TaskIPs, ips.TaskIPs...)
			v.ContainerIPs = append(v.ContainerIPs, ips.ContainerIPs...)
			summary[container] = v
		}
	}
	var results checks.Results = make([]checks.Result, 0, len(summary))
	results = append(results, results.Undefined()...)
	// See if they do not match
	for container, ips := range summary {
		// Ignore case where we weren't able to detect
		if len(ips.TaskIPs) == 0 || len(ips.ContainerIPs) == 0 {
			continue
		}
		found := false
		for _, tip := range ips.TaskIPs {
			for _, cip := range ips.ContainerIPs {
				if tip == cip {
					found = true
					break
				}
			}
			if !found {
				sort.Strings(ips.ContainerIPs)
				results = append(results,
					checks.Result{
						Status: checks.SProblem,
						Value: fmt.Sprintf("Container %s: IP %s from Mesos master \"/state\" endpoint "+
							"does not match any IP from Mesos agent \"/containers\" endpoint: %s.",
							container, tip, strings.Join(ips.ContainerIPs, ", "))})
				break
			}
		}
		if found {
			results = append(results,
				checks.Result{
					Status: checks.SOK,
					Value:  fmt.Sprintf("Container %s ipResults match", container),
				})
		}
	}
	return results
}
