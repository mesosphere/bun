package tasks

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/hako/durafmt"

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
				"Name", "ID", "State", "Launched (UTC)", "Finished (UTC)", "Duration", "Duration (seconds)", "Container Type", "Running", "CPUs", "Memory", "IPs"})
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

func isStateTerminal(state string) bool {
	return state == "TASK_FAILED" ||
		state == "TASK_FINISHED" ||
		state == "TASK_KILLED" ||
		state == "TASK_LOST" ||
		state == "TASK_DROPPED" ||
		state == "TASK_GONE" ||
		state == "TASK_GONE_BY_OPERATOR"
}

func taskLines(task *Task) []string {
	var launched float64
	var finished float64
	running := "true"
	state := "N/A"
	containerType := "N/A"
	if task.State != "" {
		state = task.State
	}
	for _, status := range task.Statuses {
		if status.State == "TASK_STARTING" {
			launched = status.Timestamp
			continue
		}
		if isStateTerminal(status.State) {
			finished = status.Timestamp
			running = "false"
			continue
		}
	}
	duration, durationS := duration(launched, finished)
	if task.Container.Type != "" {
		containerType = task.Container.Type
	}
	cpus := strconv.FormatFloat(task.Resources.Cpus, 'f', -1, 64)
	mem := strconv.FormatFloat(task.Resources.Mem, 'f', -1, 64)
	ips := findTaskIPs(task)
	return []string{task.Name, task.ID, state,
		secondsToSQLTime(launched), secondsToSQLTime(finished),
		duration, durationS,
		containerType,
		running, cpus, mem, ips}
}

func duration(launched float64, finished float64) (string, string) {
	if launched == 0.0 || finished == 0.0 {
		return "N/A", "N/A"
	}
	begin := time.Unix(int64(launched), 0)
	end := time.Unix(int64(finished), 0)
	d := end.Sub(begin)
	return durafmt.Parse(d).String(), strconv.Itoa(int(d.Round(time.Second).Seconds()))
}

func secondsToSQLTime(seconds float64) string {
	if seconds == 0.0 {
		return "N/A"
	}
	t := time.Unix(int64(seconds), 0).UTC()
	return t.Format("2006-01-02 15:04:05")
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
	for ip := range ips {
		b = append(b, ip)
	}
	return strings.Join(b, " ")
}
