package errors

import (
	"bufio"
	"fmt"
	"regexp"

	"github.com/adyatlov/bun"
)

// http://erlang.org/doc/apps/kernel/logger_chapter.html#log_level
const ErrorRegExp = `\[(?P<Level>error|emergency|critical|alert)\]`

func init() {
	builder := bun.CheckBuilder{
		Name:                    "networking-errors",
		Description:             "Identify errors in dcos-net logs",
		CollectFromMasters:      collect,
		CollectFromAgents:       collect,
		CollectFromPublicAgents: collect,
		Aggregate:               aggregate,
	}
	check := builder.Build()
	bun.RegisterCheck(check)
}

func collect(host bun.Host) (ok bool, details interface{}, err error) {
	ok = true
	file, err := host.OpenFile("net")

	if err != nil {
		ok = false
		erMsg := fmt.Sprintf("Cannot open net file %s", err)
		fmt.Println(errMsg)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	pattern := regexp.MustCompile(ErrorRegExp)
	logs := map[string]int{}
	for scanner.Scan() {
		line := scanner.Text()
		if pattern.MatchString(line) {
			ok = false
			res := pattern.FindStringSubmatch(line)
			logs[res[1]]++
		}
	}

	details = logs
	return
}

func aggregate(c *bun.Check, b bun.CheckBuilder) {
	details := []string{}

	if(len(b.Problems) == 0){
		c.OKs = details
		c.Summary = fmt.Sprintf("No dcos-net _errors_ logs found.")
		return
	}

	for _, r := range b.Problems {
		logs := r.Details.(map[string]int)
		for level, count := range logs {
			details = append(details, fmt.Sprintf("%d %v logs on %v %v",
					count, level, r.Host.Type, r.Host.IP))
		}
	}

	c.Summary = fmt.Sprintf(
		"Found dcos-net _errors_ logs on %d out of %d nodes.",
		len(details), len(b.OKs))
	c.Problems = details
}
