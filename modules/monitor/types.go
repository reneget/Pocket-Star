package monitor

type CheckType string

const (
	CheckPeerAlive    CheckType = "peer.alive"
	CheckLatency      CheckType = "peer.latency"
	CheckCPU          CheckType = "system.cpu"
	CheckRAM          CheckType = "system.ram"
	CheckDisk         CheckType = "system.disk"
	CheckPublicIP     CheckType = "network.public_ip"
	CheckContainer    CheckType = "docker.container"
	CheckContainerCPU CheckType = "docker.cpu"
	CheckContainerRAM CheckType = "docker.ram"
)

type AlertLevel int

const (
	AlertWarning  AlertLevel = iota
	AlertCritical
)

type CheckResult struct {
	Type      CheckType `json:"type"`
	OK        bool      `json:"ok"`
	Value     float64   `json:"value,omitempty"`
	Label     string    `json:"label,omitempty"`
	Timestamp int64     `json:"ts"`
	Message   string    `json:"msg,omitempty"`
}

type Alert struct {
	Level     AlertLevel `json:"level"`
	Check     CheckType  `json:"check"`
	Label     string     `json:"label"`
	Value     float64    `json:"value"`
	Threshold float64    `json:"threshold"`
	Message   string     `json:"message"`
	Time      int64      `json:"ts"`
}

type AlertRule struct {
	Check   CheckType `json:"check"`
	Label   string    `json:"label,omitempty"`
	Warning float64   `json:"warning"`
	Critical float64  `json:"critical"`
	Enabled bool      `json:"enabled"`
}

type ContainerInfo struct {
	Name       string  `json:"name"`
	Image      string  `json:"image"`
	Status     string  `json:"status"`
	Uptime     string  `json:"uptime"`
	CPUPercent float64 `json:"cpu_percent"`
	MemPercent float64 `json:"mem_percent"`
	MemUsage   string  `json:"mem_usage"`
	Restarts   int     `json:"restarts"`
	Running    bool    `json:"running"`
}

type NodeHealth struct {
	Label      string  `json:"label"`
	IP         string  `json:"ip"`
	Alive      bool    `json:"alive"`
	Latency    float64 `json:"latency_ms"`
	LastSeen   int64   `json:"last_seen"`
}

type SystemHealth struct {
	CPUPercent  float64 `json:"cpu_percent"`
	RAMPercent  float64 `json:"ram_percent"`
	RAMUsed     string  `json:"ram_used"`
	RAMTotal    string  `json:"ram_total"`
	DiskPercent float64 `json:"disk_percent"`
	DiskUsed    string  `json:"disk_used"`
	DiskTotal   string  `json:"disk_total"`
	PublicIP    string  `json:"public_ip"`
}
