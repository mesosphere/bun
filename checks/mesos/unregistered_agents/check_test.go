package unregisteredagents

import (
	"testing"

	"github.com/mesosphere/bun/v2/bundle"

	"github.com/mesosphere/bun/v2/checks"
)

func TestCheckPositive(t *testing.T) {
	check := checks.GetCheck("mesos-unregistered-agents")
	b, err := bundle.New("test_bundles/non_active_agents")
	if err != nil {
		t.Fatal(err)
	}
	problems := check.Run(b).Problems()
	if len(problems) != 2 {
		t.Fatalf("Expected 2 problems, observe %v", len(problems))
	}
	problemsMap := make(map[bundle.IP]checks.Result)
	for _, p := range problems {
		problemsMap[p.Host.IP] = p
	}
	if _, ok := problemsMap["10.0.3.6"]; !ok {
		t.Fatal("Problems should contain agent 10.0.3.6")
	}
	if p, _ := problemsMap["10.0.3.6"]; p.Host.Type != bundle.DTAgent {
		t.Fatalf("Agent 10.0.3.6 should be a private agent, it's %v instead", p.Host.Type)
	}
	if p, _ := problemsMap["10.0.3.6"]; p.Value != strAgentInactive {
		t.Fatalf("Problem with 10.0.3.6 should have value \"%v\", but has \"%v\" instead", strAgentInactive, p.Value)
	}

	if _, ok := problemsMap["10.0.6.64"]; !ok {
		t.Fatal("Problems should contain agent 10.0.6.64")
	}
	if p, _ := problemsMap["10.0.6.64"]; p.Host.Type != bundle.DTPublicAgent {
		t.Fatalf("Agent 10.0.6.64 should be a public agent, it's %v instead", p.Host.Type)
	}
	if p, _ := problemsMap["10.0.6.64"]; p.Value != strMesosAgentRecovered {
		t.Fatalf("Problem with 10.0.6.64 should have value \"%v\", but has \"%v\" instead", strMesosAgentRecovered, p.Value)
	}
}

func TestCheckNegative(t *testing.T) {
	check := checks.GetCheck("mesos-unregistered-agents")
	b, err := bundle.New("test_bundles/ok")
	if err != nil {
		t.Fatal(err)
	}
	results := check.Run(b)
	if len(results.Problems()) != 0 || len(results.Undefined()) != 0 {
		t.Fatalf("The check expected to be negative, instead %v problems and %v undefined found",
			len(results.Problems()), len(results.Undefined()))
	}
}
