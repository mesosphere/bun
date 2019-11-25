package checks

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/mesosphere/bun/v2/bundle"
)

// SearchCheck is a check which searches for the specified
// string in the the specified files. If the pattern
// is found more than Max times, the check is considered problematic.
type SearchCheck struct {
	Check                `yaml:",inline"`
	FileTypeName         bundle.FileTypeName `yaml:"fileTypeName"`         // Required
	ErrorPattern         string              `yaml:"errorPattern"`         // Required
	CurePattern          string              `yaml:"curePattern"`          // Optional
	IsErrorPatternRegexp bool                `yaml:"isErrorPatternRegexp"` // Optional, default is false
	IsCurePatternRegexp  bool                `yaml:"isCurePatternRegexp"`  // Optional, default is false
	Max                  int                 `yaml:"max"`                  // Optional, default is 0
	FailIfNotFound       bool                `yaml:"failIfNotFound"`       // Optional, default false
	errorRegexp          *regexp.Regexp
	cureRegexp           *regexp.Regexp
}

func (c SearchCheck) checkFunc() CheckFunc {
	if c.FileTypeName == "" {
		panic("FileTypeName should be specified.")
	}
	if c.ErrorPattern == "" {
		panic("ErrorPattern should be set.")
	}
	if c.FailIfNotFound && c.CurePattern != "" {
		panic("FailIfNotFound and CurePattern are mutually exclusive")
	}
	if c.IsErrorPatternRegexp {
		c.errorRegexp = regexp.MustCompile(c.ErrorPattern)
	}
	if c.IsCurePatternRegexp {
		c.cureRegexp = regexp.MustCompile(c.CurePattern)
	}
	builder := CheckFuncBuilder{}
	t := bundle.GetFileType(c.FileTypeName)
	for _, dirType := range t.DirTypes {
		switch dirType {
		case bundle.DTMaster:
			builder.CollectFromMasters = c.collect
		case bundle.DTAgent:
			builder.CollectFromAgents = c.collect
		case bundle.DTPublicAgent:
			builder.CollectFromPublicAgents = c.collect
		}
	}
	return builder.Build()
}

func (c SearchCheck) collect(host bundle.Host) Detail {
	var count int
	var lastN int
	var lastNCure int
	var matchError func(line string) bool
	if c.IsErrorPatternRegexp {
		matchError = func(line string) bool {
			return c.errorRegexp.MatchString(line)
		}
	} else {
		matchError = func(line string) bool {
			return strings.Contains(line, c.ErrorPattern)
		}
	}
	var matchCure func(line string) bool
	if c.CurePattern == "" {
		matchCure = func(line string) bool {
			return false
		}
	} else if !c.IsCurePatternRegexp {
		matchCure = func(line string) bool {
			return strings.Contains(line, c.CurePattern)
		}
	} else {
		matchCure = func(line string) bool {
			return c.cureRegexp.MatchString(line)
		}
	}
	f := func(n int, line string) bool {
		if matchError(line) {
			count++
			lastN = n
			if c.FailIfNotFound {
				return true
			}
		}
		if matchCure(line) {
			lastNCure = n
		}
		return false
	}
	err := host.ScanLines(c.FileTypeName, f)
	if c.FailIfNotFound {
		if count == 0 {
			return Detail{
				Status: SProblem,
				Err:    err,
			}
		}
	} else {
		if count > c.Max && lastN > lastNCure {
			return Detail{
				Status: SProblem,
				Value:  count,
				Err:    err,
			}
		}
	}
	return Detail{
		Status: SOK,
		Err:    err,
	}
}

func RegisterSearchChecks() {
	var searchChecks []SearchCheck
	err := yaml.Unmarshal([]byte(searchChecksYAML), &searchChecks)
	if err != nil {
		panic("Cannot read search checks YAML: " + err.Error())
	}
	for _, c := range searchChecks {
		c.CheckFunc = c.checkFunc()
		if c.FailIfNotFound {
			c.OKSummary = fmt.Sprintf("Expected pattern \"%s\" was found in file.", c.ErrorPattern)
			c.ProblemSummary = fmt.Sprintf("Expected pattern \"%s\" was not found.", c.ErrorPattern)
		} else {
			c.OKSummary = fmt.Sprintf("Error pattern \"%s\" was not found.", c.ErrorPattern)
			c.ProblemSummary = fmt.Sprintf("Error pattern \"%s\" was found.", c.ErrorPattern)
		}
		RegisterCheck(c.Check)
	}
}
