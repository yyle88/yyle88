## Hi there 👋热烈欢迎贵宾光临 👋

这是我的介绍:

- 😄 杨亦乐
- 🔭 1990
- 🌱 男
- 👯 西安电子科技大学 (2010-2014)
- 🤔 主要使用 Go 语言
- 💬 我是个兴趣使然的程序员
- 📫 偶尔有空时写写开源代码
- 😄 希望你给我点星星
- ⚡ Give me stars. Thank you!!! 谢谢大家。

我的项目列表：
[我的项目列表](https://github.com/yyle88?tab=repositories&sort=stargazers)
给我星星谢谢！

这是我的项目：

| 项目名称 | 项目描述 |
|-------------------------------------------------|--------|
| [gotrontrx](https://github.com/yyle88/gotrontrx) | 简单的Go语言版本波场创建钱包/获得测试币/发出交易的逻辑 |
| [gobtcsign](https://github.com/yyle88/gobtcsign) | 简单的Go语言版本比特币创建钱包和签名逻辑 |
| [gormcngen](https://github.com/yyle88/gormcngen) | gormcngen help gen enum code with gormcnm. 因为gormcnm能够以枚举定义字段名和字段类型，这个工具gormcngen就是自动帮你生成枚举代码的。 |
| [reggin](https://github.com/yyle88/reggin) | reggin means register gin routes. 非常简单的gin路由注册器。 |
| [gormcnm](https://github.com/yyle88/gormcnm) | gormcnm means: gorm column name. can help you use enum column name and enum type define. not use raw string. 该工具能枚举gorm的字段名和列类型，这样能避免gorm使用者在项目中写太多原始字符串，也就能避免写错，当您需要修改gorm字段时，该工具也能帮助您在编码阶段或编译阶段就发现问题，使得重构gorm类型更简单更轻松更可靠。当然该工具还提供一些简单的逻辑，使得查询操作变得更容易，但这都是建立在您已经充分掌握gorm基本增删改查操作的基础上的。我自己开发出来以后觉得非常好用，就顺带把它开源出来，以方便以后我做别的项目，也方便大家，希望大家能给星星哦，谢谢大家。 |
| [syncmap](https://github.com/yyle88/syncmap) | 该工具包100%封装sync.Map的方法，而且方法的参数和返回值都保持不变。  |
| [formatgo](https://github.com/yyle88/formatgo) | 格式化代码 gofmt 工具 format golang source code 的工具，当然顺带还能整理 import 的引用内容 |
| [done](https://github.com/yyle88/done) | 在golang代码里常有 res, err := run(); err != nil 的逻辑，错误出现的概率很小，但处理错误会让代码变得臃肿，特别是在写很小的demo时，出错就直接panic就行，这个包就是提出了个错误处理的新方案。让你的代码能够很简洁。在项目起名方面，使用checkerk等过长，而使用goerrcheck等就像要启动新的goroutine，因此最终选择使用done这个项目名，表示顺利完成某件事，虽然welldone也行，但done更简洁。 |
| [gormcls](https://github.com/yyle88/gormcls) | 跟前面的gormcnm配合使用，最终效果是很棒的 |
| [zaplog](https://github.com/yyle88/zaplog) | 自己做东西总是需要打印日志的，但似乎各个开发者用的都不一样，我本来想着能不能统一下日志，结果发现不能，因此这个包就作为个仅供自己用的日志包吧（主要是没有这个包，其它包想开源也开不起来啊，这就比较尴尬啦）。 |
| [mutexmap](https://github.com/yyle88/mutexmap) | 跟syncmap不同，这个是一个rw-mutex和一个map的组合，目的是解决map的异步读写问题，这个比较鬼扯，查别人已有的代码也行，但不如自己顺手实现个完事 |
| [syntaxgo](https://github.com/yyle88/syntaxgo) | 就是golang的ast语法分析树和golang的reflect反射包的封装，让你更方便的去分析代码，最终实现自动生成新功能代码的效果 |
| [runpath](https://github.com/yyle88/runpath) | 获取正在执行的golang代码的位置信息，即 execution location，即源代码go文件在电脑里的绝对路径和行号，使用 "runtime" 获得，因此包名起名为 "runpath" 即可，而不使用比较长的 executionlocation，但含义就是这样的 |
| [gormmom](https://github.com/yyle88/gormmom) | 添加使用母语编写gorm模型的功能 |
| [sure](https://github.com/yyle88/sure) | 在我们开发golang代码时，经常会遇到比如 res, err := a.Run() 的情况，这时假如使用 res := amust.Run() 或者 res := a.Must().Run()岂不是能够避免频繁的判断 if err != nil 啦，这个包的目的就是提供这样的便利，当然本整活大师开发的 `github.com/yyle88/done` 也能解决问题，但毕竟不是还得多一层`nice`调用嘛，而这个工具将让代码自己提供错误时panic/ignore的选项，当然包名的话在mustsoft和mustgo和flexible间选择半天，最终想到也可以和`github.com/yyle88/done`套套近乎干脆就叫`mustdone`吧 |
| [sortslice](https://github.com/yyle88/sortslice) | 简单的排序逻辑，使用泛型实现 sort.Interface 这样以后排序就不要每次都根据类型实现 sort.Interface 啦，非常方便，给个星星谢谢 |
| [erero](https://github.com/yyle88/erero) | 简单的errors包，和菠萝菠萝蜜的相同，erero，就是个简单的错误包，当发生错误的时候自动打印日志，假如名字叫errors就有点烂大街啦，还得解决包名的冲突问题，比如和标准errors或者github.com/pkg/errors的冲突。因此随便起个名字吧，假如叫ero就过短啦，还容易和变量名冲突。因此就用erero吧。 |
| [osexistpath](https://github.com/yyle88/osexistpath) | 检查路径是否存在，路径文件是否存在，路径的目录是否存在。因为没有开源包来专门做这件小事，就由我来做吧。 |
| [neatjson](https://github.com/yyle88/neatjson) | neat json make it neat to use "encoding/json" in golang. |
| [demojavabtcsign](https://github.com/yyle88/demojavabtcsign) | 使用Java给BTC签名的DEMO 当然由于我在开发时顺带也接了狗狗币dogecoin，因此这里也同样可以适用于狗狗币的签名（跟BTC签名共用逻辑，区别仅仅在于，链的网络参数不同）。 |
| [yyle88](https://github.com/yyle88/yyle88) | Go Go Go |
| [must](https://github.com/yyle88/must) | must means assert means require. while the assert/require are using in testcase. must is using in main code. |

给我星星谢谢。

这是我的账号：

| 账号名称                                            | 账号描述   |
|-------------------------------------------------|--------|
| [yyle88](https://github.com/yyle88)             | 主要开源账号 |
| [yyle66](https://github.com/yyle66)             | 备用开源账号 |
| [yangyile1990](https://github.com/yangyile1990) | 其他开源账号 |

这是我的组织：

| 组织名称                                                   | 组织描述                              |
|--------------------------------------------------------|-----------------------------------|
| [go-go-go](https://github.com/yyle88?tab=repositories) | 我的项目列表                            |
| [go-xlan](https://github.com/go-xlan)                  | 使用Go语言接其他语言/环境/网络/服务/客户端/协议/工具的代码 |
| [orzkratos](https://github.com/orzkratos)              | 使用Go-Kratos框架的心得体会工具代码            |
