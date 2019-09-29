package checks

import yaml "gopkg.in/yaml.v2"

func RegisterSearchChecks() {
	var searchChecks []SearchCheckBuilder
	err := yaml.Unmarshal([]byte(searchChecksYAML), &searchChecks)
	if err != nil {
		panic("Cannot read search checks YAML: " + err.Error())
	}
	for _, builder := range searchChecks {
		check := builder.Build()
		RegisterCheck(check)
	}
}
