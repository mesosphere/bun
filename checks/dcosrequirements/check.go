// Package dcosrequirements implements checks for DC/OS system requirements.
package dcosrequirements

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lithammer/dedent"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

// Requirement values are taken from the dcos-2.0 documentation,
// https://docs.d2iq.com/mesosphere/dcos/2.0/installing/production/system-requirements/.
var (
	cpuRequirements = map[bundle.DirType]int{
		bundle.DTMaster:      4,
		bundle.DTAgent:       2,
		bundle.DTPublicAgent: 2,
	}

	memRequirements = map[bundle.DirType]int{
		bundle.DTMaster:      32_000_000,
		bundle.DTAgent:       60_000_000,
		bundle.DTPublicAgent: 60_000_000,
	}

	diskRequirements = map[bundle.DirType][]uint64{
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
		Name:        "dcos-requirements",
		Description: "Checks that nodes meet system requirements",
		Cure: dedent.Dedent(`Adjust the node spec to meet system requirements: 
			https://docs.d2iq.com/mesosphere/dcos/2.0/installing/production/system-requirements/`),
		OKSummary:      "Nodes meet system requirements",
		ProblemSummary: "Nodes do not meet system requirements",
		Run:            builder.Build(),
	}

	checks.RegisterCheck(check)
}

func checkCpus(host bundle.Host) []error {
	cpuinfo, err := host.OpenFile("cpuinfo")

	if err != nil {
		return []error{err}
	}

	defer cpuinfo.Close()

	scanner := bufio.NewScanner(cpuinfo)

	numCpus := 0

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "processor\t:") {
			numCpus++
		}
	}

	if numCpus < cpuRequirements[host.Type] {
		return []error{
			fmt.Errorf(
				"node has less than required CPUs: %.f%% (%d vs. %d)",
				100*float64(numCpus)/float64(cpuRequirements[host.Type]),
				numCpus,
				cpuRequirements[host.Type])}
	}

	return nil
}

func checkMem(host bundle.Host) []error {
	meminfo, err := host.OpenFile("meminfo")

	if err != nil {
		return []error{err}
	}

	defer meminfo.Close()

	scanner := bufio.NewScanner(meminfo)

	mem := 0

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "MemTotal") {
			tokens := strings.Split(scanner.Text(), " ")
			mem, err = strconv.Atoi(tokens[len(tokens)-2])

			if err != nil {
				return []error{err}
			}

			break
		}
	}

	if mem < memRequirements[host.Type] {
		return []error{
			fmt.Errorf(
				"node has less than required memory: %.f%% (%d kB vs. %d kB)",
				100*float64(mem)/float64(memRequirements[host.Type]),
				mem,
				memRequirements[host.Type])}
	}

	return nil
}

type disk struct {
	ID    string
	Size  uint64
	Used  uint64
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

		size, err := strconv.ParseUint(fields[1], 10, 64)

		if err != nil {
			return nil, err
		}

		used, err := strconv.ParseUint(fields[2], 10, 64)

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

func getMountpoint(dir string, disks []disk) (*disk, error) {
	var candidate *disk = nil

	for _, disk := range disks {
		disk := disk

		rel, err := filepath.Rel(disk.Mount, dir)

		if err != nil {
			return nil, err
		}

		if !strings.HasPrefix(rel, ".") && !strings.HasPrefix(rel, "/") {
			if candidate != nil &&
				len(filepath.SplitList(rel)) <= len(filepath.SplitList(candidate.Mount)) {
				continue
			}

			candidate = &disk
		}

		if candidate == nil {
			return nil, fmt.Errorf("could not determine mount point of directory '%s'", dir)
		}
	}

	return candidate, nil
}

func checkDisksMesos(host bundle.Host) []error {
	errs := []error{}

	disks, err := getDisks(host)

	if err != nil {
		// Ignore errors due to the absence of the `df` file since it is not present in all versions.
		return nil
	}

	type Flags struct {
		WorkDir    string `json:"work_dir"`
		RuntimeDir string `json:"runtime_dir"`
	}

	type MesosFlags struct {
		Flags Flags `json:"flags"`
	}

	flagsFile := map[bundle.DirType]bundle.FileTypeName{
		bundle.DTMaster:      "mesos-master-flags",
		bundle.DTAgent:       "mesos-agent-flags",
		bundle.DTPublicAgent: "mesos-agent-flags",
	}

	var flags MesosFlags
	if err = host.ReadJSON(flagsFile[host.Type], &flags); err != nil {
		errs = append(errs, err)
		return errs
	}

	workDirDisk, err := getMountpoint(flags.Flags.WorkDir, disks)

	if err != nil {
		errs = append(errs, err)
		return errs
	}

	if workDirDisk.Size < diskRequirements[host.Type][0] {
		errs = append(errs,
			fmt.Errorf("node has less than required disk for Mesos 'work_dir': %d kB vs. %d kB",
				workDirDisk.Size,
				diskRequirements[host.Type][0]))
	}

	if len(flags.Flags.RuntimeDir) > 0 {
		runtimeDirDisk, err := getMountpoint(flags.Flags.RuntimeDir, disks)

		if err != nil {
			errs = append(errs, err)
			return errs
		}

		if runtimeDirDisk.Size < diskRequirements[host.Type][1] {
			errs = append(errs,
				fmt.Errorf("node has less than required disk for Mesos 'runtime_dir': %d kB vs. %d kB",
					runtimeDirDisk.Size,
					diskRequirements[host.Type][1]))
		}
	}

	return errs
}

func checkDisk(host bundle.Host) []error {
	errs := checkDisksMesos(host)

	return errs
}

func collect(host bundle.Host) checks.Result {
	errs := []error{}

	if err := checkCpus(host); len(err) > 0 {
		errs = append(errs, err...)
	}

	if err := checkMem(host); len(err) > 0 {
		errs = append(errs, err...)
	}

	if err := checkDisk(host); len(err) > 0 {
		errs = append(errs, err...)
	}

	if len(errs) > 0 {
		messages := make([]string, len(errs))

		for i, err := range errs {
			messages[i] = err.Error()
		}

		return checks.Result{
			Status: checks.SProblem,
			Value:  strings.Join(messages, "\n"),
		}
	}

	return checks.Result{
		Status: checks.SOK,
	}
}
