package offeredresources

type Task struct {
	Container struct {
		Docker struct {
			ForcePullImage bool   `json:"force_pull_image"`
			Image          string `json:"image"`
			Network        string `json:"network"`
			Parameters     []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"parameters"`
			PortMappings []struct {
				ContainerPort int64  `json:"container_port"`
				HostPort      int64  `json:"host_port"`
				Protocol      string `json:"protocol"`
			} `json:"port_mappings"`
			Privileged bool `json:"privileged"`
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
		NetworkInfos []struct {
			IPAddresses []struct {
				Protocol string `json:"protocol"`
			} `json:"ip_addresses"`
			Labels struct {
				Labels []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"labels"`
			} `json:"labels"`
			Name string `json:"name"`
		} `json:"network_infos"`
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
				Secret struct {
					Reference struct {
						Name string `json:"name"`
					} `json:"reference"`
					Type string `json:"type"`
				} `json:"secret"`
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
	HealthCheck struct {
		Command struct {
			Shell bool   `json:"shell"`
			Value string `json:"value"`
		} `json:"command"`
		ConsecutiveFailures int64   `json:"consecutive_failures"`
		DelaySeconds        float64 `json:"delay_seconds"`
		GracePeriodSeconds  float64 `json:"grace_period_seconds"`
		HTTP                struct {
			Path     string `json:"path"`
			Port     int64  `json:"port"`
			Protocol string `json:"protocol"`
			Scheme   string `json:"scheme"`
		} `json:"http"`
		IntervalSeconds float64 `json:"interval_seconds"`
		Tcp             struct {
			Port     int64  `json:"port"`
			Protocol string `json:"protocol"`
		} `json:"tcp"`
		TimeoutSeconds float64 `json:"timeout_seconds"`
		Type           string  `json:"type"`
	} `json:"health_check"`
	ID     string `json:"id"`
	Labels []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"labels"`
	Name      string `json:"name"`
	Resources struct {
		Cpus  float64 `json:"cpus"`
		Disk  float64 `json:"disk"`
		Gpus  float64 `json:"gpus"`
		Mem   float64 `json:"mem"`
		Ports string  `json:"ports"`
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
				Labels []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"labels"`
				Name string `json:"name"`
			} `json:"network_infos"`
		} `json:"container_status"`
		Healthy bool `json:"healthy"`
		Labels  []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"labels"`
		State     string  `json:"state"`
		Timestamp float64 `json:"timestamp"`
	} `json:"statuses"`
	User string `json:"user"`
}
