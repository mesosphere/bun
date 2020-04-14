package unregisteredagents

import (
	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
	"github.com/mesosphere/bun/v2/cluster"
)

const strAgentInactive = "Mesos agent appears to be registered but inactive"
const strMesosAgentRecovered = "Mesos agent was unregistered and recovered"

func init() {
	check := checks.Check{
		Name:        "mesos-unregistered-agents",
		Description: "Checks for unregistered Mesos agents",
		Cure: "Determine why Mesos agents cannot register by examining the Mesos agent log. " +
			"To register an agent needs to finish its recovery, and detect the master",
		OKSummary:      "All Mesos agents appear to be registered.",
		ProblemSummary: "Some Mesos agents appear to be unregistered.",
		Run:            check,
	}
	checks.RegisterCheck(check)
}

type slave struct {
	PID         string                 `json:"pid"`
	Active      bool                   `json:"active"`
	Deactivated bool                   `json:"deactivated"`
	HostName    string                 `json:"hostname"`
	Attributes  map[string]interface{} `json:"attributes"`
}

type state struct {
	Slaves          []slave `json:"slaves"`
	RecoveredSlaves []slave `json:"recovered_slaves"`
}

func check(b bundle.Bundle) checks.Results {
	c := cluster.New(b)
	leader, err := c.MesosLeader()
	if err != nil {
		return checks.Results{{Status: checks.SUndefined, Value: err}}
	}
	var s state
	if err := leader.Host.ReadJSON("mesos-master-state", &s); err != nil {
		return checks.Results{{Status: checks.SUndefined, Value: err}}
	}
	var unregistered []checks.Result
	for _, slave := range s.Slaves {
		if !slave.Active && !slave.Deactivated {
			var agent bundle.Host
			agent.IP = bundle.IP(slave.HostName)
			agent.Type = bundle.DTAgent
			if v := slave.Attributes["public_ip"]; v == "true" {
				agent.Type = bundle.DTPublicAgent
			}
			res := checks.Result{
				Status: checks.SProblem,
				Value:  strAgentInactive,
				Host:   agent,
			}
			unregistered = append(unregistered, res)
		}
	}
	for _, slave := range s.RecoveredSlaves {
		var agent bundle.Host
		agent.IP = bundle.IP(slave.HostName)
		agent.Type = bundle.DTAgent
		if v := slave.Attributes["public_ip"]; v == "true" {
			agent.Type = bundle.DTPublicAgent
		}
		res := checks.Result{
			Status: checks.SProblem,
			Value:  strMesosAgentRecovered,
			Host:   agent,
		}
		unregistered = append(unregistered, res)
	}
	if len(unregistered) > 0 {
		return unregistered
	}
	return checks.Results{{Status: checks.SOK}}
}
