package checks

import (
	"github.com/mesosphere/bun/v2/bundle"
)

// CheckFuncBuilder helps to create map/reduce-like checks.
type CheckFuncBuilder struct {
	CheckMasters      CheckHostFunc // At least one of
	CheckAgents       CheckHostFunc // the Collect... functions
	CheckPublicAgents CheckHostFunc // are required
	Aggregate         Aggregate     // Implement if the default is not sufficient
}

// CheckHostFunc checks an individual host.
type CheckHostFunc func(bundle.Host) Result

// Aggregate aggregates check results produced by CheckMaster, CheckAgents, and CheckPublicAgents functions.
type Aggregate func(results Results) Results

// Build returns a check function Run.Check
func (b CheckFuncBuilder) Build() CheckBundleFunc {
	if b.Aggregate == nil {
		b.Aggregate = DefaultAggregate
	}
	if b.CheckMasters == nil && b.CheckAgents == nil &&
		b.CheckPublicAgents == nil {
		panic("bun.CheckFuncBuilder.Build: At least one of the CheckHostFunc functions" +
			"should be specified")
	}
	return b.checkFunc
}

// Default implementation of the Aggregate function.
func DefaultAggregate(results Results) Results {
	return results
}

func checkHosts(hosts []bundle.Host, checkHostFunc CheckHostFunc, results *Results) {
	for _, host := range hosts {
		result := checkHostFunc(host)
		result.Host = host
		*results = append(*results, result)
	}
}

// Implementation of the Check.Run
func (b CheckFuncBuilder) checkFunc(bundle bundle.Bundle) (results Results) {
	results = make([]Result, 0, len(bundle.Hosts))
	if b.CheckMasters != nil {
		checkHosts(bundle.Masters(), b.CheckMasters, &results)
	}
	if b.CheckAgents != nil {
		checkHosts(bundle.Agents(), b.CheckAgents, &results)
	}
	if b.CheckPublicAgents != nil {
		checkHosts(bundle.PublicAgents(), b.CheckPublicAgents, &results)
	}
	return b.Aggregate(results)
}
