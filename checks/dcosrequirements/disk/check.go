package disk

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

// Requirement values are taken from the dcos-2.0 documentation,
// https://docs.d2iq.com/mesosphere/dcos/2.0/installing/production/system-requirements/.
var (
	diskRequirements = map[bundle.DirType][]int{
		bundle.DTMaster:      {120_000_000, 120_000_000},
		bundle.DTAgent:       {60_000_000, 60_000_000},
		bundle.DTPublicAgent: {60_000_000, 60_000_000},
	}
)

func init() {
	builder := checks.CheckFuncBuilder{
		CheckMasters:      collect,
		CheckAgents:       collect,
		CheckPublicAgents: collect,
	}
	check := checks.Check{
		Name:        "disk",
		Description: "Checks that DC/OS nodes meet disk space requirements",
		Cure: "Add disk space to the problematic node to meet system requirements: " +
			"https://docs.d2iq.com/mesosphere/dcos/2.0/installing/production/system-requirements/",
		OKSummary:      "All nodes meet disk space requirements",
		ProblemSummary: "Some nodes do not meet disk space requirements",
		Run:            builder.Build(),
	}
	checks.RegisterCheck(check)
}

type disk struct {
	ID    string
	Size  int
	Used  int
	Mount string
}

func getDisks(host bundle.Host) ([]disk, error) {
	df, err := host.OpenFile("df")
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(df)
	var disks []disk
	headerParsed := false
	for scanner.Scan() {
		if !headerParsed {
			headerParsed = true
			continue
		}
		fields := strings.Fields(scanner.Text())
		size, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		used, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, err
		}
		disks = append(disks,
			disk{
				ID:    fields[0],
				Size:  size,
				Used:  used,
				Mount: fields[5],
			})
	}
	return disks, nil
}

func getMountpoint(dir string, disks []disk) (disk, error) {
	var candidate disk
	found := false
	for _, d := range disks {
		rel, err := filepath.Rel(d.Mount, dir)
		if err != nil {
			return candidate, err
		}
		if !strings.HasPrefix(rel, ".") && !strings.HasPrefix(rel, "/") {
			if found &&
				len(filepath.SplitList(rel)) <= len(filepath.SplitList(candidate.Mount)) {
				continue
			}
			candidate = d
			found = true
		}
	}
	if !found {
		return candidate, fmt.Errorf("could not determine mount point of directory '%s'", dir)
	}
	return candidate, nil
}

func checkWorkDir(host bundle.Host, disks []disk, dir string) checks.Result {
	workDirDisk, err := getMountpoint(dir, disks)
	if err != nil {
		return checks.Result{
			Status: checks.SUndefined,
			Host:   host,
			Value:  "Couldn't get Mountpoint wile checking work disk requirement: " + err.Error(),
		}
	}
	if workDirDisk.Size < diskRequirements[host.Type][0] {
		return checks.Result{
			Status: checks.SProblem,
			Host:   host,
			Value: fmt.Sprintf("node has less than required disk for Mesos 'work_dir': %.f%% (%.2f GB vs. %.2f GB)",
				100*float64(workDirDisk.Size)/float64(diskRequirements[host.Type][0]),
				convertKBtoGB(workDirDisk.Size),
				convertKBtoGB(diskRequirements[host.Type][0])),
		}
	}
	return checks.Result{
		Status: checks.SOK,
		Host:   host,
	}
}

func checkRuntimeDir(host bundle.Host, disks []disk, dir string) checks.Result {
	runtimeDirDisk, err := getMountpoint(dir, disks)
	if err != nil {
		return checks.Result{
			Status: checks.SUndefined,
			Host:   host,
			Value:  "Couldn't get Mountpoint wile checking runtime disk requirement: " + err.Error(),
		}
	}
	if runtimeDirDisk.Size < diskRequirements[host.Type][1] {
		return checks.Result{
			Status: checks.SProblem,
			Host:   host,
			Value: fmt.Sprintf("node has less than required disk for Mesos 'runtime_dir': %.f%% (%.2f GB vs. %.2f GB)",
				100*float64(runtimeDirDisk.Size)/float64(diskRequirements[host.Type][1]),
				convertKBtoGB(runtimeDirDisk.Size),
				convertKBtoGB(diskRequirements[host.Type][1])),
		}
	}
	return checks.Result{
		Status: checks.SOK,
		Host:   host,
	}
}

func collect(host bundle.Host) checks.Result {
	type Flags struct {
		WorkDir    string `json:"work_dir"`
		RuntimeDir string `json:"runtime_dir"`
	}
	type MesosFlags struct {
		Flags Flags `json:"flags"`
	}

	disks, err := getDisks(host)
	if err != nil {
		return checks.Result{
			Status: checks.SUndefined,
			Host:   host,
			Value:  "Couldn't check disk requirement: " + err.Error(),
		}
	}
	flagsFile := map[bundle.DirType]bundle.FileTypeName{
		bundle.DTMaster:      "mesos-master-flags",
		bundle.DTAgent:       "mesos-agent-flags",
		bundle.DTPublicAgent: "mesos-agent-flags",
	}
	var flags MesosFlags
	if err = host.ReadJSON(flagsFile[host.Type], &flags); err != nil {
		return checks.Result{
			Status: checks.SUndefined,
			Host:   host,
			Value:  "Couldn't read JSON while checking disk requirement: " + err.Error(),
		}
	}
	var results checks.Results
	results = append(results, checkWorkDir(host, disks, flags.Flags.WorkDir))
	if len(flags.Flags.RuntimeDir) > 0 {
		results = append(results, checkRuntimeDir(host, disks, flags.Flags.RuntimeDir))
	}
	result := checks.Result{
		Host:   host,
		Status: results.Status(),
	}
	if result.Status == checks.SOK {
		return result
	}
	messages := make([]string, 0)
	for _, p := range results.Problems() {
		messages = append(messages, p.Value.(string))
	}
	for _, u := range results.Undefined() {
		messages = append(messages, u.Value.(string))
	}
	return checks.Result{
		Value: strings.Join(messages, "\n"),
	}
}

func convertKBtoGB(kb int) float64 {
	return float64(kb) / 1e6
}
