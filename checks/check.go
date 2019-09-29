package checks

import "github.com/mesosphere/bun/v2/bundle"

// Status defines possible check outcomes.
type Status string

const (
	// SUndefined means that the check wasn't performed successfully.
	SUndefined Status = "UNDEFINED"
	// SOK means that the bundle passed the check.
	SOK = "OK"
	// SProblem means that the bundle failed to pass the check.
	SProblem = "PROBLEM"
)

// cluster analyzing its diagnostics bundle. Check supposed to populate fields
// of the CheckResult.
// Each check should implement this interface.

// Check checks some aspect of the DC/OS cluster analyzing its diagnostics
// bundle.
// Checks can be registered in the check registry with the registerCheck function.
// Check is not supposed to be run more than one time.
type Check struct {
	Name        string                      // Required
	Description string                      // Optional
	CheckFunc   func(*Check, bundle.Bundle) // Required
	Status      Status                      // Do not set
	Summary     string                      // Do not set
	Problems    []string                    // Do not set
	Errors      []string                    // Do not set
	OKs         []string                    // Do not set
}

// Run runs the check.
func (c *Check) Run(b bundle.Bundle) {
	c.CheckFunc(c, b)
}
