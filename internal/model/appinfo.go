package model

type AppInfo struct {
	Version     string `json:"version"`
	GitCommit   string `json:"git_commit"`
	BuildTime   string `json:"build_time"`
	Environment string `json:"environment"`
	PodName     string `json:"pod_name"`
	NodeName    string `json:"node_name"`
	HostIP      string `json:"host_ip"`
	GoVersion   string `json:"go_version"`
	Uptime      string `json:"uptime"`
}
