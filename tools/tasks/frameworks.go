package tasks

type Framework struct {
	Active         bool     `json:"active"`
	Capabilities   []string `json:"capabilities"`
	Checkpoint     bool     `json:"checkpoint"`
	CompletedTasks []Task   `json:"completed_tasks"`
	Connected      bool     `json:"connected"`
	Executors      []struct {
		Command struct {
			Argv        []string `json:"argv"`
			Environment struct {
				Variables []struct {
					Name  string `json:"name"`
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"variables"`
			} `json:"environment"`
			Shell bool `json:"shell"`
			Uris  []struct {
				Executable bool   `json:"executable"`
				Value      string `json:"value"`
			} `json:"uris"`
			Value string `json:"value"`
		} `json:"command"`
		ExecutorID  string `json:"executor_id"`
		FrameworkID string `json:"framework_id"`
		Labels      []struct {
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
		Role    string `json:"role"`
		SlaveID string `json:"slave_id"`
		Type    string `json:"type"`
	} `json:"executors"`
	FailoverTimeout  float64 `json:"failover_timeout"`
	Hostname         string  `json:"hostname"`
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	OfferedResources struct {
		Cpus  float64 `json:"cpus"`
		Disk  float64 `json:"disk"`
		Gpus  float64 `json:"gpus"`
		Mem   float64 `json:"mem"`
		Ports string  `json:"ports"`
	} `json:"offered_resources"`
	Offers []struct {
		AllocationInfo struct {
			Role string `json:"role"`
		} `json:"allocation_info"`
		FrameworkID string `json:"framework_id"`
		ID          string `json:"id"`
		Resources   struct {
			Cpus  float64 `json:"cpus"`
			Disk  float64 `json:"disk"`
			Gpus  float64 `json:"gpus"`
			Mem   float64 `json:"mem"`
			Ports string  `json:"ports"`
		} `json:"resources"`
		SlaveID string `json:"slave_id"`
	} `json:"offers"`
	Pid              string  `json:"pid"`
	Principal        string  `json:"principal"`
	Recovered        bool    `json:"recovered"`
	RegisteredTime   float64 `json:"registered_time"`
	ReregisteredTime float64 `json:"reregistered_time"`
	Resources        struct {
		Cpus  float64 `json:"cpus"`
		Disk  float64 `json:"disk"`
		Gpus  float64 `json:"gpus"`
		Mem   float64 `json:"mem"`
		Ports string  `json:"ports"`
	} `json:"resources"`
	Role             string  `json:"role"`
	Tasks            []Task  `json:"tasks"`
	UnreachableTasks []Task  `json:"unreachable_tasks"`
	UnregisteredTime float64 `json:"unregistered_time"`
	UsedResources    struct {
		Cpus  float64 `json:"cpus"`
		Disk  float64 `json:"disk"`
		Gpus  float64 `json:"gpus"`
		Mem   float64 `json:"mem"`
		Ports string  `json:"ports"`
	} `json:"used_resources"`
	User     string `json:"user"`
	WebuiURL string `json:"webui_url"`
}

type Frameworks struct {
	CompletedFrameworks    []Framework `json:"completed_frameworks"`
	Frameworks             []Framework `json:"frameworks"`
	UnregisteredFrameworks []Framework `json:"unregistered_frameworks"`
}
