package overlay

type overlay struct {
	ConfigurationAttempts int64  `json:"configuration_attempts"`
	IP                    string `json:"ip"`
	Overlays              []struct {
		Backend struct {
			Vxlan struct {
				Vni      int64  `json:"vni"`
				VtepIP   string `json:"vtep_ip"`
				VtepMac  string `json:"vtep_mac"`
				VtepName string `json:"vtep_name"`
			} `json:"vxlan"`
		} `json:"backend"`
		DockerBridge struct {
			IP   string `json:"ip"`
			Name string `json:"name"`
		} `json:"docker_bridge"`
		Info struct {
			Enabled bool   `json:"enabled"`
			Name    string `json:"name"`
			Prefix  int64  `json:"prefix"`
			Subnet  string `json:"subnet"`
		} `json:"info"`
		MesosBridge struct {
			IP   string `json:"ip"`
			Name string `json:"name"`
		} `json:"mesos_bridge"`
		State struct {
			Error  string `json:"error"`
			Status string `json:"status"`
		} `json:"state"`
		Subnet string `json:"subnet"`
	} `json:"overlays"`
}
