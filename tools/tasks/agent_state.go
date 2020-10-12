package tasks

type AgentTask struct {
	Container struct {
		Docker struct {
			ForcePullImage bool   `json:"force_pull_image"`
			Image          string `json:"image"`
			Network        string `json:"network"`
			Privileged     bool   `json:"privileged"`
		} `json:"docker"`
		LinuxInfo struct {
			SharePidNamespace bool `json:"share_pid_namespace"`
		} `json:"linux_info"`
		Mesos struct {
			Image struct {
				Cached bool `json:"cached"`
				Docker struct {
					Name string `json:"name"`
				} `json:"docker"`
				Type string `json:"type"`
			} `json:"image"`
		} `json:"mesos"`
		RlimitInfo struct {
			Rlimits []struct {
				Hard int64  `json:"hard"`
				Soft int64  `json:"soft"`
				Type string `json:"type"`
			} `json:"rlimits"`
		} `json:"rlimit_info"`
		Type    string `json:"type"`
		Volumes []struct {
			ContainerPath string `json:"container_path"`
			HostPath      string `json:"host_path"`
			Mode          string `json:"mode"`
			Source        struct {
				SandboxPath struct {
					Path string `json:"path"`
					Type string `json:"type"`
				} `json:"sandbox_path"`
				Type string `json:"type"`
			} `json:"source"`
		} `json:"volumes"`
	} `json:"container"`
	Discovery struct {
		Name  string `json:"name"`
		Ports struct {
			Ports []struct {
				Labels struct {
					Labels []struct {
						Key   string `json:"key"`
						Value string `json:"value"`
					} `json:"labels"`
				} `json:"labels"`
				Name       string `json:"name"`
				Number     int64  `json:"number"`
				Protocol   string `json:"protocol"`
				Visibility string `json:"visibility"`
			} `json:"ports"`
		} `json:"ports"`
		Visibility string `json:"visibility"`
	} `json:"discovery"`
	ExecutorID  string `json:"executor_id"`
	FrameworkID string `json:"framework_id"`
	ID          string `json:"id"`
	Labels      []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"labels"`
	Name      string `json:"name"`
	Resources struct {
		Cpus  int64  `json:"cpus"`
		Disk  int64  `json:"disk"`
		Gpus  int64  `json:"gpus"`
		Mem   int64  `json:"mem"`
		Ports string `json:"ports"`
	} `json:"resources"`
	Role     string `json:"role"`
	SlaveID  string `json:"slave_id"`
	State    string `json:"state"`
	Statuses []struct {
		ContainerStatus struct {
			ContainerID struct {
				Parent struct {
					Value string `json:"value"`
				} `json:"parent"`
				Value string `json:"value"`
			} `json:"container_id"`
			NetworkInfos []struct {
				IPAddresses []struct {
					IPAddress string `json:"ip_address"`
					Protocol  string `json:"protocol"`
				} `json:"ip_addresses"`
			} `json:"network_infos"`
		} `json:"container_status"`
		State     string  `json:"state"`
		Timestamp float64 `json:"timestamp"`
	} `json:"statuses"`
	User string `json:"user"`
}

type AgentState struct {
	CompletedFrameworks []struct {
		CompletedExecutors []struct {
			CompletedTasks []AgentTask `json:"completed_tasks"`
			Container      string      `json:"container"`
			Directory      string      `json:"directory"`
			ID             string      `json:"id"`
			Labels         []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"labels"`
			Name        string        `json:"name"`
			QueuedTasks []interface{} `json:"queued_tasks"`
			Resources   struct {
				Cpus float64 `json:"cpus"`
				Disk int64   `json:"disk"`
				Gpus int64   `json:"gpus"`
				Mem  int64   `json:"mem"`
			} `json:"resources"`
			Role   string      `json:"role"`
			Source string      `json:"source"`
			Tasks  []AgentTask `json:"tasks"`
			Type   string      `json:"type"`
		} `json:"completed_executors"`
		Executors []interface{} `json:"executors"`
		ID        string        `json:"id"`
		Name      string        `json:"name"`
	} `json:"completed_frameworks"`
	Frameworks []struct {
		CompletedExecutors []struct {
			CompletedTasks []AgentTask `json:"completed_tasks"`
			Container      string      `json:"container"`
			Directory      string      `json:"directory"`
			ID             string      `json:"id"`
			Labels         []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"labels"`
			Name        string        `json:"name"`
			QueuedTasks []interface{} `json:"queued_tasks"`
			Resources   struct {
				Cpus float64 `json:"cpus"`
				Disk int64   `json:"disk"`
				Gpus int64   `json:"gpus"`
				Mem  int64   `json:"mem"`
			} `json:"resources"`
			Role   string      `json:"role"`
			Source string      `json:"source"`
			Tasks  []AgentTask `json:"tasks"`
			Type   string      `json:"type"`
		} `json:"completed_executors"`
		Executors []struct {
			CompletedTasks []AgentTask `json:"completed_tasks"`
			Tasks          []AgentTask `json:"tasks"`
		} `json:"executors"`
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"frameworks"`
}
