package cmd

import (
	_ "github.com/mesosphere/bun/checks/dcosversion"
	_ "github.com/mesosphere/bun/checks/health"
	_ "github.com/mesosphere/bun/checks/marathon/deployments"
	_ "github.com/mesosphere/bun/checks/mesos/actormailboxes"
	_ "github.com/mesosphere/bun/checks/mesos/mesos9868"
	_ "github.com/mesosphere/bun/checks/networking/errors"
	_ "github.com/mesosphere/bun/checks/nodecount"
)
