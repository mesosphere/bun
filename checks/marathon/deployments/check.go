package deployments

import (
	"fmt"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

// max number of Marathon deployments considered healthy
const maxDeployments = 10

func init() {
	builder := checks.CheckBuilder{
		Name:               "marathon-deployments",
		Description:        "Check for too many running Marathon app deployments",
		CollectFromMasters: collect,
		Aggregate:          checks.DefaultAggregate,
	}
	check := builder.Build()
	checks.RegisterCheck(check)
}

func collect(host bundle.Host) (ok bool, details interface{}, err error) {
	deployments := []struct{}{}
	if err = host.ReadJSON("marathon-deployments", &deployments); err != nil {
		return
	}
	if len(deployments) > maxDeployments {
		details = fmt.Sprintf("Too many deployments: %v", len(deployments))
		return
	}
	ok = true
	return
}
