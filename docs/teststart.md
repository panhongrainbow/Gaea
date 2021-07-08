 

# 开发当下的测试

## 执行测试

可以依照以下步骤，上传程式码前，进行自动化的测试

```bash
# 安装执行 make 所需要的套件
apt-get install build-essential

# 检查 go 语言的版本
go version
# go version go1.16.5 linux/amd64

# 在 Gaea 目录下执行测试
cd Gaea
make test
# 会印出以下内容
# go test -coverprofile=.coverage.out ./...
# ok      github.com/XiaoMi/Gaea/backend  0.037s  coverage: 9.1% of statements
# ?       github.com/XiaoMi/Gaea/backend/mocks    [no test files]
# ?       github.com/XiaoMi/Gaea/cc       [no test files]
# ?       github.com/XiaoMi/Gaea/cc/proxy [no test files]
# ?       github.com/XiaoMi/Gaea/cc/service       [no test files]
# ok      github.com/XiaoMi/Gaea/cmd/gaea 0.042s  coverage: 14.0% of statements
# ?       github.com/XiaoMi/Gaea/cmd/gaea-cc      [no test files]
# ?       github.com/XiaoMi/Gaea/core     [no test files]
# ?       github.com/XiaoMi/Gaea/core/errors      [no test files]
# ?       github.com/XiaoMi/Gaea/log      [no test files]
# ?       github.com/XiaoMi/Gaea/log/xlog [no test files]
# ok      github.com/XiaoMi/Gaea/models   0.011s  coverage: 76.4% of statements
# ok      github.com/XiaoMi/Gaea/models/etcd      0.036s  coverage: 4.4% of statements
# ?       github.com/XiaoMi/Gaea/models/file      [no test files]
# ok      github.com/XiaoMi/Gaea/mysql    0.034s  coverage: 34.3% of statements
# ok      github.com/XiaoMi/Gaea/parser   0.290s  coverage: 84.1% of statements
# ok      github.com/XiaoMi/Gaea/parser/ast       0.121s  coverage: 57.1% of statements
# ok      github.com/XiaoMi/Gaea/parser/auth      0.040s  coverage: 0.0% of statements [no tests to run]
# ok      github.com/XiaoMi/Gaea/parser/format    0.014s  coverage: 81.3% of statements
# ok      github.com/XiaoMi/Gaea/parser/model     0.033s  coverage: 53.6% of statements
# ok      github.com/XiaoMi/Gaea/parser/opcode    0.028s  coverage: 55.6% of statements
# ?       github.com/XiaoMi/Gaea/parser/stmtctx   [no test files]
# ok      github.com/XiaoMi/Gaea/parser/terror    0.064s  coverage: 75.9% of statements
# ok      github.com/XiaoMi/Gaea/parser/tidb-types        0.091s  coverage: 67.3% of statements
# ok      github.com/XiaoMi/Gaea/parser/tidb-types/json   0.030s  coverage: 81.4% of statements
# ok      github.com/XiaoMi/Gaea/parser/tidb-types/parser_driver  0.054s  coverage: 24.7% of statements
# ?       github.com/XiaoMi/Gaea/parser/types     [no test files]
# ok      github.com/XiaoMi/Gaea/proxy/plan       0.558s  coverage: 63.0% of statements
# ok      github.com/XiaoMi/Gaea/proxy/router     0.075s  coverage: 57.4% of statements
# ?       github.com/XiaoMi/Gaea/proxy/sequence   [no test files]
# ok      github.com/XiaoMi/Gaea/proxy/server     0.044s  coverage: 15.2% of statements
# ok      github.com/XiaoMi/Gaea/stats    0.009s  coverage: 70.2% of statements
# ok      github.com/XiaoMi/Gaea/stats/prometheus 0.111s  coverage: 97.1% of statements
# ok      github.com/XiaoMi/Gaea/util     9.342s  coverage: 50.3% of statements
# ok      github.com/XiaoMi/Gaea/util/bucketpool  0.211s  coverage: 93.9% of statements
# ok      github.com/XiaoMi/Gaea/util/cache       0.019s  coverage: 98.1% of statements
# ok      github.com/XiaoMi/Gaea/util/crypto      0.007s  coverage: 69.7% of statements
# ok      github.com/XiaoMi/Gaea/util/hack        0.002s  coverage: 55.0% of statements
# ?       github.com/XiaoMi/Gaea/util/math        [no test files]
# ok      github.com/XiaoMi/Gaea/util/requests    0.004s  coverage: 45.7% of statements
# ok      github.com/XiaoMi/Gaea/util/sync2       0.028s  coverage: 60.3% of statements
# ?       github.com/XiaoMi/Gaea/util/testleak    [no test files]
# ok      github.com/XiaoMi/Gaea/util/timer       0.623s  coverage: 96.2% of statements
# go tool cover -func=.coverage.out -o .coverage.func
# tail -1 .coverage.func
# total:                                                                          (statements)                                            59.2%
# go tool cover -html=.coverage.out -o .coverage.html
```

