package cluster

import (
	"fmt"

	"github.com/mesosphere/bun/v2/bundle"
)

type MesosMaster struct {
	Node
	IsLeader bool
}

type masterState struct {
	Hostname string `json:"hostname"`
}

func NewMesosMaster(m bundle.Host) (MesosMaster, error) {
	master := MesosMaster{}
	if m.Type != bundle.DTMaster {
		return master, fmt.Errorf("host %v is not a master, it's %v", m.IP, m.Type)
	}
	state := masterState{}
	if err := m.ReadJSON("mesos-master-state", &state); err != nil {
		return master, err
	}
	master.Host = m
	// All masters will become the same ID
	master.IsLeader = m.IP == bundle.IP(state.Hostname)
	return master, nil
}
