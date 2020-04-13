package cluster

import (
	"fmt"

	"github.com/mesosphere/bun/v2/bundle"
)

type Cluster struct {
	b bundle.Bundle
}

type Node struct {
	Host bundle.Host
}

func (n Node) IP() bundle.IP {
	return n.Host.IP
}

func New(b bundle.Bundle) Cluster {
	return Cluster{b: b}
}

func (c Cluster) MesosLeader() (MesosMaster, error) {
	var leader MesosMaster
	masters, err := c.MesosMasters()
	if err != nil {
		return leader, err
	}
	leaderFound := false
	for _, m := range masters {
		if !m.IsLeader {
			continue
		}
		if !leaderFound {
			leader = m
			leaderFound = true
			continue
		}
		return leader, fmt.Errorf("at least two Mesos leaders found: %v and %v", leader.IP(), m.IP())
	}
	if !leaderFound {
		return leader, fmt.Errorf("no Mesos leader found")
	}
	return leader, nil
}

func (c Cluster) MesosMasters() ([]MesosMaster, error) {
	masters := make([]MesosMaster, 0, len(c.b.Masters()))
	for _, m := range c.b.Masters() {
		master, err := NewMesosMaster(m)
		if err != nil {
			return nil, err
		}
		masters = append(masters, master)
	}
	return masters, nil
}
