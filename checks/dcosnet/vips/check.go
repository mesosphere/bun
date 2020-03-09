package unregisteredagents

import (
	"fmt"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"

	"github.com/lithammer/dedent"
)

func init() {
	builder := checks.CheckFuncBuilder{
		CheckMasters:      collectMasters,
		CheckAgents:       collectVIPs,
		CheckPublicAgents: collectVIPs,
		Aggregate:         aggregate,
	}

	check := checks.Check{
		Name:        "dcosnet-vips",
		Description: "Checks for wrong VIP backends",
		Cure: dedent.Dedent(`
			Checks if for every VIP in the cluster there is a corresponding container
			in running state. It also implicitly checks if the mesos and VIP states
			are synchronized.`),
		OKSummary:      "All VIPs have a corresponding live backends",
		ProblemSummary: "Some VIPs do not have a corresponding live backend",
		Run:            builder.Build(),
	}

	checks.RegisterCheck(check)
}

type containerId struct {
	Value string `json:"value"`
}
type ipAddress struct {
	Protocol  string `json:"protocol"`
	IPAddress string `json:"ip_address"`
}

type containerNetworkInfo struct {
	IPAddresses []ipAddress `json:"ip_addresses"`
}

type containerStatus struct {
	ContainerID  containerId            `json:"container_id"`
	NetworkInfos []containerNetworkInfo `json:"network_infos"`
}

type taskStatus struct {
	ContainerStatus containerStatus `json:"container_status"`
}

type frameworkTask struct {
	Statuses []taskStatus `json:"statuses"`
	State    string       `json:"state"`
	Name     string       `json:"name"`
}

type framework struct {
	Tasks []frameworkTask `json:"tasks"`
	Name  string          `json:"name"`
	ID    string          `json:"id"`
}

type state struct {
	AgentID    string      `json:"id"`
	Frameworks []framework `json:"frameworks"`
}

type vipBackend struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type vip struct {
	Name     string       `json:"vip"`
	Backends []vipBackend `json:"backend"`
}

type ipMappingInfo struct {
	Framework *framework
	Container *containerStatus
	AgentID   string
}

type scanResult struct {
	IPMapping map[string]ipMappingInfo
	VIPs      []vip
	Faults    []checks.Result
}

func collectVIPs(host bundle.Host) checks.Result {
	var vips []vip

	if err := host.ReadJSON("vips", &vips); err != nil {
		return checks.Result{
			Status: checks.SUndefined,
			Value:  err,
		}
	}

	return checks.Result{
		Status: checks.SOK,
		Value:  &scanResult{nil, vips, nil},
	}
}

func collectMasters(host bundle.Host) checks.Result {
	var state state
	var faults []checks.Result

	if err := host.ReadJSON("mesos-master-state", &state); err != nil {
		return checks.Result{
			Status: checks.SUndefined,
			Value:  err,
		}
	}

	ipMapping := make(map[string]ipMappingInfo)
	for _, fw := range state.Frameworks {
		for _, task := range fw.Tasks {
			if task.State != "TASK_RUNNING" {
				continue
			}

			lastStatus := task.Statuses[len(task.Statuses)-1]
			container := lastStatus.ContainerStatus
			for _, netInfo := range container.NetworkInfos {
				if len(netInfo.IPAddresses) == 0 {
					continue
				}

				for _, addr := range netInfo.IPAddresses {
					if value, ok := ipMapping[addr.IPAddress]; ok {
						faults = append(faults, checks.Result{
							Status: checks.SUndefined,
							Value: fmt.Errorf("More than one containers are using the same IP (%s and %s)",
								value.Container.ContainerID.Value, container.ContainerID.Value),
						})
					} else {
						ipMapping[addr.IPAddress] = ipMappingInfo{
							Framework: &fw,
							Container: &container,
							AgentID:   state.AgentID,
						}
					}
				}
			}
		}
	}

	ret := collectVIPs(host)
	if ret.Status != checks.SOK {
		return ret
	}

	vipStatus := ret.Value.(*scanResult)

	return checks.Result{
		Status: checks.SOK,
		Value:  &scanResult{ipMapping, vipStatus.VIPs, append(faults, vipStatus.Faults...)},
	}
}

func aggregate(r checks.Results) checks.Results {
	var results checks.Results = nil

	// Merge all IP Mappings
	ipMappings := make(map[string]ipMappingInfo)
	for _, d := range r.OKs() {
		result := d.Value.(*scanResult)

		for ip, info := range result.IPMapping {
			if value, ok := ipMappings[ip]; ok {
				if value.Container.ContainerID.Value != info.Container.ContainerID.Value {
					results = append(results,
						checks.Result{
							Status: checks.SProblem,
							Value: fmt.Sprintf(
								"Mesos on agent %s reports that container %s has IP %s, while mesos on agent %s reports that this IP belongs to container %s",
								value.AgentID, value.Container.ContainerID.Value, ip,
								info.AgentID, info.Container.ContainerID.Value,
							)},
					)
				}
			} else {
				ipMappings[ip] = info
			}
		}
	}

	// Lookup vips
	for _, d := range r.OKs() {
		result := d.Value.(*scanResult)
		for _, vip := range result.VIPs {
			if len(vip.Backends) == 0 {
				results = append(results,
					checks.Result{
						Status: checks.SProblem,
						Value: fmt.Sprintf(
							"The VIP '%s' has no back-ends defined",
							vip.Name,
						)},
				)
				continue
			}

			for _, ip := range vip.Backends {
				if _, ok := ipMappings[ip.IP]; !ok {
					results = append(results,
						checks.Result{
							Status: checks.SProblem,
							Value: fmt.Sprintf(
								"The backend %s of VIP '%s' was not found on any mesos container",
								ip.IP, vip.Name,
							)},
					)
				}
			}
		}
	}

	return results
}
