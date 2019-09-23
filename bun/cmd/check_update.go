package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-version"
	"net/http"
	"time"
)

func CheckUpdate() {
	thisVersion, err := version.NewSemver("v" + Version)
	if err != nil {
		return
	}
	gh := gitHub{"mesosphere", "bun"}
	latestVersion, err := gh.latestVersion()
	if err != nil {
		fmt.Println(err)
		return
	}
	if thisVersion.LessThan(latestVersion) {
		fmt.Printf("Good news, everyone! New Bun is out, please download v%s from https://github.com/mesosphere/bun/releases/latest\n\n",
			latestVersion.String())
	}
}

type gitHub struct {
	Account string
	Repo    string
}

type tag struct {
	Name string
}

func (g gitHub) latestVersion() (*version.Version, error) {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/" + g.Account + "/" + g.Repo + "/tags")
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(resp.Body)
	tags := make([]tag, 0)
	if err := decoder.Decode(&tags); err != nil {
		return nil, err
	}
	latest := version.Must(version.NewSemver("v0.0.0"))
	for _, tag := range tags {
		v, err := version.NewSemver(tag.Name)
		if err != nil {
			continue
		}
		if latest.LessThan(v) {
			latest = v
		}
	}
	return latest, nil
}
