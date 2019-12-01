package checks

import (
	"fmt"
	"strings"
	"sync"
)

var (
	checkRegistry   = make(map[string]Check)
	checkRegistryMu sync.RWMutex
)

// RegisterCheck registers a new check to make it discoverable for consumers.
func RegisterCheck(c Check) {
	checkRegistryMu.Lock()
	defer checkRegistryMu.Unlock()
	if c.Name == "" {
		panic("bun.checks.RegisterCheck: Check Name should not be empty")
	}
	if c.Description == "" {
		panic("bun.checks.RegisterCheck: Check Description should not be empty")
	}
	checkS(c)
	if c.Cure == "" {
		panic("bun.checks.RegisterCheck: Check Cure should not be empty")
	}
	if _, exists := checkRegistry[c.Name]; exists {
		panic(fmt.Sprintf("bun.RegisterCheck: called twice for check %v", c.Name))
	}
	if c.ProblemSummary == "" {
		c.ProblemSummary = "Problems were found."
	}
	if c.OKSummary == "" {
		c.OKSummary = "No problems were found."
	}
	checkRegistry[c.Name] = c
}

func checkS(c Check) {
	verb := strings.Fields(c.Description)[0]
	if !strings.HasSuffix(verb, "s") {
		panic(fmt.Sprintf("Wrong description for the check \"%s\"."+
			" Check description should start with a present tense third person singular verb.", c.Name))
	}
}

// Checks Returns all registered checks.
func Checks() []Check {
	checkRegistryMu.RLock()
	defer checkRegistryMu.RUnlock()
	checks := make([]Check, 0, len(checkRegistry))
	for _, c := range checkRegistry {
		checks = append(checks, c)
	}
	return checks
}

// GetCheck returns check by name.
func GetCheck(name string) Check {
	checkRegistryMu.RLock()
	check, ok := checkRegistry[name]
	checkRegistryMu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("bun.GetCheck: don't have check %v", name))
	}
	return check
}
