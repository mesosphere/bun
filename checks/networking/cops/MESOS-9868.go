package errors

import (
	"fmt"
	"strings"

	"github.com/mesosphere/bun"
)

type NetworkInfo struct {
	IPAddresses []struct {
		IPAddress string `json:"ip_address,omitempty"`
	} `json:"ip_addresses,omitempty"`
}

type MesosContainer struct {
	ContainerID string `json:"container_id"`
	Status      struct {
		NetworkInfos []NetworkInfo `json:"network_infos,omitempty"`
	} `json:"status,omitempty"`
}

type MesosTaskStatus struct {
	Timestamp float64 `json:"timestamp"`

	ContainerStatus struct {
		ContainerID struct {
			Value  string
			Parent struct {
				Value string `json:"value"`
			} `json:"parent"`
		} `json:"container_id"`
		NetworkInfos []NetworkInfo `json:"network_infos,omitempty"`
	} `json:"container_status"`
}

type MesosState struct {
	Frameworks []struct {
		Name  string `json:"name"`
		Tasks []struct {
			Name     string            `json:"name"`
			Statuses []MesosTaskStatus `json:"statuses"`
		} `json:"tasks"`
	} `json:"frameworks"`
}

type ContainerIPs struct {
	TaskIPs      []string
	ContainerIPs []string
}

func init() {
	builder := bun.CheckBuilder{
		Name:                    "MESOS-9868",
		Description:             "Affected by MESOS-9868 error",
		CollectFromMasters:      collectMasters,
		CollectFromAgents:       collectAgents,
		CollectFromPublicAgents: collectAgents,
		Aggregate:               aggregate,
	}
	check := builder.Build()
	bun.RegisterCheck(check)
}

func collectMasters(host bun.Host) (ok bool, details interface{}, err error) {
	var mesosState MesosState

	ok = true
	err = host.ReadJSON("mesos-master-state", &mesosState)
	if err != nil {
		ok = false
		fmt.Println(fmt.Sprintf("Unable to parse: %s", err))
		return
	}

	summary := make(map[string]ContainerIPs)

	for _, fw := range mesosState.Frameworks {
		for _, t := range fw.Tasks {

			var lastStatus *MesosTaskStatus = nil
			for _, s := range t.Statuses {
				if lastStatus == nil || s.Timestamp > lastStatus.Timestamp {
					lastStatus = &s
				}
			}

			var ips ContainerIPs
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

func collectAgents(host bun.Host) (ok bool, details interface{}, err error) {
	var containerIPs []MesosContainer

	ok = true
	err = host.ReadJSON("mesos-agent-containers", &containerIPs)
	if err != nil {
		ok = false
		fmt.Println(fmt.Sprintf("Unable to parse: %s", err))
		return
	}

	summary := make(map[string]ContainerIPs)

	for _, c := range containerIPs {
		var ips ContainerIPs
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

func aggregate(c *bun.Check, b bun.CheckBuilder) {
	oks := []string{}
	probs := []string{}

	// Collect all the IPs from tasks and containers
	summary := make(map[string]ContainerIPs)
	for _, r := range b.OKs {
		for agent, ips := range r.Details.(map[string]ContainerIPs) {
			v := summary[agent]
			v.TaskIPs = append(v.TaskIPs, ips.TaskIPs...)
			v.ContainerIPs = append(v.ContainerIPs, ips.ContainerIPs...)
			summary[agent] = v
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
				probs = append(probs, fmt.Sprintf("Container %s 5050-state IP %s dont match 5051-containers IP: %s",
					container, tip, strings.Join(ips.ContainerIPs, ", ")))
				break
			}
		}

		if found {
			oks = append(oks, "Container %s IPs match", container)
		}
	}

	if len(probs) > 0 {
		c.Summary = fmt.Sprintf("%d containers affected by MESOS-9868", len(probs))
	} else {
		c.Summary = "Not affected by MESOS-9868"
	}

	c.Problems = probs
	c.OKs = oks
}
