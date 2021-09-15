package backend

import (
	"fmt"
	"github.com/XiaoMi/Gaea/mysql"
	"log"
	"strconv"
	"strings"
)

// InitSwitchTrans 函式 🧚 为 在单元测试数据库时决定要使用何种数据库模拟资料
func (dc *DirectConnection) initSwitchTrans() (string, error) {
	// 得知要使用的数据库 (正确的做法，手動指定)
	/*if err := dc.Trans.UseDB("novel"); err != nil { // 29本小说资料
		return dc, err
	}*/

	// 得知要使用的数据库 (错误的做法，自动载入)
	/*if err := dc.Trans.UseDB(dc.db); err != nil { // 因为上层函式并不会传送数据库名称到 dc.db 变数里
		return dc, err
	}*/

	// 由网路位置取出埠号
	tmp := strings.Split(dc.addr, ":")
	port, err := strconv.Atoi(tmp[1])
	if err != nil {
		return "", err
	}

	// 根据测试埠号去载入相关模拟数据库
	switch {
	// 将来要抽换制造假资料的方法，就直接在这里抽换就好，这是唯一要修改的地方
	case (3309 <= port) && (port <= 3311): // 第二丛集 Port 3309 ~ 3311
		// 决定要使用何种数据库模拟资料
		dc.Trans = new(novelData)                           // 29本小说资料
		if err := dc.Trans.UseFakeDB("novel"); err != nil { // 29本小说资料
			return "", err
		}
		// 初始化数据库模拟资料
		if _, ok := fakeDBInstance["novel"]; !ok { // 看 fakeDBInstance map 里的 key 存不存在就知道模拟数据是否有初始化完成
			if fakeDBInstance == nil {
				fakeDBInstance = make(map[string]*fakeDB)
			}
			fakeDBInstance["novel"] = new(fakeDB)
			fakeDBInstance["novel"].MockDataInDB = make([]mysql.Result, 2, 2) // Slice 不用在扩张了，小说资料只会被分成二个切片
			fmt.Println("切片長度", len(fakeDBInstance["novel"].MockDataInDB))
		}
		return "novel", err
		// 将来要抽换制造假资料的方法，就直接在这里抽换就好，这是唯一要修改的地方
	case (3312 <= port) && (port <= 3314): // 第三丛集 Port 3312 ~ 3314
		// 决定要使用何种数据库模拟资料
		dc.Trans = new(novelData)                           // 29本小说资料
		if err := dc.Trans.UseFakeDB("novel"); err != nil { // 29本小说资料
			return "", err
		}
		// 初始化数据库模拟资料
		if _, ok := fakeDBInstance["novel"]; !ok { // 看 fakeDBInstance map 里的 key 存不存在就知道模拟数据是否有初始化完成
			if fakeDBInstance == nil {
				fakeDBInstance = make(map[string]*fakeDB)
			}
			fakeDBInstance["novel"] = new(fakeDB)
			fakeDBInstance["novel"].MockDataInDB = make([]mysql.Result, 2, 2) // Slice 不用在扩张了，小说资料只会被分成二个切片
			fmt.Println("切片長度", len(fakeDBInstance["novel"].MockDataInDB))
		}
		return "novel", err
	}

	// 都没命中埠的事后的处理
	log.Fatal("没有命中模拟测试数据库的埠号为: ", port) // 中断，因为测试程式有问题
	return "", nil
}
