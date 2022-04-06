package containerdTest

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

// 进行介面的设置

// Load 为用来执行容器服务
type Load struct {
	client map[string]Run
	prefix string
}

// ContainerdPath 返回容器服务设定目录的路径
func (r *Load) ContainerdPath() string {
	return filepath.Join(r.prefix, "containerd")
}

// listContainerD 列出 ContainerD 的設定檔列表
func (r *Load) listContainerD() ([]string, error) {
	path := r.ContainerdPath()

	list := make([]string, 0)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		list = append(list, f.Name())
	}

	return list, nil
}

// loadContainerD 载入指定的設定檔轉成 ContainerD 設定值
func (r *Load) loadContainerD(file string) (ContainerD, error) {
	path := r.ContainerdPath()
	config := filepath.Join(path, file)
	b, err := ioutil.ReadFile(config)
	if err != nil {
		return ContainerD{}, err
	}
	c := ContainerD{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return ContainerD{}, err
	}
	return c, nil
}

// loadAllContainerD 把所有的設定檔轉成 ContainerD 設定值
func (r *Load) loadAllContainerD() (map[string]ContainerD, error) {
	files, err := r.listContainerD()
	if err != nil {
		return nil, err
	}
	configs := make(map[string]ContainerD, len(files))
	for _, f := range files {
		config, err := r.loadContainerD(f)
		if err != nil {
			return nil, err
		}
		configs[config.Name] = config
	}

	return configs, nil
}
