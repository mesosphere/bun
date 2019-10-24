package checks

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mesosphere/bun/v2/bundle"
)

// SearchCheckBuilder builds a check which searches for the specified
// string in the the specified files. If the pattern
// is found more than Max times, the check is considered problematic.
// The amount of found lines appear in the Check.Problems of the check.
type SearchCheckBuilder struct {
	Name                 string              `yaml:"name"`                 // Required
	Description          string              `yaml:"description"`          // Optional
	Cure                 string              `yaml:"cure"`                 // Required
	FileTypeName         bundle.FileTypeName `yaml:"fileTypeName"`         // Required
	ErrorPattern         string              `yaml:"errorPattern"`         // Required
	CurePattern          string              `yaml:"curePattern"`          // Optional
	IsErrorPatternRegexp bool                `yaml:"isErrorPatternRegexp"` // Optional, default is false
	IsCurePatternRegexp  bool                `yaml:"isCurePatternRegexp"`  // Optional, default is false
	Max                  int                 `yaml:"max"`                  // Optional, default is 0
	errorRegexp          *regexp.Regexp
	cureRegexp           *regexp.Regexp
}

// Build creates a bun.Check.
func (b SearchCheckBuilder) Build() Check {
	if b.FileTypeName == "" {
		panic("FileTypeName should be specified.")
	}
	if b.ErrorPattern == "" {
		panic("ErrorPattern should be set.")
	}
	if b.Cure == "" {
		panic("Cure should be set.")
	}
	if b.IsErrorPatternRegexp {
		b.errorRegexp = regexp.MustCompile(b.ErrorPattern)
	}
	if b.IsCurePatternRegexp {
		b.cureRegexp = regexp.MustCompile(b.CurePattern)
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
	var count int
	var lastN int
	var lastNCure int
	var matchError func(line string) bool
	if b.IsErrorPatternRegexp {
		matchError = func(line string) bool {
			return b.errorRegexp.MatchString(line)
		}
	} else {
		matchError = func(line string) bool {
			return strings.Contains(line, b.ErrorPattern)
		}
	}
	var matchCure func(line string) bool
	if b.CurePattern == "" {
		matchCure = func(line string) bool {
			return false
		}
	} else if !b.IsCurePatternRegexp {
		matchCure = func(line string) bool {
			return strings.Contains(line, b.CurePattern)
		}
	} else {
		matchCure = func(line string) bool {
			return b.cureRegexp.MatchString(line)
		}
	}
	f := func(n int, line string) bool {
		if matchError(line) {
			count++
			lastN = n
		}
		if matchCure(line) {
			lastNCure = n
		}
		return false
	}
	err = host.ScanLines(b.FileTypeName, f)
	if count > b.Max && lastN > lastNCure {
		details = fmt.Sprintf("%v problems occurred in the logs.", count)
		return
	}
	details = "No problems were found"
	ok = true
	return
}
