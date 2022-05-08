package containerTest

import (
	"encoding/json"
	"fmt"
	"github.com/XiaoMi/Gaea/util"
	"github.com/XiaoMi/Gaea/util/mocks/containerTest/builder/containerd"
	"github.com/XiaoMi/Gaea/util/mocks/containerTest/builder/containerd/run"
	"io/ioutil"
	"net"
	"path/filepath"
	"strconv"
	"strings"
)

// 进行介面的设置

// Load 为用来执行容器服务
type Load struct {
	client map[string]run.Run
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
func (r *Load) loadContainerD(file string) (containerd.ContainerD, error) {
	path := r.ContainerdPath()
	config := filepath.Join(path, file)
	b, err := ioutil.ReadFile(config)
	if err != nil {
		return containerd.ContainerD{}, err
	}
	c := containerd.ContainerD{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return containerd.ContainerD{}, err
	}
	return c, nil
}

// loadAllContainerD 把所有的設定檔轉成 ContainerD 設定值
func (r *Load) loadAllContainerD() (map[string]containerd.ContainerD, error) {
	files, err := r.listContainerD()
	if err != nil {
		return nil, err
	}
	configs := make(map[string]containerd.ContainerD, len(files))
	for _, f := range files {
		config, err := r.loadContainerD(f)
		if err != nil {
			return nil, err
		}
		configs[config.Name] = config
	}

	return configs, nil
}

// >>>>> >>>>> >>>>> 扩展容器服务的设定 extend the containerd config

// correctRange 是用来修正表示范围的字符串
// correctRange is used to correct the range string
func correctRange(name string) string {
	// 先修正字符串 fix the string
	name = strings.TrimSpace(name)
	name = strings.Replace(name, "TO", "To", -1)
	name = strings.Replace(name, "to", "To", -1)
	return name
}

// extendContainerName 是用来扩展容器名称的
// extendContainerName is used to extend the container name
func extendContainerName(name string) ([]string, error) {
	// 先修正字符串 fix the string
	name = correctRange(name)

	// 如没有 "To" "{" 或 "}" 如果有就返回错误 return error
	// if there is no "To" "{" or "}", return error
	if !strings.Contains(name, "To") ||
		!strings.Contains(name, "{") ||
		!strings.Contains(name, "}") {
		return nil, fmt.Errorf("the name %s is not valid", name)
	}

	// 先决定 { To 和 } 的位置 find the position of {, To, and }
	indexOpenCurly := strings.LastIndex(name, "{")
	indexCloseCurly := strings.LastIndex(name, "}")
	indexTo := strings.LastIndex(name, "To")

	// 决定前缀、中缀、后缀 prefix, middle, suffix
	prefix := name[:indexOpenCurly]
	middle := name[indexOpenCurly+1 : indexTo]
	suffix := name[indexTo+2 : indexCloseCurly]

	// 把 middleInt 和 middleFloat 转成 int 和 float convert middleInt and middleFloat to int and float
	middleInt, err := strconv.Atoi(middle)
	if err != nil {
		return nil, err
	}

	suffixInt, err := strconv.Atoi(suffix)
	if err != nil {
		return nil, err
	}

	// 收集所有的名称 collect all names
	var names = make([]string, suffixInt-middleInt+1)
	var j = 0

	// 开始收集 names。 collect names
	for i := middleInt; i <= suffixInt; i++ {
		names[j] = prefix + strconv.Itoa(i)
		j++
	}

	// 回传 names 回传 names。 return names
	return names, nil
}

// separateIPandPort 分离 IP 和 Port
// separateIPandPort is used to seperate IP and Port
func separateIPandPort(ip string) (string, string) {
	// 找到分离 IP 和 Port 的位置 find the position of seperating IP and Port.
	index := strings.LastIndex(ip, ":")

	// 如果 index 是 -1，那么就没有 ":" 在 ip 中 if index is -1, it means that there is no ":" in the ip.
	if index == -1 {
		return ip, ""
	}

	// 回传 IP 和 Port 回传 IP 和 Port。 return IP and Port
	return ip[:index], ip[index+1:]
}

// extendContainerIP 是用来扩展容器网路位置的
// extendContainerIP is used to extend the container ip
func extendContainerIP(ipStr string) ([]string, error) {
	// 先修正字符串 fix the string
	// 分离 IP 和 Port。 separate IP and Port
	ipPortRange := strings.Split(correctRange(ipStr), "To")
	firstIP, firstPort := separateIPandPort(ipPortRange[0])
	endIP, endPort := separateIPandPort(ipPortRange[1])
	if firstPort != endPort {
		return nil, fmt.Errorf("the port %s is not match", endPort)
	}

	// ip 阵列。 ip array
	ips := make([]string, 0)

	// 下一个 IP 和 Port。 next IP and Port
	nextNetIP := net.ParseIP(firstIP)
	for {
		ips = append(ips, nextNetIP.String()+":"+firstPort)
		if nextNetIP.String() == endIP {
			// 如果 nextIP 字符串和 endIP 相同，那么就退出退出循环 break the loop
			// if nextIP string is equal to endIP, then exit the loop
			break
		}
		// 下一个 IP 地址 nextIP
		nextNetIP = util.IncrementIP(nextNetIP)
	}

	// 回传 IP 阵列。 return IP array
	return ips, nil
}

func extendContainerConfig(config containerd.ContainerD) ([]containerd.ContainerD, error) {

	// names
	names, err := extendContainerName(config.Name)
	if err != nil {
		return nil, err
	}
	// ip
	ips, err := extendContainerIP(config.IP)

	// snapshot
	snapshots, err := extendContainerName(config.SnapShot)

	extendConfigs := make([]containerd.ContainerD, len(names))

	for i, name := range names {
		extendConfigs[i] = containerd.ContainerD{
			Sock:      config.Sock,
			Type:      config.Type,
			Name:      name,
			NameSpace: snapshots[i],
			Image:     config.Image,
			Task:      config.Task,
			NetworkNs: config.NetworkNs,
			IP:        ips[i],
			SnapShot:  config.SnapShot,
			Schema:    config.Schema,
			User:      config.User,
			Password:  config.Password,
		}
	}

	return extendConfigs, nil
}
