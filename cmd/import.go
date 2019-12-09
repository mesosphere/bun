package cmd

import (
	_ "github.com/mesosphere/bun/v2/checks/dcosrequirements"
	_ "github.com/mesosphere/bun/v2/checks/dcosversion"
	_ "github.com/mesosphere/bun/v2/checks/health"
	_ "github.com/mesosphere/bun/v2/checks/marathon/deployments"
	_ "github.com/mesosphere/bun/v2/checks/mesos/actormailboxes"
	_ "github.com/mesosphere/bun/v2/checks/mesos/mesos9868"
	_ "github.com/mesosphere/bun/v2/checks/mesos/unregistered_agents"
	_ "github.com/mesosphere/bun/v2/checks/nodecount"
)
