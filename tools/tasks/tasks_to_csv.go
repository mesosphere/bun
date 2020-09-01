package tasks

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/mesosphere/bun/v2/bundle"
)

func ToCSV(b *bundle.Bundle, writer io.Writer) error {
	var err error
	b.ForEachFile("mesos-master-frameworks",
		func(f bundle.File) (stop bool) {
			defer func() { _ = f.Close() }()
			decoder := json.NewDecoder(f)
			var frameworks Frameworks
			err = decoder.Decode(&frameworks)
			if err != nil {
				return false
			}
			w := csv.NewWriter(writer)
			err = w.Write([]string{"Framework name", "Framework ID", "Framework Active", "Framework Status",
				"Name", "ID", "Launched (UTC)", "Finished (UTC)", "Running", "CPUs", "Memory", "IPs"})
			if err != nil {
				return true
			}
			frameworksSet := map[string][]Framework{
				"registered":   frameworks.Frameworks,
				"completed":    frameworks.CompletedFrameworks,
				"unregistered": frameworks.UnregisteredFrameworks,
			}
			for status, frameworks := range frameworksSet {
				err = writeFrameworks(w, frameworks, status)
				if err != nil {
					return true
				}
			}
			w.Flush()
			return true
		},
	)
	return err
}

func writeFrameworks(w *csv.Writer, frameworks []Framework, status string) error {
	for _, framework := range frameworks {
		frameworkPart := []string{framework.Name, framework.ID,
			strconv.FormatBool(framework.Active), status}
		tasksSet := [][]Task{framework.Tasks, framework.UnreachableTasks, framework.CompletedTasks}
		for _, tasks := range tasksSet {
			for _, task := range tasks {
				line := append(frameworkPart, taskLines(&task)...)
				if err := w.Write(line); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func taskLines(task *Task) []string {
	launched := "N/A"
	finished := "N/A"
	running := "true"
	for _, status := range task.Statuses {
		if status.State == "TASK_STARTING" {
			launched = secondsToSQLTime(status.Timestamp)
			continue
		}
		if status.State == "TASK_FAILED" ||
			status.State == "TASK_FINISHED" ||
			status.State == "TASK_KILLED" {
			finished = secondsToSQLTime(status.Timestamp)
			running = "false"
			continue
		}
	}
	cpus := strconv.FormatFloat(task.Resources.Cpus, 'f', -1, 64)
	mem := strconv.FormatFloat(task.Resources.Mem, 'f', -1, 64)
	ips := findTaskIPs(task)
	return []string{task.Name, task.ID, launched, finished, running, cpus, mem, ips}
}

func secondsToSQLTime(seconds float64) string {
	t := time.Unix(int64(seconds), 0)
	return t.UTC().Format("2006-01-02 15:04:05")
}

func findTaskIPs(task *Task) string {
	ips := make(map[string]interface{})
	for _, statuses := range task.Statuses {
		for _, networkInfo := range statuses.ContainerStatus.NetworkInfos {
			for _, ipAddress := range networkInfo.IPAddresses {
				ips[ipAddress.IPAddress] = nil
			}
		}
	}
	b := make([]string, 0, len(ips))
	for ip, _ := range ips {
		b = append(b, ip)
	}
	return strings.Join(b, " ")
}
