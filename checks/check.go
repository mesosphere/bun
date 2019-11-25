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

// MsgErr is a standard message used in the check summary when errors
// occurs during the check.
const MsgErr = "Error(s) occurred while performing the check."

// cluster analyzing its diagnostics bundle. Check supposed to populate fields
// of the CheckResult.
// Each check should implement this interface.

// Check checks some aspect of the DC/OS cluster analyzing its diagnostics
// bundle.
// Checks can be registered in the check registry with the registerCheck function.
// Check is not supposed to be run more than one time.
type Check struct {
	Name           string    `yaml:"name"`           // Required
	Description    string    `yaml:"description"`    // Required
	Cure           string    `yaml:"cure"`           // Required
	OKSummary      string    `yaml:"okSummary"`      // Optional
	ProblemSummary string    `yaml:"problemSummary"` // Optional
	CheckFunc      CheckFunc // Required
}

type CheckFunc func(bundle.Bundle) Details

type Detail struct {
	Status Status
	Value  interface{}
	Host   *bundle.Host
	Err    error
}

type Details []Detail

func (d Details) filter(status Status) Details {
	details := make([]Detail, 0, len(d))
	for _, detail := range d {
		if detail.Status == status {
			details = append(details, detail)
		}
	}
	return details
}

func (d Details) Problems() Details {
	return d.filter(SProblem)
}

func (d Details) Failures() Details {
	return d.filter(SUndefined)
}

func (d Details) OKs() Details {
	return d.filter(SOK)
}

// Run runs the check.
func (c *Check) Run(b bundle.Bundle) (status Status, summary string, details Details) {
	details = c.CheckFunc(b)
	if len(details.Problems()) > 0 {
		status = SProblem
		summary = c.ProblemSummary
		return
	}
	if len(details.Failures()) > 0 {
		status = SUndefined
		summary = MsgErr
		return
	}
	summary = c.OKSummary
	status = SOK
	return
}
