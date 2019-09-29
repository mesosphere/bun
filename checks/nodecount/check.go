package nodecount

import (
	"fmt"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

func init() {
	check := checks.Check{
		Name: "node-count",
		Description: "Count nodes of each type, checks if the amount of " +
			"master nodes is odd",
		CheckFunc: checkFunc,
	}
	checks.RegisterCheck(check)
}

func checkFunc(c *checks.Check, b bundle.Bundle) {
	lenMasters := len(b.Masters)
	stats := fmt.Sprintf(
		"Masters: %v, Agents: %v, Public Agents: %v, Total: %v",
		lenMasters, len(b.Agents), len(b.PublicAgents), len(b.Hosts))
	if lenMasters == 1 || lenMasters == 3 || lenMasters == 5 {
		c.Status = checks.SOK
		c.Summary = stats
		return
	}
	c.Status = checks.SProblem
	c.Summary = fmt.Sprintf("Number of masters is not valid. %v", stats)
	c.Problems = append(c.Problems, "It should be 1, 3 or 5 masters.")
}
