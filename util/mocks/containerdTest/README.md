现在有一个大问题，数据库资料要如何灌到容器
目前是做一个新的本地印象档，容器启动时就自动载入，能少一个依赖少一个依赖

决定使用 containerd 加入测试，docker 未来没有这么重要了
如果把测试和 docker 整合在一起，这样没有跟上科技的脚步
经测试后，决定把 containerd 整个到测试中，并把 soarTest 改名成 containerdTest