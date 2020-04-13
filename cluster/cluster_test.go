package cluster

import (
	"testing"

	"github.com/mesosphere/bun/v2/bundle"
)

func TestMesosLeaderFindsLeader(t *testing.T) {
	b, err := bundle.New("test_bundles/ok")
	if err != nil {
		t.Fatal(err)
	}
	c := New(b)
	leader, err := c.MesosLeader()
	if err != nil {
		t.Fatal(err)
	}
	expected := bundle.IP("10.0.5.2")
	if expected != leader.IP() {
		t.Fatalf("Expected %v, observed %v", expected, leader.IP())
	}
}

func TestMesosLeaderErrorWhenDuplicateIPs(t *testing.T) {
	b, err := bundle.New("test_bundles/duplicate_ids")
	if err != nil {
		t.Fatal(err)
	}
	c := New(b)
	_, err = c.MesosLeader()
	if err == nil {
		t.Fatal("Expected duplicate error")
	}
}
