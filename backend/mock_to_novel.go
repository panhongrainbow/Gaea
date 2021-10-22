package backend

import (
	"fmt"
	"github.com/XiaoMi/Gaea/mysql"
	"log"
	"os"
)

// switchMockResult 函式 🧚 为到不同的模拟数据库去找寻回应的讯息
func (fdb *fakeDB) switchMockResult(db string, key uint32) (*mysql.Result, error) {
	switch db {
	case "novel": // 29 本小说部份
		return fdb.switchNovelResult(key) // 在小说模拟数据库时去找对应到 SQL 字串的回应讯息
	}
	log.Fatal("没有命中模拟测试数据名称为: ", db) // 中断，因为测试程式有问题
	return &mysql.Result{}, nil
}

// switchNovelResult 函式 🧚 为在小说模拟数据库时去找对应到 SQL 字串的回应讯息
// 这里可以参考位于 Gaea/backend/mock_key_test.go 的测试函式 TestMockNovelKey
func (fdb *fakeDB) switchNovelResult(key uint32) (*mysql.Result, error) {
	switch key {
	// >>>>> >>>>> >>>>> >>>>> >>>>> 向多台数据库进行查询
	case 3717314451, 1196547673, 4270781616, 3210257868, 3606398974:
		/*
			向第一个切片进行查询
			所对应各切片 SQL 执行字串 以及 切片相关资讯
			数据库名称: novel
			模拟数据库的网路位置: 192.168.122.2:3309 或 192.168.122.2:3310 或 192.168.122.2:3311
			数据库执行字串: SELECT *,`BookID` FROM `novel`.`Book_0000` ORDER BY `BookID`
			数据库执行时所对应的 Key: 3717314451 或 1196547673 或 4270781616
		*/
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return fdb.MockDataInDB[0].result, nil
	case 2403537350, 1260331735, 1401931444:
		/*
			向第二个切片进行查询
			所对应各切片 SQL 执行字串 以及 切片相关资讯
			数据库名称: novel
			模拟数据库的网路位置: 192.168.122.2:3312 或 192.168.122.2:3313 或 192.168.122.2:3314
			数据库执行字串: SELECT *,`BookID` FROM `novel`.`Book_0001` ORDER BY `BookID`
			数据库执行时所对应的 Key: 2403537350 或 1260331735 或 1401931444
		*/
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return fdb.MockDataInDB[1].result, nil
	// >>>>> >>>>> >>>>> >>>>> >>>>> 向多台数据库进行写入
	case 1389454267: // 写入第一本小说到数据库 三国演义 (会分配到 Slice-1)
		ret, err := fdb.MockDataInDB[1].result.InsertFirstNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 618120042: // 写入第二本小说到数据库 水浒传 (会分配到 Slice-0)
		ret, err := fdb.MockDataInDB[0].result.InsertSecondNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 514659115: // 写入第三本小说到数据库 西游记 (会分配到 Slice-1)
		ret, err := fdb.MockDataInDB[1].result.InsertThirdNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 4273731942: // 写入第四本小说到数据库 红楼梦 (会分配到 Slice-0)
		ret, err := fdb.MockDataInDB[0].result.InsertFourthNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 4076192191: // 写入第五本小说到数据库 金瓶梅 (会分配到 Slice-1)
		ret, err := fdb.MockDataInDB[1].result.InsertFifthNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 1926088204: // 写入第六本小说到数据库 儒林外史 (会分配到 Slice-0)
		ret, err := fdb.MockDataInDB[0].result.InsertSixthNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 1572904758: // 写入第七本小说到数据库 初刻拍案惊奇 (会分配到 Slice-1)
		ret, err := fdb.MockDataInDB[1].result.InsertSeventhNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 1708424148: // 写入第八本小说到数据库 二刻拍案惊奇 (会分配到 Slice-0)
		ret, err := fdb.MockDataInDB[0].result.InsertEighthNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 3188314210: // 写入第九本小说到数据库 封神演义 (会分配到 Slice-1)
		ret, err := fdb.MockDataInDB[1].result.InsertNinethNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 3303343655: // 写入第十本小说到数据库 镜花缘 (会分配到 Slice-0)
		ret, err := fdb.MockDataInDB[0].result.InsertTenthNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 3599615497: // 写入第十一本小说到数据库 喻世明言 (会分配到 Slice-1)
		ret, err := fdb.MockDataInDB[1].result.InsertEleventhNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 600352469: // 写入第十二本小说到数据库 说岳全传 (会分配到 Slice-0)
		ret, err := fdb.MockDataInDB[0].result.InsertTwelfthNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 709958148: // 写入第十三本小说到数据库 杨家将 (会分配到 Slice-1)
		ret, err := fdb.MockDataInDB[1].result.InsertThirteenthNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 1226676578: // 写入第十四本小说到数据库 说唐 (会分配到 Slice-0)
		ret, err := fdb.MockDataInDB[0].result.InsertFourteenthNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 56203336: // 写入第十五本小说到数据库 七侠五义 (会分配到 Slice-1)
		ret, err := fdb.MockDataInDB[1].result.InsertFifteenthNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 3585696861: // 写入第十六本小说到数据库 施公案 (会分配到 Slice-0)
		ret, err := fdb.MockDataInDB[0].result.InsertSixteenthNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 3821388015: // 写入第十七本小说到数据库 青楼梦 (会分配到 Slice-1)
		ret, err := fdb.MockDataInDB[1].result.InsertSeventeenthNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 1792929480: // 写入第十八本小说到数据库 歧路灯 (会分配到 Slice-0)
		ret, err := fdb.MockDataInDB[0].result.InsertEighteenthNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 398747927: // 写入第十九本小说到数据库 老残游记 (会分配到 Slice-1)
		ret, err := fdb.MockDataInDB[1].result.InsertNineteenthNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 1187323765: // 写入第二十本小说到数据库 二十年目睹之怪现状 (会分配到 Slice-0)
		ret, err := fdb.MockDataInDB[0].result.InsertTwentiethNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 1498815330: // 写入第二十一本小说到数据库 孽海花 (会分配到 Slice-1)
		ret, err := fdb.MockDataInDB[1].result.InsertTwentyFirstNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 2032678570: // 写入第二十二本小说到数据库 官场现形记 (会分配到 Slice-0)
		ret, err := fdb.MockDataInDB[0].result.InsertTwentySecondNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 2614046017: // 写入第二十三本小说到数据库 觉世名言十二楼 (会分配到 Slice-1)
		ret, err := fdb.MockDataInDB[1].result.InsertTwentyThirdNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 2457093340: // 写入第二十四本小说到数据库 无声戏 (会分配到 Slice-0)
		ret, err := fdb.MockDataInDB[0].result.InsertTwentyFourthNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 238477972: // 写入第二十五本小说到数据库 肉蒲团 (会分配到 Slice-1)
		ret, err := fdb.MockDataInDB[1].result.InsertTwentyFifthNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 4020693348: // 写入第二十六本小说到数据库 浮生六记 (会分配到 Slice-0)
		ret, err := fdb.MockDataInDB[0].result.InsertTwentySixthNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 2745523730: // 写入第二十七本小说到数据库 野叟曝言 (会分配到 Slice-1)
		ret, err := fdb.MockDataInDB[1].result.InsertTwentySeventhNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 1776512190: // 写入第二十八本小说到数据库 九尾龟 (会分配到 Slice-0)
		ret, err := fdb.MockDataInDB[0].result.InsertTwentyEighthNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	case 424563096: // 写入第二十九本小说到数据库 品花宝鉴 (会分配到 Slice-1)
		ret, err := fdb.MockDataInDB[1].result.InsertTwentyNinethNovelResult()
		if err != nil {
			return nil, err
		}
		tool := os.Getenv("IDE_TOOL")
		if tool == "jetbrains" {
			fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		}
		return ret, nil
	}

	// >>>>> >>>>> >>>>> >>>>> >>>>> 如果没有命中 key 值的时候，就直接中断整个测试
	log.Fatalf("\u001B[35m 没有命中模拟测试 Key 为: %d\n", key) // 中断，因为测试程式有问题
	return nil, nil
}
