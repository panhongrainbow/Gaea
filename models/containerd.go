package models

// Containerd 容器服务设定值
type Containerd struct {
	ConfigPath string `ini:"config_path"`
}

// CMConfig 容器服务管理器配置
type CMConfig struct {
	Enable string `ini:"containerd_enable"` // 容器服务是否启用
	Path   string `ini:"containerd_path"`   // 容器服务路径
}

// ParseCMConfig 解析 Containerd 容器服务管理员配置
/*func ParseCMConfig(cfgFile string) (*CMConfig, error) {
	cfg, err := ini.Load(cfgFile)

	if err != nil {
		return nil, err
	}

	ccConfig := new(CCConfig)
	err = cfg.MapTo(ccConfig)
	if ccConfig.DefaultCluster == "" && ccConfig.CoordinatorRoot != "" {
		ccConfig.DefaultCluster = strings.TrimPrefix(ccConfig.CoordinatorRoot, "/")
	}
	if ccConfig.CoordinatorType == "" {
		ccConfig.CoordinatorType = ConfigEtcd
	}
	return ccConfig, err
}*/
