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

// Check checks some aspect of the DC/OS cluster analyzing its diagnostics
// bundle.
// Checks can be registered in the check registry with the registerCheck function.
type Check struct {
	Name           string          `yaml:"name"`           // Required
	Description    string          `yaml:"description"`    // Required
	Cure           string          `yaml:"cure"`           // Required
	OKSummary      string          `yaml:"okSummary"`      // Optional
	ProblemSummary string          `yaml:"problemSummary"` // Optional
	Run            CheckBundleFunc // Required
}

type CheckBundleFunc func(bundle.Bundle) Results

// Result represents check result.
type Result struct {
	Status Status
	Value  interface{}
	Host   bundle.Host
}

func (r Result) IsHostSet() bool {
	return r.Host.IP != ""
}

type Results []Result

func (d Results) filter(status Status) Results {
	results := make([]Result, 0, len(d))
	for _, result := range d {
		if result.Status == status {
			results = append(results, result)
		}
	}
	return results
}

func (d Results) Problems() Results {
	return d.filter(SProblem)
}

func (d Results) Undefined() Results {
	return d.filter(SUndefined)
}

func (d Results) OKs() Results {
	return d.filter(SOK)
}

func (d Results) Status() Status {
	if len(d.Problems()) > 0 {
		return SProblem
	}
	if len(d.Undefined()) > 0 {
		return SUndefined
	}
	return SOK
}
