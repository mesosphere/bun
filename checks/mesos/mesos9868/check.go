package mesos9868

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mesosphere/bun/bundle"
	"github.com/mesosphere/bun/checks"
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
	builder := checks.CheckBuilder{
		Name:                    "mesos-9868",
		Description:             "Check if the cluster is affected by the MESOS-9868 bug",
		CollectFromMasters:      collectMasters,
		CollectFromAgents:       collectAgents,
		CollectFromPublicAgents: collectAgents,
		Aggregate:               aggregate,
	}
	check := builder.Build()
	checks.RegisterCheck(check)
}

func collectMasters(host bundle.Host) (ok bool, details interface{}, err error) {
	var mesosState mesosState
	ok = true
	err = host.ReadJSON("mesos-master-state", &mesosState)
	if err != nil {
		err = fmt.Errorf("unable to parse: %s", err.Error())
		return
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
				err = fmt.Errorf("no statuses found for task %s", t.Name)
				return
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
	details = summary
	return
}

func collectAgents(host bundle.Host) (ok bool, details interface{}, err error) {
	var mesosContainers []mesosContainer
	ok = true
	err = host.ReadJSON("mesos-agent-containers", &mesosContainers)
	if err != nil {
		err = fmt.Errorf("unable to parse: %s", err.Error())
		return
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
	details = summary
	return
}

func aggregate(c *checks.Check, b checks.CheckBuilder) {
	var oks []string
	var problems []string
	// Collect all the ipResults from tasks and containers
	summary := make(map[string]ipResults)
	for _, r := range b.OKs {
		for container, ips := range r.Details.(map[string]ipResults) {
			v := summary[container]
			v.TaskIPs = append(v.TaskIPs, ips.TaskIPs...)
			v.ContainerIPs = append(v.ContainerIPs, ips.ContainerIPs...)
			summary[container] = v
		}
	}
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
				problems = append(problems,
					fmt.Sprintf("Container %s: IP %s from Mesos master \"/state\" endpoint "+
						"does not match any IP from Mesos agent \"/containers\" endpoint: %s.",
						container, tip, strings.Join(ips.ContainerIPs, ", ")))
				break
			}
		}

		if found {
			oks = append(oks, "Container %s ipResults match", container)
		}
	}

	if len(problems) > 0 {
		c.Summary = fmt.Sprintf("%d containers affected by MESOS-9868", len(problems))
	} else if len(c.Errors) > 0 {
		c.Summary = "Errors occurred during the check."
	}

	c.Problems = problems
	c.OKs = oks
}
