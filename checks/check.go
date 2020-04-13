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

func (r Results) filter(status Status) Results {
	results := make([]Result, 0, len(r))
	for _, result := range r {
		if result.Status == status {
			results = append(results, result)
		}
	}
	return results
}

func (r Results) Problems() Results {
	return r.filter(SProblem)
}

func (r Results) Undefined() Results {
	return r.filter(SUndefined)
}

func (r Results) OKs() Results {
	return r.filter(SOK)
}

func (r Results) Status() Status {
	if len(r.Problems()) > 0 {
		return SProblem
	}
	if len(r.Undefined()) > 0 {
		return SUndefined
	}
	return SOK
}
