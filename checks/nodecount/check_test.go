package nodecount

import (
	"testing"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

func TestOK(t *testing.T) {
	b, err := bundle.New("test_bundles/ok")
	if err != nil {
		t.Fatal(err)
	}
	c := checks.GetCheck("node-count")
	results := c.Run(b)
	if len(results.Undefined()) > 0 {
		t.Errorf("Expected zero undefined results, observed %v.", results.Undefined())
	}
	if results.Status() != checks.SOK {
		t.Errorf("Expected Status = OK, observed Status = %v.", results.Status())
	}
	if len(results.Problems()) > 0 {
		t.Errorf("Expected Problems is empty, observed Problems has size %v.", len(results.Problems()))
	}
}

func TestMastersCount(t *testing.T) {
	b, err := bundle.New("test_bundles/problem")
	if err != nil {
		t.Fatal(err)
	}
	c := checks.GetCheck("node-count")
	results := c.Run(b)
	if len(results.Undefined()) > 0 {
		t.Errorf("Expected zero undefined results, observed %v.", results.Undefined())
	}
	if results.Status() != checks.SProblem {
		t.Errorf("Expected Status = Problem, observed Status = %v.", results.Status())
	}
	if len(results.Problems()) != 1 {
		t.Errorf("Expected one Problem, observed Problems has size %v.", len(results.Problems()))
	}
}
