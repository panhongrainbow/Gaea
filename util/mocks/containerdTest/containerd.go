package containerdTest

// 先進行 Parse 操作

// ContainerD 为 Containerd 容器服务所要 Parse 对应的内容
type ContainerD struct {
	Type      string `json:"type"`      // 容器服务类型
	Name      string `json:"name"`      // 容器服务名称
	Image     string `json:"image"`     // 容器服务镜像
	Container string `json:"container"` // 容器服务容器名称
	Task      string `json:"task"`      // 容器服务任务名称
	IP        string `json:"ip"`        // 容器服务 IP
	Schema    string `json:"schema"`    // 容器服务 Schema，用於數據庫設定
	User      string `json:"user"`      // 容器服务用户名
	Password  string `json:"password"`  // 容器服务密码
}
