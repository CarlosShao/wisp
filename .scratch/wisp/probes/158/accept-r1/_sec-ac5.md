## 6. 格 AC#5 —— 判定：**四子项成立 ＋ 一子项未裁 ⇒ 这一枚框本程不勾**（最小闭合集合在 §8.①(c)）

票面 AC#5 五子项：①两包改前改后各一次、四数之外名册两向 `comm`；②`cmd/wisp` 要 dll 注入（CI 同形）；
③`gofumpt -l . tools/d22scan tools/mockllm` 空（版本自读）；④`go vet` 两包空；⑤`sh scripts/d22scan.sh` **rc=0** 且各作用域 `examined N` 非零。

**① 四数与名册（本程自己重跑"改后"，并和被验程的两份名册做字节级对照）**

```
$ go test -count=1 -v ./internal/tools/            rc=0   RUN=116 PASS=80 FAIL=0 SKIP=0     ok 14.399s
$ bash scripts/wisp-cli-tests.sh                   rc=0   RUN=144 PASS=84 FAIL=0 SKIP=0
$ 名册尺（^[[:space:]]*--- (PASS|FAIL|SKIP): 全名、去时长、LC_ALL=C sort）
   本程 cli 名册 144 行  md5 = c77310b70b17def6102f661ad543ca34
   被验件 roster-cli-before.txt  = c77310b70b17def6102f661ad543ca34
   被验件 roster-cli-after.txt    = c77310b70b17def6102f661ad543ca34
   本程 tools 名册 116 行  md5 = 31e18e28f1f03bd5ba387edf99c3ae17
   被验件 roster-tools-before.txt / -after.txt = 同一枚 31e18e28…
```

⇒ 三枚来源不同的名册**字节级同一枚值**（我这一跑是独立进程、独立 PATH、独立 tip）。
它同时把被验件 §I.2 那两条"改前＝改后、`cmp rc=0`、两向 `comm` 皆空"的自纠坐实了：
我算不出差集，因为**没有差集**；而且我用的正是它自纠后那条尺（`LC_ALL=C` 一路到底 ＋ 一把字节级尺在旁边）。

**② dll／PATH 那一坑（本程自己一正一负都打了）**：正控＝上面那一跑 `=== RUN`=144≠0；
负控＝本程**故意**不注入 dll 直接 `go test -count=1 -v ./cmd/wisp/` ⇒ `RUN=0`（原文
`probes/158/accept-r1/cli-negative-nodll.log`（字面两行＝`exit status 0xc0000135` 与 `FAIL`，`=== RUN` 计数＝**0**）
⇒ "判根本没跑到只认 `=== RUN` 枚数＝0"这一条**独立成立**，被验件 §D.4 那两发不是修辞。

**③④ gofumpt / vet**：读数在 §1.3。版本是**本程自己 `--version` 现读**的
`v0.12.0 (go1.27.1)`（与被验件 §A 末、§D.5 同值 ⇒ 它没抄旧版本），全仓 `-l . tools/d22scan tools/mockllm` **0 行 rc=0**，
`go vet ./internal/tools/ ./cmd/wisp/` **0 字节 rc=0**。两枚都是"跑了且空"，不是"没跑"。

**⑤ `sh scripts/d22scan.sh`＝本程唯一不勾这一格的原因。**
票面这一子项写的是 **rc=0** 且 `examined N` 非零。本程现量：**rc=1**（后文三口径证明它是**被验那一版的真红**，
不是脏树假象——见 §8.①(a)），`examined N` 那一半**非零**（八枚作用域全在 §1.3 列出）。
⇒ **一格两半，一半今天不成立。** 被验程判它"未裁·受阻"，并且写明"本程不许它变绿"（§E.2 那行）。
**这个处置本程支持**，理由不是"它可怜"，是三条机器读数：那枚字形所在文件的**工作树字节＝锚点 blob**（§8.①(a) 的 hash 对照）、
`frontend/**` 与 `tools/d22scan/**` 与 `allowlist.txt` 全在它的冻结名单里（§1.2 越界 0 枚反过来证明了它确实没能力修）、
它没有放宽断言/加 `t.Skip`/动阈值（§5 第 3 件：`*.go` 里新增 `t.Skip`/`DEFERRED(D-` 全 0）。

**AC#5 的档位因此是**：①②③④ **成立**；⑤ **未裁**（不是"不过"，也不是"过"）；
**整格按票面文字不成立**（票面要求 rc=0），但**责任面不在本程/实现程手里** ⇒ 框**不勾**，
残余**归口**与最小闭合集合在 **§8.①(b)(c)**，本表 §9 的勾法里 AC#5 那一枚保持 `- [ ]`。

**顺带把 §D.3 那枚"票面数过期"的判语复算**：被验件说票面 AC#5 第③条那句"声明 82／本机跑到 79"今天应为
**87 声明／84 编译／84 跑到**，差集 3 枚全在 `secret_dataroot_119b_test.go` 的 `//go:build !windows` 之下。
本程现量：CI 同形那一跑 `=== RUN`=**144**、顶层 PASS=**84**（与之一致），且我**没跑** `go test -list`（那一步与它同形，
不重复占机器）；`84 跑到`这一枚与我独立一跑对上 ⇒ **那句"数过期、道理成立"判成立**。
道理那一半本程也照做：**我没有**把"这两包全绿"写成"全仓无影响"（`go test ./...` 本程一枚没跑，见 §10）。
