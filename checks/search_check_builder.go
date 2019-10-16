package checks

import (
	"fmt"
	"github.com/mesosphere/bun/v2/bundle"
	"regexp"
)

// SearchCheckBuilder builds a check which searches for the specified
// string in the the specified files. If the pattern
// is found more than Max times, the check is considered problematic.
// The amount of found lines appear in the Check.Problems of the check.
type SearchCheckBuilder struct {
	Name         string              `yaml:"name"`         // Required
	Description  string              `yaml:"description"`  // Optional
	Cure         string              `yaml:"cure"`         // Optional
	FileTypeName bundle.FileTypeName `yaml:"fileTypeName"` // Required
	SearchString string              `yaml:"searchString"` // Required
	IsRegexp     bool                `yaml:"isRegexp"`     // Optional, default is false
	Max          uint                `yaml:"max"`          // Optional, default is 0
	regexp       *regexp.Regexp
}

// Build creates a bun.Check.
func (b SearchCheckBuilder) Build() Check {
	if b.FileTypeName == "" {
		panic("FileTypeName should be specified.")
	}
	if b.SearchString == "" {
		panic("SearchString should be set.")
	}
	if b.IsRegexp {
		b.regexp = regexp.MustCompile(b.SearchString)
	}
	builder := CheckBuilder{
		Name:        b.Name,
		Description: b.Description,
		Cure:        b.Cure,
		Aggregate:   DefaultAggregate,
	}
	t := bundle.GetFileType(b.FileTypeName)
	for _, dirType := range t.DirTypes {
		switch dirType {
		case bundle.DTMaster:
			builder.CollectFromMasters = b.collect
		case bundle.DTAgent:
			builder.CollectFromAgents = b.collect
		case bundle.DTPublicAgent:
			builder.CollectFromPublicAgents = b.collect
		}
	}
	return builder.Build()
}

func (b SearchCheckBuilder) collect(host bundle.Host) (ok bool, details interface{}, err error) {
	var n uint
	if b.IsRegexp {
		n, err = host.FindLineRegexp(b.FileTypeName, b.regexp)
	} else {
		n, err = host.FindLine(b.FileTypeName, b.SearchString)
	}
	if err != nil {
		return
	}
	if n > b.Max {
		details = fmt.Sprintf("%v problems occurred in the logs.", n)
		return
	}
	ok = true
	return
}
