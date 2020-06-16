package cmd

import (
	_ "github.com/mesosphere/bun/v2/checks/dcosnet/vips"
	_ "github.com/mesosphere/bun/v2/checks/dcosrequirements/cpu"
	_ "github.com/mesosphere/bun/v2/checks/dcosrequirements/disk"
	_ "github.com/mesosphere/bun/v2/checks/dcosrequirements/mem"
	_ "github.com/mesosphere/bun/v2/checks/dcosversion"
	_ "github.com/mesosphere/bun/v2/checks/health"
	_ "github.com/mesosphere/bun/v2/checks/marathon/deployments"
	_ "github.com/mesosphere/bun/v2/checks/marathon/marathon_lb_1.14.1"
	_ "github.com/mesosphere/bun/v2/checks/mesos/actormailboxes"
	_ "github.com/mesosphere/bun/v2/checks/mesos/mesos9868"
	_ "github.com/mesosphere/bun/v2/checks/mesos/unregistered_agents"
	_ "github.com/mesosphere/bun/v2/checks/nodecount"
)
