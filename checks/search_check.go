package checks

import (
	_ "embed"
	"fmt"
	"path"
	"regexp"
	"sort"
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

func (c SearchCheck) checkFunc() CheckBundleFunc {
	builder := CheckFuncBuilder{}
	t := bundle.GetFileType(c.FileTypeName)
	for _, dirType := range t.DirTypes {
		switch dirType {
		case bundle.DTMaster:
			builder.CheckMasters = c.search
		case bundle.DTAgent:
			builder.CheckAgents = c.search
		case bundle.DTPublicAgent:
			builder.CheckPublicAgents = c.search
		}
	}
	builder.Aggregate = aggregate
	return builder.Build()
}

type problemValue struct {
	count int
	file  bundle.File
}

func aggregate(r Results) Results {
	var results Results = make([]Result, 0, len(r))
	problems := r.Problems()
	if len(problems) > 0 {
		if _, ok := problems[0].Value.(problemValue); ok {
			sort.Slice(problems, func(i int, j int) bool {
				countI := problems[i].Value.(problemValue).count
				countJ := problems[j].Value.(problemValue).count
				return countI > countJ
			})
			for i := range problems {
				problem := problems[i].Value.(problemValue)
				problems[i].Value = fmt.Sprintf("Error pattern occurred %v time(s) in file %s",
					problem.count, path.Base(problem.file.Name()))
			}
		}
	}
	results = append(results, problems...)
	results = append(results, r.Undefined()...)
	results = append(results, r.OKs()...)
	return results
}

func (c SearchCheck) search(host bundle.Host) Result {
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

	file, err := host.ScanLines(c.FileTypeName, f)
	if err != nil {
		return Result{
			Status: SUndefined,
			Host:   host,
			Value:  "Couldn't check. Error: " + err.Error(),
		}
	}
	if c.FailIfNotFound {
		if count == 0 {
			return Result{
				Status: SProblem,
				Value:  "Expected pattern not found in " + file.Name(),
			}
		}
	} else {
		if count > c.Max && lastN > lastNCure {
			return Result{
				Status: SProblem,
				Value:  problemValue{count, file},
			}
		}
	}
	return Result{
		Status: SOK,
	}
}

//go:embed search_checks.yaml
var searchChecksYAML []byte

func RegisterSearchChecks() {
	var searchChecks []SearchCheck
	err := yaml.Unmarshal(searchChecksYAML, &searchChecks)
	if err != nil {
		panic("Cannot read search checks YAML: " + err.Error())
	}
	for _, c := range searchChecks {
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
		c.Run = c.checkFunc()
		if c.FailIfNotFound {
			c.OKSummary = fmt.Sprintf("Expected pattern \"%s\" found.", c.ErrorPattern)
			c.ProblemSummary = fmt.Sprintf("Expected pattern \"%s\" not found.", c.ErrorPattern)
		} else {
			c.OKSummary = fmt.Sprintf("Error pattern \"%s\" not found.", c.ErrorPattern)
			c.ProblemSummary = fmt.Sprintf("Error pattern \"%s\" found.", c.ErrorPattern)
		}
		RegisterCheck(c.Check)
	}
}
