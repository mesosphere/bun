package checks

import (
	"github.com/mesosphere/bun/v2/bundle"
)

// CheckHost checks an individual host.
type CheckHost func(bundle.Host) Detail

// Aggregate aggregates heck results.
type Aggregate func(details Details) Details

// CheckFuncBuilder helps to create map/reduce-like checks.
type CheckFuncBuilder struct {
	CollectFromMasters      CheckHost // At least one of
	CollectFromAgents       CheckHost // the Collect... functions
	CollectFromPublicAgents CheckHost // are required
	Aggregate               Aggregate // Implement if the default is not sufficient
}

// Build returns a check function CheckFunc.Check
func (b *CheckFuncBuilder) Build() CheckFunc {
	if b.Aggregate == nil {
		b.Aggregate = DefaultAggregate
	}
	if b.CollectFromMasters == nil && b.CollectFromAgents == nil &&
		b.CollectFromPublicAgents == nil {
		panic("bun.CheckFuncBuilder.Build: At least one of the ForEach functions" +
			"should be specified")
	}
	return b.checkFunc
}

// Default implementation of the Aggregate function.
func DefaultAggregate(details Details) Details {
	return details
}

func (b *CheckFuncBuilder) checkHosts(h []bundle.Host, ch CheckHost, details *Details) {
	for _, host := range h {
		detail := ch(host)
		detail.Host = &host
		*details = append(*details, detail)
	}
}

// Implementation of the Check.CheckFunc
func (b *CheckFuncBuilder) checkFunc(bundle bundle.Bundle) (details Details) {
	details = make([]Detail, 0, len(bundle.Hosts))
	if b.CollectFromMasters != nil {
		b.checkHosts(bundle.Masters(), b.CollectFromMasters, &details)
	}
	if b.CollectFromAgents != nil {
		b.checkHosts(bundle.Agents(), b.CollectFromAgents, &details)
	}
	if b.CollectFromPublicAgents != nil {
		b.checkHosts(bundle.PublicAgents(), b.CollectFromPublicAgents, &details)
	}
	return b.Aggregate(details)
}
