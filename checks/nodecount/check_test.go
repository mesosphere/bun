package nodecount

import (
	"testing"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

func TestOK(t *testing.T) {
	b, err := bundle.NewBundle("test_bundles/ok")
	if err != nil {
		t.Fatal(err)
	}
	c := checks.GetCheck("node-count")
	c.Run(b)
	if len(c.Errors) > 0 {
		t.Fatal(c.Errors)
	}
	if c.Status != checks.SOK {
		t.Errorf("Expected Status = OK, observed Status = %v.", c.Status)
	}
	if len(c.Problems) > 0 {
		t.Errorf("Expected Problems is empty, observed Problems has size %v.",
			len(c.Problems))
	}
}

func TestMastersCount(t *testing.T) {
	b, err := bundle.NewBundle("test_bundles/problem")
	if err != nil {
		t.Fatal(err)
	}
	c := checks.GetCheck("node-count")
	c.Run(b)
	if len(c.Errors) > 0 {
		t.Fatal(c.Errors)
	}
	if c.Status != checks.SProblem {
		t.Errorf("Expected Status = Problem, observed Status = %v.", c.Status)
	}
	if len(c.Problems) != 1 {
		t.Errorf("Expected one Problem, observed Problems has size %v.",
			len(c.Problems))
	}
}

func TestOKSummary(t *testing.T) {
	const summary = "Masters: 3, Agents: 7, Public Agents: 2, Total: 12"
	b, err := bundle.NewBundle("test_bundles/ok")
	if err != nil {
		t.Fatal(err)
	}
	c := checks.GetCheck("node-count")
	c.Run(b)
	if len(c.Errors) > 0 {
		t.Fatal(c.Errors)
	}
	if c.Summary != summary {
		t.Errorf("Expected summary: \"%v\", summary, observed summary; \"%v\".",
			summary, c.Summary)
	}
}
